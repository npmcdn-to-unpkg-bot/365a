package module

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

	"github.com/gin-gonic/gin"
	"github.com/thesyncim/365/server/user"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	//"github.com/asdine/storm"
	"github.com/blevesearch/bleve"
	"github.com/asdine/storm"
)

type FieldType int

const (
	FieldText FieldType = iota
	FieldNumer
)

type Field struct {
	Name     string
	Value    string
	Required bool
}

//TODO do not edit after any created
type ModuleConfig struct {
	Id                 int `storm:"id"`
	Name               string  `storm:"unique"`
	Description        string
	ExcludeClients     bool
	ExcludeAttach      bool
	ExcludeTitle       bool
	ExcludeDescription bool
	ExcludeWorkers     bool
	Fields             []Field
}

type Module struct {
	Id                int `storm:"id"`
	Title             string `storm:"index"`
	Description       string
	CreatedBy         int `storm:"index"`
	ClientID          int `storm:"index"`
	Users             []user.User `storm:"index"`
	Client            *user.User
	Attachments       []*struct {
		Id   string
		Tags string
		File *Base64Upload
	}
	ModuleConfig      *ModuleConfig `storm:"index"`
	ExtraFields       []Field
	ExtraFieldsHeader []Field
	Created           time.Time `storm:"index"`
	Updated           time.Time `storm:"index"`
}
type Base64Upload struct {
	Filesize int
	Filetype string
	Filename string
	ID       string
	Base64   string
}

func getExtraFields(id string) ([]Field, error) {
	var moduleInfo ModuleConfig
	err := stormdb.One("Id", id, &moduleInfo)
	return moduleInfo.Fields, err
}

func SaveUpload(upload *Base64Upload, dst string) error {
	err := os.MkdirAll(dst, 0777)
	if err != nil {
		return err
	}

	id := bson.NewObjectId().Hex()
	ext := filepath.Ext(upload.Filename)
	out, err := os.Create(filepath.Join(dst, id + ext))
	if err != nil {
		return err
	}

	res, err := base64.StdEncoding.DecodeString(upload.Base64)
	if err != nil {

		return err

	}

	reader := bytes.NewReader(res)

	defer out.Close()
	n, err := io.Copy(out, reader)
	if err != nil {
		return err
	}
	if int(n) != upload.Filesize {
		return fmt.Errorf("wanted", upload.Filesize, "got", n)
	}

	//clear upload data
	upload.Base64 = ""
	upload.ID = id + ext
	return nil
}

type ModuleController struct {
}

func (ModuleController) Create(c *gin.Context) {

	modulestr := c.Param("module")
	moduleid, err := strconv.Atoi(modulestr)
	if err != nil {
		c.JSON(500, err)
		return
	}

	var p = &Module{}
	err = c.BindJSON(p)
	if err != nil {
		c.Error(err)
		return
	}

	sess, ok := c.Get("Session")
	if !ok {
		c.Error(fmt.Errorf("session not found"))
		c.JSON(500, "session not found")
		return
	}
	var moduleinfo ModuleConfig
	err = stormdb.One("Id", modulestr, &moduleinfo)
	if err != nil {
		c.Error(err)
		return
	}
	p.ModuleConfig = &moduleinfo
	p.CreatedBy = sess.(user.Session).UserID
	p.Created = time.Now()
	p.Updated = p.Created

	for i := range p.Attachments {
		p.Attachments[i].Id = bson.NewObjectId().Hex()
		err = SaveUpload(p.Attachments[i].File, filepath.Join("uploads", "modules", fmt.Sprint(moduleid)))
		if err != nil {
			c.Error(err)
			return
		}
	}

	err = stormdb.Save(p)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	index := GetBleveIndex("365.module." + modulestr)
	err = index.Index(fmt.Sprint(p.Id), p)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}

}

