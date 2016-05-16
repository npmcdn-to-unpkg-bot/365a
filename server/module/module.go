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
	"strings"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"github.com/thesyncim/365/server/client"
	"github.com/thesyncim/365/server/user"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	//"github.com/asdine/storm"
)

type FieldType int

const (
	FieldText FieldType = iota
	FieldNumer



)
type Field struct{
	Name string
	Value string
}

//TODO do not edit after any created
type ModuleConfig struct {
	Id int `storm:"id"`
	Name string  `storm:"unique"`
	Description string
	ExcludeClients bool
	ExcludeAttach bool
	ExcludeTitle bool
	ExcludeDescription bool
	ExcludeWorkers bool
	Fields []Field
}

type Module struct {
	Id          int `storm:"id"`
	Title       string `storm:"index"`
	Description string
	CreatedBy   int
	Client      *client.Client `storm:"index"`
	ClientID    int
	Users       []user.User `storm:"index"`
	Attachments []*struct {
		Id   string
		Tags string
		File *Base64Upload
	}
	ModuleConfig *ModuleConfig `storm:"index"`
	ExtraFields []Field
	ExtraFieldsHeader []Field
	Created time.Time `storm:"index"`
	Updated time.Time `storm:"index"`
}
type Base64Upload struct {
	Filesize int
	Filetype string
	Filename string
	ID       string
	Base64   string
}


func getExtraFields(id int)([]Field,error){
	var moduleInfo ModuleConfig
	err:=stormdb.One("Id",id,&moduleInfo)
	return moduleInfo.Fields,err
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
	moduleid,err:=strconv.Atoi(modulestr)
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
		return
	}
	var moduleinfo ModuleConfig
	err=stormdb.One("Id",moduleid,&moduleinfo)
	if err != nil {
		c.Error(err)
		return
	}
	p.ModuleConfig=&moduleinfo
	p.CreatedBy = sess.(user.Session).UserID
	p.Created = time.Now()
	p.Updated = p.Created

	for i := range p.Attachments {
		p.Attachments[i].Id = bson.NewObjectId().Hex()
		err = SaveUpload(p.Attachments[i].File,  filepath.Join("uploads", "modules",fmt.Sprint(moduleid)))
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

}

func (ModuleController) Edit(c *gin.Context) {
	modulestr := c.Param("module")
	moduleid,err:=strconv.Atoi(modulestr)
	if err != nil {
		c.JSON(500, err)
		return
	}


	idstr := c.Param("id")
	id,err:=strconv.Atoi(idstr)
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
			err = SaveUpload(p.Attachments[i].File, filepath.Join("uploads", "modules",fmt.Sprint(moduleid)))
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

}

func (ModuleController) GetId(c *gin.Context) {

	idstr := c.Param("id")
	id,err:=strconv.Atoi(idstr)
	if err != nil {
		c.JSON(500, err)
		return
	}
	var pub Module

	err = stormdb.One("Id", id, &pub)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}





	c.JSON(http.StatusOK, pub)
}

func (ModuleController) GetModuleSchema(c *gin.Context) {

	idstr := c.Param("module")
	id,err:=strconv.Atoi(idstr)
	if err != nil {
		c.JSON(500, err)
		return
	}
	var pub Module

	var moduleinfo ModuleConfig
	err=stormdb.One("Id",idstr,&moduleinfo)
	if err != nil {
		c.JSON(500, err)
		return
	}

	fields,err:=getExtraFields(id)
	if err != nil   {
		c.JSON(500, err.Error())
		return
	}
	pub.ModuleConfig=&moduleinfo

	pub.ExtraFieldsHeader=fields

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

	err=stormdb.Save(&p)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}

}

//TODO delete images an attach on delete
func (ModuleController) DeleteModule(c *gin.Context) {

	idstr := c.Param("id")
	id,err:=strconv.Atoi(idstr)
	if err != nil {
		c.JSON(500, err)
		return
	}
	p := ModuleConfig{Id: id}





	err = stormdb.Remove(p)
	if err != nil {
		c.JSON(500, err)
		return
	}
	os.RemoveAll(filepath.Join("uploads","modules",idstr))
}


func (ModuleController) GetIdModule(c *gin.Context) {

	idstr := c.Param("id")
	id,err:=strconv.Atoi(idstr)
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

	for i, j := 0, len(moduleEntry)-1; i < j; i, j = i+1, j-1 {
		moduleEntry[i], moduleEntry[j] = moduleEntry[j], moduleEntry[i]
	}

	c.JSON(http.StatusOK, moduleEntry)
}



func (ModuleController) Delete(c *gin.Context) {

	idstr := c.Param("id")
	id,err:=strconv.Atoi(idstr)
	if err != nil {
		c.JSON(500, err)
		return
	}
	p := Module{Id: id}

	err = stormdb.Remove(p)
	if err != nil {
		c.JSON(500, err)
		return
	}
}




func (ModuleController) ListAll(c *gin.Context) {
	var moduleEntry []Module

	modulestr := c.Param("module")
	moduleid,err:=strconv.Atoi(modulestr)
	if err != nil {
		c.JSON(500, err)
		return
	}

	err = stormdb.Find("ModuleConfig",moduleid,&moduleEntry)
	if err != nil {
		log.Println(err)
		c.JSON(500, err)
		return
	}

	for i, j := 0, len(moduleEntry)-1; i < j; i, j = i+1, j-1 {
		moduleEntry[i], moduleEntry[j] = moduleEntry[j], moduleEntry[i]
	}
	c.JSON(http.StatusOK, moduleEntry)
}

func (ModuleController) Search(c *gin.Context) {
	var Modules []Module

	modulestr := c.Param("module")
	moduleid,err:=strconv.Atoi(modulestr)
	if err != nil {
		c.JSON(500, err)
		return
	}

	keyword := c.Param("keyword")

	max := 30

	stormdb.Bolt.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("Module"))
		c := bucket.Cursor()

		var Module Module

		for k, v := c.First(); k != nil; k, v = c.Next() {
			if strings.Contains(strings.ToLower(string(v)), strings.ToLower(keyword)) {
				err := stormdb.Get("Module", k, &Module)
				if err != nil {
					log.Println(err)
					continue
				}
				if Module.ModuleConfig.Id!=moduleid{
					continue
				}
				Modules = append(Modules, Module)
			}
		}

		return nil
	})

	if len(Modules) < max {
		max = len(Modules)
	}

	c.JSON(200, Modules[:max])
}

func addDefaultModule() error {
	var p = &Module{}
	p.Title = "registo xpto"

	log.Println(stormdb.Save(p))
	return nil

}
