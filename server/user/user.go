package user

import (
	"time"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/boltdb/bolt"
	"gopkg.in/mgo.v2/bson"
	"strconv"
)

type UserRole uint8

const (
	UserAdmin UserRole = iota+1
	UserSupervisor
	UserWorker
)

type User struct {
	Id             int `storm:"id"`
	Email          string `storm:"unique"`
	Salt           string
	Role           UserRole
	Name           string

	Password string        `storm:"index"`
	Image    *Base64Upload `storm:"inline"`
	Created  time.Time     `storm:"index"`
	Updated  time.Time     `storm:"index"`
}

type Base64Upload struct {
	Filesize int
	Filetype string
	Filename string
	ID       string
	Base64   string
}

func SaveUpload(upload *Base64Upload, dst string) error {
	err:=os.MkdirAll(dst,0777)
	if err != nil {
		return err
	}

	id := bson.NewObjectId().Hex()
	ext := filepath.Ext(upload.Filename)
	out, err := os.Create(filepath.Join(dst, id+ext))
	if err != nil {
		return err
	}

	res, err := base64.StdEncoding.DecodeString(upload.Base64)
	if err != nil {
		if err != nil {
			return err
		}
	}

	reader := bytes.NewReader(res)

	defer out.Close()
	n, err := io.Copy(out, reader)
	if err != nil {
		if err != nil {
			return err
		}
	}
	if int(n) != upload.Filesize {
		return fmt.Errorf("wanted", upload.Filesize, "got", n)
	}

	//clear upload data
	upload.Base64 = ""
	upload.ID = id + ext
	return nil
}

func (pub *User) ID() int { return pub.Id }

func (pub *User) PASSWORD() string { return pub.Password }

func (pub *User) FindUser(email string) (IUser, error) {
	err := stormdb.One("Email", email, pub)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return pub, nil
}

type UserController struct {
}

func (UserController) Create(c *gin.Context) {
	var p = &User{}
	err := c.BindJSON(p)
	if err != nil {
		c.Error(err)
		return
	}

	p.Created = time.Now()
	p.Updated = p.Created
	p.Password = NewSha512Password(p.Password)

	if p.Image != nil {
		err = SaveUpload(p.Image, "uploads")
		if err != nil {
			c.JSON(500, err.Error())
			return
		}

	}

	err = stormdb.Save(p)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
}

func (UserController) Edit(c *gin.Context) {

	idstr := c.Param("id")
	id,err:=strconv.Atoi(idstr) 
	if err != nil {
		c.JSON(500, err)
		return
	}

	var p User
	err = c.BindJSON(&p)
	if err != nil {
		if err != nil {
			c.JSON(500, err.Error())
			return
		}
	}

	var old User

	err = stormdb.One("Id", id, &old)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	if p.Image != nil {
		err = SaveUpload(p.Image, filepath.Join("uploads","users"))
		if err != nil {
			c.JSON(500, err.Error())
			return
		}

	}

	p.Updated = time.Now()
	if p.Password != old.Password {
		p.Password = NewSha512Password(p.Password)
	}

	err = stormdb.Save(p)
	if err != nil {
		c.JSON(500, err)
		return
	}

}

func (UserController) GetId(c *gin.Context) {

	idstr := c.Param("id")
	id,err:=strconv.Atoi(idstr)
	if err != nil {
		c.JSON(500, err)
		return
	}
	var pub User

	err = stormdb.One("Id", id, &pub)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	c.JSON(http.StatusOK, pub)
}

func (UserController) GetImage(c *gin.Context) {

	idstr := c.Param("id")
	id,err:=strconv.Atoi(idstr)
	if err != nil {
		c.JSON(500, err)
		return
	}

	var pub User
	err = stormdb.One("Id", id, &pub)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	c.JSON(http.StatusOK, pub.Image.Base64)

}

func (UserController) Delete(c *gin.Context) {

	idstr := c.Param("id")
	id,err:=strconv.Atoi(idstr)
	if err != nil {
		c.JSON(500, err)
		return
	}
	p := User{Id: id}

	err = stormdb.Remove(p)
	if err != nil {
		c.JSON(500, err)
		return
	}
}

/*
func MustGetUserRole(c *gin.Context) UserRole {
return UserAdmin
}
*/

func (UserController) ListAll(c *gin.Context) {
	var publishers []User

	err := stormdb.All(&publishers)
	if err != nil {
		c.JSON(500, err)
		return
	}

	for i, j := 0, len(publishers)-1; i < j; i, j = i+1, j-1 {
		publishers[i], publishers[j] = publishers[j], publishers[i]
	}

	c.JSON(http.StatusOK, publishers)
}
func (UserController) Search(c *gin.Context) {
	var users []User

	keyword := c.Param("keyword")

	max := 30

	stormdb.Bolt.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("User"))
		c := bucket.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			if strings.Contains(strings.ToLower(string(v)), strings.ToLower(keyword)) {
				var user User

				err := stormdb.Get("User", k, &user)
				if err != nil {
					log.Println(err)
					continue
				}
				users = append(users, user)
			}

		}

		return nil
	})

	if len(users) < max {
		max = len(users)
	}

	c.JSON(200, users[:max])
}

func addDefaultPub() error {
	var p = &User{}
	p.Name = "Marcelo Pires"
	p.Password = "Kirk1zodiak"
	p.Email = "thesyncim@gmail.com"
	p.Role = UserAdmin
	p.Password = NewSha512Password(p.Password)

	//todo validate existing email

	log.Println(stormdb.Save(p))
	return nil

}


