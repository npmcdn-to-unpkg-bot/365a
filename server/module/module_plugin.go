//file sstream_plugin.go
package module

import (
	"github.com/gin-gonic/gin"
	"github.com/asdine/storm"
	"log"

	"github.com/blevesearch/bleve"
	"sync"
)

var (
	stormdb *storm.DB
	index bleve.Index
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

func openIndex(path string) bleve.Index {
	index, err := bleve.Open(path)
	if err == bleve.ErrorIndexPathDoesNotExist {
		log.Printf("Creating new index...")
		// create a mapping
		indexMapping := bleve.NewIndexMapping()
		index, err = bleve.New(path, indexMapping)
		if err != nil {
			log.Fatal(err)
		}
	} else if err == nil {
		log.Printf("Opening existing index...")
	} else {
		log.Fatal(err)
	}
	return index
}

var instance =map[string]*Bleveinstance{}

type Bleveinstance struct {
	Once  sync.Once
	Index bleve.Index
}

func GetBleveIndex(id string) bleve.Index {
	if _, ok := instance[id]; !ok {
		instance[id]=new(Bleveinstance)
	}
	instance[id].Once.Do(func() {
		instance[id].Index = openIndex(id)
	})
	return instance[id].Index
}