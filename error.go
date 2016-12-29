// coalesce/error.go

package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GET /error
func ErrorHome(c *gin.Context) {
	c.HTML(http.StatusOK, "misc/error.html", gin.H{
		"Site": GetConf(),
		"User": GetUser(c),
		"Errs": c.Errors.Errors(),
	})
}

// Render Error
func RenderErr(c *gin.Context, err error) {
	c.HTML(500, "misc/error.html", gin.H{
		"Error": err,
		"User":  GetUser(c),
	})
}
