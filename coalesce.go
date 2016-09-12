// coalesce/coalesce.go

package main

import (
	"runtime"

	gcfg "gopkg.in/gcfg.v1"
	mgo "gopkg.in/mgo.v2"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

var cfg = Config{}
var err = gcfg.ReadFileInto(&cfg, "coalesce.cfg")
var globalSession, _ = mgo.Dial(cfg.Database.Host)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	//gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	// middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// session management
	secret := "wowzaverysecretivelucraciobuddy"
	store := sessions.NewCookieStore([]byte(secret))
	r.Use(sessions.Sessions(cfg.Site.Title, store))

	// templates
	r.LoadHTMLGlob(cfg.Server.Template)

	// routes
	r.Static("/static", cfg.Server.Static)
	r.GET("/", PagesHome)
	r.GET("/img", ImgHome)
	r.GET("/img/thumb/:id", ImgThumb)
	r.GET("/img/view/:id", ImgView)

	r.POST("/img/new", ImgUpload)

	r.GET("/posts", PostsHome)
	r.GET("/posts/view/:id", PostsView)
	authorized := r.Group("/posts/new", AnyUserAuth())
	authorized.GET("/", PostsNew)
	authorized.POST("/", PostsTryNew)

	r.GET("/auth/sign-in", AuthSignIn)
	r.GET("/auth/register", AuthRegister)
	r.POST("/auth/sign-in", AuthTrySignIn)
	r.POST("/auth/register", AuthTryRegister)

	r.Run()
	// r.Run(":3000") for a hard coded port
}