func (ModuleController) Edit(c *gin.Context) {
	modulestr := c.Param("module")

	idstr := c.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		c.JSON(500, err)
		return
	}
	var p Module
	err = c.BindJSON(&p)
	if err != nil {
		if err != nil {
			c.JSON(500, err.Error())
			return
		}
	}

	var old Module

	sess, ok := c.Get("Session")
	if !ok {
		c.Error(fmt.Errorf("session not found"))
		return
	}
	p.CreatedBy = sess.(user.Session).UserID

	err = stormdb.One("Id", id, &old)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	p.Updated = time.Now()

	for i := range p.Attachments {
		//if we provide a a file replace it
		if p.Attachments[i].File.Base64 != "" {
			p.Attachments[i].Id = bson.NewObjectId().Hex()
			err = SaveUpload(p.Attachments[i].File, filepath.Join("uploads", "modules", fmt.Sprint(modulestr)))
			if err != nil {
				c.Error(err)
				return
			}
		}

	}

	err = stormdb.Save(p)
	if err != nil {
		c.JSON(500, err)
		return
	}
	index := GetBleveIndex("365.module." + modulestr)

	err = index.Index(fmt.Sprint(p.Id), p)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}

}

func (ModuleController) GetId(c *gin.Context) {

	idstr := c.Param("id")

	var pub Module

	err := stormdb.One("Id", idstr, &pub)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	c.JSON(http.StatusOK, pub)
}

func (ModuleController) GetModuleSchema(c *gin.Context) {

	idstr := c.Param("module")

	var pub Module

	var moduleinfo ModuleConfig
	err := stormdb.One("Id", idstr, &moduleinfo)
	if err != nil {
		c.JSON(500, err)
		return
	}

	fields, err := getExtraFields(idstr)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	pub.ModuleConfig = &moduleinfo

	pub.ExtraFieldsHeader = fields

	c.JSON(http.StatusOK, pub)
}

func (ModuleController) CreateModule(c *gin.Context) {
	var p = &ModuleConfig{}
	err := c.BindJSON(p)
	if err != nil {
		c.Error(err)
		return
	}

	err = stormdb.Save(p)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
}

func (ModuleController) EditModule(c *gin.Context) {

	var p ModuleConfig
	err := c.BindJSON(&p)
	if err != nil {
		if err != nil {
			c.JSON(500, err.Error())
			return
		}
	}

	var old ModuleConfig

	err = stormdb.One("Id", p.Id, &old)
	if err != nil {
		c.JSON(500, err.Error())
		c.Error(err)
		return
	}
	var modules []Module

	err = stormdb.Find("ModuleConfig", &old, &modules)
	if err != nil && err!=storm.ErrNotFound {
		c.JSON(500, err.Error())
		c.Error(err)
		return
	}
	tx, err := stormdb.Begin(true)
	if err != nil {
		c.JSON(500, err.Error())
		c.Error(err)
		return
	}

	for i := range modules {
		modules[i].ModuleConfig = &p
		err = tx.Save(modules[i])
		if err != nil {
			tx.Rollback()
			c.JSON(500, err.Error())
			c.Error(err)
			return
		}
	}

	err = tx.Save(&p)
	if err != nil {
		tx.Rollback()
		c.JSON(500, err.Error())
		c.Error(err)
		return
	}

	err = tx.Commit()
	if err != nil {
		c.JSON(500, err.Error())
		c.Error(err)
		return
	}

	//update all modules
	//var moduleentries []Module


}

