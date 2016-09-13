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

	CreateAdmin()

	//gin.SetMode(gin.ReleaseMode)

	pub := gin.New()

	// middleware
	pub.Use(gin.Logger())
	pub.Use(gin.Recovery())

	// session management
	secret := "wowzaverysecretivelucraciobuddy"
	store := sessions.NewCookieStore([]byte(secret))
	pub.Use(sessions.Sessions(cfg.Site.Title, store))

	// templates
	pub.LoadHTMLGlob(cfg.Server.Template)

	// routes
	users := pub.Group("/", AccessLevelAuth(1))
	editors := pub.Group("/", AccessLevelAuth(2))

	pub.Static("/static", cfg.Server.Static)

	pub.GET("/", PagesHome)
	pub.GET("/img", ImgHome)
	pub.GET("/img/thumb/:id", ImgThumb)
	pub.GET("/img/view/:id", ImgView)
	editors.POST("/img/new", ImgUpload)

	pub.GET("/posts", PostsHome)
	pub.GET("/posts/view/:id", PostsView)
	editors.GET("/posts/new", PostsNew)
	editors.POST("/posts/new", PostsTryNew)

	users.POST("/comments/new", CommentsTryNew)
	users.POST("/comments/reply", CommentsTryReply)

	pub.GET("/auth/sign-in", AuthSignIn)
	pub.POST("/auth/sign-in", AuthTrySignIn)
	pub.GET("/auth/register", AuthRegister)
	pub.POST("/auth/register", AuthTryRegister)

	pub.Run()
	// pub.Run(":3000") for a hard coded port
}
