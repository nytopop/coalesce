// coalesce/controllers/error.go

package controllers

import "github.com/gin-gonic/gin"

// GET /error
func ErrorHome(c *gin.Context) {
	c.HTML(200, "misc/error.html", gin.H{
		"User": GetUser(c),
	})
}

// Render Error
func RenderErr(c *gin.Context, err error) {
	c.HTML(500, "misc/error.html", gin.H{
		"Error": err,
		"User":  GetUser(c),
	})
}
