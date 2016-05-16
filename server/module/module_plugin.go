//file sstream_plugin.go
package module

import (

	"github.com/gin-gonic/gin"
	"github.com/asdine/storm"
)

type ModulePlugin struct {
	BaseUrl          string
	Authenticator    gin.HandlerFunc
	ModuleController *ModuleController
}

func NewModulePlugin(baseURL string, authenticator gin.HandlerFunc, db *storm.DB) *ModulePlugin {
	stormdb = db
	stormdb.Init(Module{})
	stormdb.Init(ModuleConfig{})

	return &ModulePlugin{
		BaseUrl:          baseURL,
		Authenticator:    authenticator,
		ModuleController: &ModuleController{},
	}
}

func (i ModulePlugin) GetName() string {
	return "Client Plugin"
}

func (i ModulePlugin) GetDescription() string {
	return ""
}

func (i *ModulePlugin) Register(s *gin.Engine) {
	pub := s.Group(i.BaseUrl)
	pub.Use(i.Authenticator)


	pub.POST("/createmodule/", i.ModuleController.CreateModule)
	pub.POST("/editmodule/:id", i.ModuleController.EditModule)
	pub.POST("/deletemodule/:id", i.ModuleController.DeleteModule)
	pub.GET("/getidmodule/:id", i.ModuleController.GetIdModule)
	pub.GET("/listallmodule", i.ModuleController.ListAllModule)
	pub.GET("/getmoduleschema/:module", i.ModuleController.GetModuleSchema)


	pub.POST("/create/:module", i.ModuleController.Create)
	pub.POST("/edit/:module/:id", i.ModuleController.Edit)
	pub.GET("/getid/:module/:id", i.ModuleController.GetId)
	pub.POST("/delete/:module/:id", i.ModuleController.Delete)
	pub.GET("/listall/:module", i.ModuleController.ListAll)
	pub.GET("/search/:module/:keyword", i.ModuleController.Search)
}
