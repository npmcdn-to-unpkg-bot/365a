//file sstream_plugin.go
package client

import (

	"github.com/gin-gonic/gin"

	"github.com/asdine/storm"
)

type ClientPlugin struct {
	BaseUrl          string
	Authenticator    gin.HandlerFunc
	ClientController *ClientController
}

func NewUserPlugin(baseURL string, authenticator gin.HandlerFunc, db *storm.DB) *ClientPlugin {
	stormdb = db
	stormdb.Init(Client{})

	return &ClientPlugin{
		BaseUrl:          baseURL,
		Authenticator:    authenticator,
		ClientController: &ClientController{},
	}
}

func (i ClientPlugin) GetName() string {
	return "Client Plugin"
}

func (i ClientPlugin) GetDescription() string {
	return ""
}

func (i *ClientPlugin) Register(s *gin.Engine) {
	pub := s.Group(i.BaseUrl)
	pub.Use(i.Authenticator)
	pub.POST("/create", i.ClientController.Create)
	pub.POST("/edit/:id", i.ClientController.Edit)
	pub.GET("/getid/:id", i.ClientController.GetId)
	pub.POST("/delete/:id", i.ClientController.Delete)
	pub.GET("/listall", i.ClientController.ListAll)
	pub.GET("/search/:keyword", i.ClientController.Search)
}
