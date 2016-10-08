// coalesce/error.go

package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GET /error
func ErrorHome(c *gin.Context) {
	c.HTML(http.StatusOK, "misc/error.html", gin.H{
		"Site": cfg.Site,
		"User": GetUser(c),
	})
}
