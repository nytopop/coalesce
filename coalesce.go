// coalesce/coalesce.go

package main

import (
	"os"
	"runtime"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nytopop/coalesce/controllers"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

//var globalSession, _ = mgo.Dial(os.Getenv("DATABASE_PORT_27017_TCP_ADDR"))
//var dbname = os.Getenv("DB_NAME")

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// TODO: initial config mode if db empty

	gin.SetMode(gin.ReleaseMode)

	pub := gin.New()

	// middleware
	pub.Use(gin.Logger())
	pub.Use(gin.Recovery())

	// session management
	secret := []byte(os.Getenv("SESSION_SECRET"))
	store := sessions.NewCookieStore(secret)
	pub.Use(sessions.Sessions("coalesce", store))

	// authentication
	pub.Use(controllers.AuthCheckpoint())
	commentors := pub.Group("/", controllers.AccessLevelAuth(1))
	editors := pub.Group("/", controllers.AccessLevelAuth(2))
	admins := pub.Group("/", controllers.AccessLevelAuth(3))

	// templates
	pub.LoadHTMLGlob("resources/templates/**/*.html")

	// routes
	pub.Static("/static", "resources/static")
	pub.GET("/", controllers.Home)

	// /img
	editors.GET("/img", controllers.ImgAll)
	pub.GET("/img/thumb/:id", controllers.ImgThumb)
	pub.GET("/img/view/:id", controllers.ImgView)
	editors.GET("/img/new", controllers.ImgNew)
	editors.POST("/img/new", controllers.ImgTryNew)
	editors.GET("/img/del/:id", controllers.ImgTryDelete)

	// /posts
	pub.GET("/posts", controllers.PostsPage)
	pub.GET("/posts/view/:id", controllers.PostsView)
	editors.GET("/posts/new", controllers.PostsNew)
	editors.POST("/posts/new", controllers.PostsTryNew)
	editors.GET("/posts/edit/:id", controllers.PostsEdit)
	editors.POST("/posts/edit/:id", controllers.PostsTryEdit)
	editors.GET("/posts/del/:id", controllers.PostsTryDelete)

	// /comments
	commentors.POST("/comments/new", controllers.CommentsTryNew)
	commentors.POST("/comments/reply", controllers.CommentsTryReply)
	commentors.GET("/comments/del/:id", controllers.CommentsTryDelete)

	// /auth
	pub.GET("/auth/sign-in", controllers.AuthSignIn)
	pub.POST("/auth/sign-in", controllers.AuthTrySignIn)
	pub.GET("/auth/sign-out", controllers.AuthSignOut)
	pub.GET("/auth/register", controllers.AuthRegister)
	pub.POST("/auth/register", controllers.AuthTryRegister)

	// /users
	commentors.GET("/users/me", controllers.UsersMe)
	commentors.GET("/users/myposts", controllers.UsersMyPosts)
	commentors.GET("/users/mycomments", controllers.UsersMyComments)
	admins.GET("/users/all", controllers.UsersAll)
	admins.GET("/users/promote/:id", controllers.UsersTryPromote)
	admins.GET("/users/demote/:id", controllers.UsersTryDemote)
	admins.GET("/users/del/:id", controllers.UsersTryDelete)

	// /config
	admins.GET("/config", controllers.ConfigEdit)
	admins.POST("/config/edit", controllers.ConfigTryEdit)
	admins.POST("/config/reset", controllers.ConfigTryReset)

	// /error
	pub.GET("/error", controllers.ErrorHome)

	pub.Run()
	// pub.Run(":3000") for a hard coded port
}
