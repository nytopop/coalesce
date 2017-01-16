// coalesce/coalesce.go

package main

import (
	"flag"
	"log"
	"os"
	"runtime"

	gcfg "gopkg.in/gcfg.v1"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nytopop/coalesce/controllers"
	"github.com/nytopop/coalesce/models"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// process flags
	var configFile = flag.String("cfg", "/etc/coalesce.conf", "path to configuration file")
	flag.Parse()

	// Configuration file
	err := gcfg.ReadFileInto(&cfg, *configFile)
	if err != nil {
		log.Fatalln(err)
	}

	// Logging
	errlog, err := os.Create(cfg.System.ErrorLog)
	if err != nil {
		log.Fatalln(err)
	}
	acclog, err := os.Create(cfg.System.AccessLog)
	if err != nil {
		log.Fatalln(err)
	}
	logs := controllers.Logs{
		Error:  log.New(errlog, "", log.Ldate|log.Ltime),
		Access: log.New(acclog, "", log.Ldate|log.Ltime),
	}

	// Database initialization
	err = models.InitDB(cfg.System.Database, cfg.System.DatabaseInit)
	if err != nil {
		logs.Error.Fatalln(err)
	}
	defer models.CloseDB()

	models.Site.SiteTitle = cfg.Site.Title
	models.Site.Copyright = cfg.Site.Copyright
	models.Site.Email = cfg.Site.Email

	gin.SetMode(gin.ReleaseMode)
	pub := gin.New()

	// middleware
	pub.Use(controllers.Logger(logs))
	pub.Use(gin.Recovery())

	// session management TODO secret
	secret := []byte(os.Getenv("SESSION_SECRET"))
	store := sessions.NewCookieStore(secret)
	pub.Use(sessions.Sessions("coalesce", store))

	// authentication
	pub.Use(controllers.AuthCheckpoint())
	commentors := pub.Group("/", controllers.AccessLevelAuth(1))
	editors := pub.Group("/", controllers.AccessLevelAuth(2))
	admins := pub.Group("/", controllers.AccessLevelAuth(3))

	// templates
	pub.LoadHTMLGlob(cfg.System.ResourceDir + "/templates/**/*.html")

	// routes
	pub.Static("/static", cfg.System.ResourceDir+"/static")
	pub.StaticFile("favicon.png", cfg.System.ResourceDir+"/favicon.png")
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
	commentors.POST("/users/passchange", controllers.UsersTryPassChange)
	commentors.GET("/users/myposts", controllers.UsersMyPosts)
	admins.GET("/users/all", controllers.UsersAll)
	admins.GET("/users/promote/:id", controllers.UsersTryPromote)
	admins.GET("/users/demote/:id", controllers.UsersTryDemote)
	admins.GET("/users/del/:id", controllers.UsersTryDelete)

	// /config
	admins.GET("/config", controllers.ConfigEdit)
	admins.POST("/config/edit", controllers.ConfigTryEdit)
	admins.POST("/config/reset", controllers.ConfigTryReset)

	pub.Run(cfg.System.Listen)
}
