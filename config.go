// coalesce/config.go

package main

import "github.com/gin-gonic/gin"

// config file
type Config struct {
	Server struct {
		ApiKey        string
		AdminPassword string
	}
	Site struct {
		Title       string
		Description string
		Owner       string
		Github      string
		Email       string
	}
}

// GET /config
func ConfigEdit(c *gin.Context) {
	//	session := globalSession.Copy()
	//	s := session.DB(dbname).C("config")

	// geet those configs
}

// POST /config/edit
func ConfigTryEdit(c *gin.Context) {

}

// POST /config/reset
func ConfigTryReset(c *gin.Context) {

}
