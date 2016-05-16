package client

import (
	"time"

	"github.com/gin-gonic/gin"

	"log"
	"net/http"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/thesyncim/365/server/user"
	"strconv"
)

type Client struct {
	Id      int `storm:"id"`
	Name    string
	Nif     string `storm:"unique"`
	Created time.Time
	Updated time.Time
}

func MustGetUserRole(c *gin.Context) user.UserRole {
	s, err := user.ReadSession(c)
	if err != nil {
		panic(err)
	}
	var user user.User
	err = stormdb.One("Id", s.UserID, &user)
	if err != nil {
		panic(err)
	}

	return user.Role
}

type ClientController struct {
}

func (ClientController) Create(c *gin.Context) {
	var p = &Client{}
	err := c.BindJSON(p)
	if err != nil {
		c.Error(err)
		return
	}

	log.Println(MustGetUserRole(c))

	p.Created = time.Now()

	err = stormdb.Save(p)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
}

func (ClientController) Edit(c *gin.Context) {

	idstr := c.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		c.JSON(500, err)
		return
	}
	var p Client
	err = c.BindJSON(&p)
	if err != nil {
		if err != nil {
			c.JSON(500, err.Error())
			return
		}
	}

	var old Client

	err = stormdb.One("Id", id, &old)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	p.Updated = time.Now()

	err = stormdb.Save(p)
	if err != nil {
		c.JSON(500, err)
		return
	}

}

func (ClientController) GetId(c *gin.Context) {

	idstr := c.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		c.JSON(500, err)
		return
	}
	var pub Client

	err = stormdb.One("Id", id, &pub)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	c.JSON(http.StatusOK, pub)
}

func (ClientController) Delete(c *gin.Context) {

	idstr := c.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		c.JSON(500, err)
		return
	}
	p := Client{Id: id}

	err = stormdb.Remove(p)
	if err != nil {
		c.JSON(500, err)
		return
	}
}

func (ClientController) ListAll(c *gin.Context) {
	var clients[]Client

	err := stormdb.All(&clients)
	if err != nil {
		c.JSON(500, err)
		return
	}
	for i, j := 0, len(clients)-1; i < j; i, j = i+1, j-1 {
		clients[i], clients[j] = clients[j], clients[i]
	}

	c.JSON(http.StatusOK, clients)
}
func (ClientController) Search(c *gin.Context) {
	var clients []Client

	keyword := c.Param("keyword")

	max := 30

	stormdb.Bolt.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("Client"))
		c := bucket.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			if strings.Contains(strings.ToLower(string(v)), strings.ToLower(keyword)) {
				var client Client

				err := stormdb.Get("Client", k, &client)
				if err != nil {
					log.Println(err)
					continue
				}
				clients = append(clients, client)
			}

		}

		return nil
	})

	if len(clients) < max {
		max = len(clients)
	}

	c.JSON(200, clients[:max])
}

