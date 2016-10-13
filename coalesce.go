// coalesce/coalesce.go

package main

import (
	"os"
	"runtime"

	gcfg "gopkg.in/gcfg.v1"
	mgo "gopkg.in/mgo.v2"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

var cfg = Config{}
var err = gcfg.ReadFileInto(&cfg, "coalesce.cfg")

var globalSession, _ = mgo.Dial(os.Getenv("DATABASE_PORT_27017_TCP_ADDR"))
var dbname = os.Getenv("DB_NAME")

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// BUG: panics if no database connection
	CreateAdmin()

	gin.SetMode(gin.ReleaseMode)

	pub := gin.New()

	// middleware
	pub.Use(gin.Logger())
	pub.Use(gin.Recovery())

	// session management
	secret := "wowzaverysecretivelucraciobuddy"
	store := sessions.NewCookieStore([]byte(secret))
	pub.Use(sessions.Sessions(cfg.Site.Title, store))

	// authentication
	pub.Use(AuthCheckpoint())
	commentors := pub.Group("/", AccessLevelAuth(1))
	editors := pub.Group("/", AccessLevelAuth(2))
	admins := pub.Group("/", AccessLevelAuth(3))

	// templates
	pub.LoadHTMLGlob("resources/templates/**/*.html")

	// routes
	pub.Static("/static", "resources/static")
	pub.GET("/", PagesHome)

	// /img
	editors.GET("/img", ImgAll)
	pub.GET("/img/thumb/:id", ImgThumb)
	pub.GET("/img/view/:id", ImgView)
	editors.GET("/img/new", ImgNew)
	editors.POST("/img/new", ImgTryNew)

	// /posts
	pub.GET("/posts", PostsAll)
	pub.GET("/posts/view/:id", PostsView)
	commentors.GET("/posts/me", PostsMe)
	editors.GET("/posts/new", PostsNew)
	editors.POST("/posts/new", PostsTryNew)
	editors.GET("/posts/edit/:id", PostsEdit)
	editors.POST("/posts/edit", PostsTryEdit)
	editors.GET("/posts/del/:id", PostsTryDelete)

	// /comments
	pub.POST("/comments/new", CommentsTryNew)
	//pub.POST("/comments/reply", CommentsTryReply)

	// /auth
	pub.GET("/auth/sign-in", AuthSignIn)
	pub.POST("/auth/sign-in", AuthTrySignIn)
	pub.GET("/auth/sign-out", AuthSignOut)
	pub.GET("/auth/register", AuthRegister)
	pub.POST("/auth/register", AuthTryRegister)

	// /users
	admins.GET("/users/all", UsersAll)
	admins.GET("/users/promote/:name", UsersTryPromote)
	admins.GET("/users/demote/:name", UsersTryDemote)
	admins.GET("/users/del/:name", UsersTryDelete)

	// /config
	admins.GET("/config", ConfigEdit)
	admins.POST("/config/edit", ConfigTryEdit)

	// /error
	pub.GET("/error", ErrorHome)

	pub.Run()
	// pub.Run(":3000") for a hard coded port
}
