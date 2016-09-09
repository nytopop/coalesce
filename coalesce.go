// coalesce/coalesce.go

package main

import (
	"runtime"

	gcfg "gopkg.in/gcfg.v1"
	mgo "gopkg.in/mgo.v2"

	"github.com/gin-gonic/gin"
)

var cfg = Config{}
var err = gcfg.ReadFileInto(&cfg, "coalesce.cfg")
var globalSession, _ = mgo.Dial(cfg.Database.Host)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	//gin.SetMode(gin.ReleaseMode)

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	r := gin.Default()

	// templates
	r.LoadHTMLGlob(cfg.Server.Template)

	// routes
	r.Static("/static", cfg.Server.Static)
	r.GET("/img", ImgHome)
	r.GET("/img/thumb/:id", ImgThumb)
	r.GET("/img/view/:id", ImgView)
	r.POST("/img/new", ImgUpload)

	r.GET("/posts", PostsHome)
	r.GET("/posts/view/:id", PostsView)
	r.GET("/posts/new", PostsNew)
	r.POST("/posts/new", PostsCreate)

	r.Run()
	// r.Run(":3000") for a hard coded port
}