//TODO delete images an attach on delete
func (ModuleController) DeleteModule(c *gin.Context) {

	idstr := c.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		c.JSON(500, err)
		return
	}
	p := ModuleConfig{Id: id}

	err = stormdb.One("Id", p.Id, &p)
	if err != nil {
		c.JSON(500, err.Error())
		c.Error(err)
		return
	}

	var modules []Module
	err = stormdb.Find("ModuleConfig", &p, &modules)
	if err != nil && err!=storm.ErrNotFound{
		c.JSON(500, err.Error())
		c.Error(err)
		return
	}
	tx, err := stormdb.Begin(true)
	for i := range modules {
		err = tx.Remove(modules[i])
		if err != nil {
			tx.Rollback()
			c.JSON(500, err.Error())
			c.Error(err)
			return
		}
	}

	err = tx.Remove(p)
	if err != nil {
		tx.Rollback()
		c.JSON(500, err.Error())
		c.Error(err)
		return
	}
	err = tx.Commit()
	if err != nil {
		c.JSON(500, err.Error())
		c.Error(err)
		return
	}

	err = os.RemoveAll(filepath.Join("uploads", "modules", idstr))
	if err != nil {
		c.JSON(500, err.Error())
		c.Error(err)
	}
}

func (ModuleController) GetIdModule(c *gin.Context) {

	idstr := c.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		c.JSON(500, err)
		return
	}
	var pub ModuleConfig

	err = stormdb.One("Id", id, &pub)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	c.JSON(http.StatusOK, pub)
}

func (ModuleController) ListAllModule(c *gin.Context) {
	var moduleEntry []ModuleConfig

	err := stormdb.All(&moduleEntry)
	if err != nil {
		c.JSON(500, err)
		return
	}

	for i, j := 0, len(moduleEntry) - 1; i < j; i, j = i + 1, j - 1 {
		moduleEntry[i], moduleEntry[j] = moduleEntry[j], moduleEntry[i]
	}

	c.JSON(http.StatusOK, moduleEntry)
}

func (ModuleController) Delete(c *gin.Context) {

	idstr := c.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		c.JSON(500, err)
		return
	}

	module := c.Param("module")

	p := Module{Id: id}

	err = stormdb.Remove(p)
	if err != nil {
		c.JSON(500, err)
		return
	}

	index := GetBleveIndex("365.module." + fmt.Sprint(module))
	err = index.Delete(idstr)
	if err != nil {
		c.JSON(500, err)
		return
	}
}

func (ModuleController) ListAll(c *gin.Context) {
	var moduleEntry []Module

	modulestr := c.Param("module")
	moduleid, err := strconv.Atoi(modulestr)
	if err != nil {
		c.JSON(500, err)
		return
	}
	var modulec ModuleConfig
	stormdb.One("Id", moduleid, &modulec)

	err = stormdb.Find("ModuleConfig", &modulec, &moduleEntry)
	if err != nil && err!=storm.ErrNotFound {
		log.Println(err)
		c.JSON(500, err)
		return
	}

	for i, j := 0, len(moduleEntry) - 1; i < j; i, j = i + 1, j - 1 {
		moduleEntry[i], moduleEntry[j] = moduleEntry[j], moduleEntry[i]
	}
	c.JSON(http.StatusOK, moduleEntry)
}

func (ModuleController) Search(c *gin.Context) {

	modulestr := c.Param("module")
	moduleid, err := strconv.Atoi(modulestr)
	if err != nil {
		c.JSON(500, err)
		return
	}

	var modulec ModuleConfig

	err = stormdb.One("Id", moduleid, &modulec)
	if err != nil {
		c.JSON(500, err)
		return
	}

	keyword := c.Param("keyword")

	index := GetBleveIndex("365.module." + modulestr)
	query := bleve.NewQueryStringQuery(keyword)
	searchRequest := bleve.NewSearchRequest(query)
	res, err := index.Search(searchRequest)
	if err != nil {
		log.Println(err)
		c.JSON(500, err.Error())
		return
	}

	var results []Module
	for i := range res.Hits {
		var result Module
		err = stormdb.One("Id", res.Hits[i].ID, &result)
		if err != nil {
			c.JSON(500, err)
			return
		}
		results = append(results, result)
	}

	c.JSON(200, results)
}

func addDefaultModule() error {
	var p = &Module{}
	p.Title = "registo xpto"

	log.Println(stormdb.Save(p))
	return nil

}
