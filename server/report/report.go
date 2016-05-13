package report

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
)

type Report struct {
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

type ReportController struct {
}

func (ReportController) Create(c *gin.Context) {
	var p = &Report{}
	err := c.BindJSON(p)
	if err != nil {
		c.Error(err)
		return
	}

	sess, ok := c.Get("Session")
	if !ok {
		c.Error(fmt.Errorf("session not found"))
		return
	}
	p.CreatedBy = sess.(user.Session).UserID
	p.Created = time.Now()
	p.Updated = p.Created

	for i := range p.Attachments {
		p.Attachments[i].Id = bson.NewObjectId().Hex()
		err = SaveUpload(p.Attachments[i].File, filepath.Join("uploads", "reports"))
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

func (ReportController) Edit(c *gin.Context) {

	idstr := c.Param("id")
	id,err:=strconv.Atoi(idstr)
	if err != nil {
		c.JSON(500, err)
		return
	}
	var p Report
	err = c.BindJSON(&p)
	if err != nil {
		if err != nil {
			c.JSON(500, err.Error())
			return
		}
	}

	var old Report

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
			err = SaveUpload(p.Attachments[i].File, filepath.Join("uploads", "reports"))
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

func (ReportController) GetId(c *gin.Context) {

	idstr := c.Param("id")
	id,err:=strconv.Atoi(idstr)
	if err != nil {
		c.JSON(500, err)
		return
	}
	var pub Report

	err = stormdb.One("Id", id, &pub)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	c.JSON(http.StatusOK, pub)
}

func (ReportController) Delete(c *gin.Context) {

	idstr := c.Param("id")
	id,err:=strconv.Atoi(idstr)
	if err != nil {
		c.JSON(500, err)
		return
	}
	p := Report{Id: id}

	err = stormdb.Remove(p)
	if err != nil {
		c.JSON(500, err)
		return
	}
}

func (ReportController) ListAll(c *gin.Context) {
	var publishers []Report

	err := stormdb.All(&publishers)
	if err != nil {
		c.JSON(500, err)
		return
	}

	c.JSON(http.StatusOK, publishers)
}

func (ReportController) Search(c *gin.Context) {
	var reports []Report

	keyword := c.Param("keyword")

	max := 30

	stormdb.Bolt.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("Report"))
		c := bucket.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			if strings.Contains(strings.ToLower(string(v)), strings.ToLower(keyword)) {
				var report Report

				err := stormdb.Get("Report", k, &report)
				if err != nil {
					log.Println(err)
					continue
				}
				reports = append(reports, report)
			}

		}

		return nil
	})

	if len(reports) < max {
		max = len(reports)
	}

	c.JSON(200, reports[:max])
}

func addDefaultReport() error {
	var p = &Report{}
	p.Title = "registo xpto"

	log.Println(stormdb.Save(p))
	return nil

}
