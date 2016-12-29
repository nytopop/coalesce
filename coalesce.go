// coalesce/coalesce.go

package main

import (
	"log"
	"os"
	"runtime"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

//var globalSession, _ = mgo.Dial(os.Getenv("DATABASE_PORT_27017_TCP_ADDR"))
//var dbname = os.Getenv("DB_NAME")

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var err error
	sqdb, err = sql.Open("sqlite3", "coalesce.db")
	if err != nil {
		log.Fatal(err)
	}
	defer sqdb.Close()

	// BUG: panics if no database connection
	//CreateAdmin()

	if err := initDB(); err != nil {
		log.Fatal(err)
	}

	// TODO: initial config mode if db empty

	gin.SetMode(gin.ReleaseMode)

	pub := gin.New()

	// middleware
	pub.Use(gin.Logger())
	pub.Use(gin.Recovery())

	// session management
	//secret := "wowzaverysecretivelucraciobuddy"
	secret := []byte(os.Getenv("SESSION_SECRET"))
	store := sessions.NewCookieStore(secret)
	pub.Use(sessions.Sessions("coalesce", store))

	// authentication
	pub.Use(AuthCheckpoint())
	commentors := pub.Group("/", AccessLevelAuth(1))
	editors := pub.Group("/", AccessLevelAuth(2))
	admins := pub.Group("/", AccessLevelAuth(3))

	// templates
	pub.LoadHTMLGlob("resources/templates/**/*.html")

	// routes
	pub.Static("/static", "resources/static")
	pub.GET("/", Home)

	// /img
	editors.GET("/img", ImgAll)
	pub.GET("/img/thumb/:id", ImgThumb)
	pub.GET("/img/view/:id", ImgView)
	editors.GET("/img/new", ImgNew)
	editors.POST("/img/new", ImgTryNew)
	editors.GET("/img/del/:id", ImgTryDelete)

	// /posts
	pub.GET("/posts", PostsPage)
	pub.GET("/posts/view/:id", PostsView)
	editors.GET("/posts/new", PostsNew)
	editors.POST("/posts/new", PostsTryNew)
	editors.GET("/posts/edit/:id", PostsEdit)
	editors.POST("/posts/edit/:id", PostsTryEdit)
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
	commentors.GET("/users/me", UsersMe)
	commentors.GET("/users/myposts", UsersMyPosts)
	commentors.GET("/users/mycomments", UsersMyComments)
	admins.GET("/users/all", UsersAll)
	admins.GET("/users/promote/:id", UsersTryPromote)
	admins.GET("/users/demote/:id", UsersTryDemote)
	admins.GET("/users/del/:id", UsersTryDelete)

	// /config
	admins.GET("/config", ConfigEdit)
	admins.POST("/config/edit", ConfigTryEdit)
	admins.POST("/config/reset", ConfigTryReset)

	// /error
	pub.GET("/error", ErrorHome)

	pub.Run()
	// pub.Run(":3000") for a hard coded port
}
