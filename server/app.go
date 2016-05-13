package main

import (
	"github.com/gin-gonic/gin"

	"github.com/kardianos/service"
	"github.com/thesyncim/365/server/client"
	"github.com/thesyncim/365/server/report"
	"github.com/thesyncim/365/server/module"
	"github.com/thesyncim/365/server/user"
	"time"

	"github.com/asdine/storm"
	"io"

	"github.com/braintree/manners"

	"github.com/itsjamie/gin-cors"
	"github.com/asdine/storm/codec/json"
	"net/http"
	"github.com/elazarl/go-bindata-assetfs"
	"strings"
	"path/filepath"
	"fmt"
	"os"
)

type app struct {
	Logfile io.ReadWriteCloser
}

func (a *app) run() error {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	fmt.Println(dir)

	os.MkdirAll(filepath.Join("uploads","reports"),0755)
	os.MkdirAll("uploads",0755)


	db, err :=storm.Open("my.db",storm.AutoIncrement(),storm.Codec(json.Codec))
	if err != nil {
		return err
	}



	authenticator := user.AngularAuth(db)
	_=authenticator

	api := gin.Default()
	api.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))

	userplugin := user.NewUserPlugin("/users", authenticator, db)

	userplugin.Register(api)

	clientsPlugin := client.NewUserPlugin("/clients", authenticator, db)
	clientsPlugin.Register(api)

	reportsPlugin := report.NewReportPlugin("/reports", authenticator, db)
	reportsPlugin.Register(api)

	modulesPlugin := module.NewModulePlugin("/modules", authenticator, db)
	modulesPlugin.Register(api)

	api.POST("/auth/login", user.AngularSignIn(db, (&user.User{}).FindUser, user.NewSha512Password, time.Hour*48))
	api.StaticFS("/static", BinaryFileSystem(""))
	api.GET("/", func(c *gin.Context) {
		c.Redirect(302,"/static/index.html")
	})
	api.GET("/index.html", func(c *gin.Context) {
		c.Redirect(302,"/static/index.html")
	})

	api.Static("/uploads","uploads")



	return manners.ListenAndServe(":7894", api)
}

func (a *app) Start(s service.Service) error {

	go a.run()
	return nil
}

func (a *app) Stop(s service.Service) error {
	manners.Close()
	a.Logfile.Close()
	return nil
}

type binaryFileSystem struct {
	fs http.FileSystem
}

func (b *binaryFileSystem) Open(name string) (http.File, error) {
	return b.fs.Open(name)
}

func (b *binaryFileSystem) Exists(prefix string, filepath string) bool {

	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		if _, err := b.fs.Open(p); err != nil {
			return false
		}
		return true
	}
	return false
}

func BinaryFileSystem(root string) *binaryFileSystem {
	fs := &assetfs.AssetFS{Asset, AssetDir,AssetInfo ,root}
	return &binaryFileSystem{
		fs,
	}
}
