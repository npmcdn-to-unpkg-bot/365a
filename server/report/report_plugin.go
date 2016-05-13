//file sstream_plugin.go
package report

import (

	"github.com/gin-gonic/gin"
	"github.com/asdine/storm"
)

type ReportPlugin struct {
	BaseUrl          string
	Authenticator    gin.HandlerFunc
	ReportController *ReportController
}

func NewReportPlugin(baseURL string, authenticator gin.HandlerFunc, db *storm.DB) *ReportPlugin {
	stormdb = db
	stormdb.Init(Report{})

	return &ReportPlugin{
		BaseUrl:          baseURL,
		Authenticator:    authenticator,
		ReportController: &ReportController{},
	}
}

func (i ReportPlugin) GetName() string {
	return "Client Plugin"
}

func (i ReportPlugin) GetDescription() string {
	return ""
}

func (i *ReportPlugin) Register(s *gin.Engine) {
	pub := s.Group(i.BaseUrl)
	pub.Use(i.Authenticator)

	pub.POST("/create", i.ReportController.Create)
	pub.POST("/edit/:id", i.ReportController.Edit)
	pub.GET("/getid/:id", i.ReportController.GetId)
	pub.POST("/delete/:id", i.ReportController.Delete)
	pub.GET("/listall", i.ReportController.ListAll)
	pub.GET("/search/:keyword", i.ReportController.Search)
}
