//file sstream_plugin.go
package user

import (

	"github.com/gin-gonic/gin"
	"github.com/asdine/storm"
)

type UserPlugin struct {
	BaseUrl        string
	Authenticator  gin.HandlerFunc
	UserController *UserController
}

func NewUserPlugin(baseURL string, authenticator gin.HandlerFunc, db *storm.DB) *UserPlugin {
	stormdb = db
	stormdb.Init(User{})
	stormdb.Init(Session{})
	addDefaultPub()

	return &UserPlugin{
		BaseUrl:        baseURL,
		Authenticator:  authenticator,
		UserController: &UserController{},
	}
}

func (i UserPlugin) GetName() string {
	return "Opinion Publisher"
}

func (i UserPlugin) GetDescription() string {
	return "Azorestv Opinion App Publisher"
}

func (i *UserPlugin) Register(s *gin.Engine) {
	pub := s.Group(i.BaseUrl)
	pub.Use(i.Authenticator)

	pub.POST("/create", i.UserController.Create)
	pub.POST("/edit/:id", i.UserController.Edit)
	pub.GET("/getid/:id", i.UserController.GetId)
	pub.POST("/delete/:id", i.UserController.Delete)
	pub.GET("/listall", i.UserController.ListAll)
	pub.GET("/publisher/image/:id", i.UserController.GetImage)
	pub.GET("/search/:keyword", i.UserController.Search)

}
