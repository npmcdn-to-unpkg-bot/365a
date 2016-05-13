package user

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"fmt"
	"net/http"
	"time"
	"log"

	"github.com/asdine/storm"
	"github.com/gin-gonic/gin"
)

const (
	IrisContextField = "Session"
	XSRFCookieName   = "XSRF-TOKEN"
	TokenHeaderField = "X-XSRF-TOKEN"
)

var (
	SignInErr = errors.New("Sign in error")
)

type (
	SuccessResponse struct {
		Status string
		Data   interface{}
	}

	FailResponse struct {
		Status string
		Err    string
	}

	UserIDData struct {
		ID string
	}

	Session struct {
		Id int  `storm:"id"`
		Token   string  `storm:"index"`
		UserID  int
		Expires time.Time
	}

	IUser interface {
		ID() int
		PASSWORD() string
	}

	FindUser        func(string) (IUser, error)
	ConvertPassword func(string) string
)


func IsAllowed(roles ...UserRole)  func(c *gin.Context){

	return func(c *gin.Context){
		v, ok := c.Get(IrisContextField)
		if !ok {
			c.AbortWithStatus(401)
			return
		}

		s, ok := v.(Session)
		if !ok {
			c.AbortWithStatus(401)
			return

		}
		var user User
		err:=stormdb.One("Id",s.UserID,&user)
		if err!=nil{
			c.AbortWithError(500,err)
			return
		}

		var valid bool


		for i:= range roles{

            //lower rule higher permission
			if user.Role<=roles[i]{
				valid=true
			}
		}

		if !valid{
			c.AbortWithStatus(401)
			return

		}


		c.Next()
	}


}

func NewSuccessResponse(data interface{}) SuccessResponse {
	return SuccessResponse{
		Status: "success",
		Data:   data,
	}
}

func NewFailResponse(err interface{}) FailResponse {
	return FailResponse{
		Status: "fail",
		Err:    fmt.Sprintf("%v", err),
	}
}

func NewSessionToken() (string, error) {
	buf := make([]byte, 20)

	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}

	c := sha256.New()
	hash := fmt.Sprintf("%x", c.Sum(buf))

	return hash, nil
}

func NewSha512Password(pass string) string {
	hash := sha512.New()
	tmp := hash.Sum([]byte(pass))
	passHash := fmt.Sprintf("%x", tmp)
	return passHash
}

func ReadSession(ctx *gin.Context) (Session, error) {

	v, ok := ctx.Get(IrisContextField)
	if !ok {
		return Session{}, errors.New("Wrong session in cookie")
	}

	s, ok := v.(Session)
	if !ok {
		return Session{}, errors.New("Wrong session in cookie")
	}

	return s, nil
}

func AngularAuth(db *storm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := Auther(c, db)
		if err != nil {
			c.JSON(http.StatusUnauthorized, NewFailResponse(err))
			return
		}
	}
}

func Auther(c *gin.Context, db *storm.DB) error {
	token := c.Request.Header.Get(TokenHeaderField)
	if token == "" {
		cookie, err := c.Request.Cookie(XSRFCookieName)
		if err != nil {
			return errors.New("Cookie not found")
		}
		token = cookie.Value
		if token == "" {
			return errors.New("Header not found")
		}
	}

	var sess Session

	err := db.One("Token", token, &sess)
	if err != nil {
		return err
	}

	if &sess == nil {
		return errors.New("Session not found")

	}

	if sess.Expires.Before(time.Now()) {
		return errors.New("Session expired")
	}

	c.Set(IrisContextField, sess)
	c.Next()
	return nil
}

func AngularSignIn(coll *storm.DB, findUser FindUser, cPass ConvertPassword, expireTime time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := Signer(c, coll, findUser, cPass, expireTime)
		if err != nil {
			c.JSON(http.StatusUnauthorized, NewFailResponse(err))
			log.Println(err)
		}
	}
}

func Signer(c *gin.Context, db *storm.DB, findUser FindUser, convertPassword ConvertPassword, expireTime time.Duration) error {

	type auth struct {
		Email    string
		Password string
	}

	var a auth

	err := c.BindJSON(&a)
	if err != nil {
		return err
	}

	passHash := convertPassword(a.Password)

	user, err := findUser(a.Email)
	if err != nil {

		return err
	}



	if user.PASSWORD() != passHash {
		return SignInErr
	}





	resp := NewSuccessResponse(user.(*User))

	sessionToken, err := NewSessionToken()
	if err != nil {
		return err
	}

	expire := time.Now().Add(expireTime)

	session := Session{
		UserID:  user.ID(),
		Token:   sessionToken,
		Expires: expire,
	}

	err = db.Save(&session)
	if err != nil {
		c.JSON(500, err.Error())
		return err
	}

	cookie := http.Cookie{
		Name:    XSRFCookieName,
		Value:   sessionToken,
		Expires: expire,
		Path: "/",
	}


	http.SetCookie(c.Writer.(http.ResponseWriter), &cookie)

	c.JSON(200, resp)
	return nil
}
