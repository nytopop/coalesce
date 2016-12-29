// coalesce/error.go

package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GET /error
func ErrorHome(c *gin.Context) {
	log.Println(c.Errors.Errors())

	c.HTML(http.StatusOK, "misc/error.html", gin.H{
		"Site": GetConf(),
		"User": GetUser(c),
		"Errs": c.Errors.Errors(),
	})
}
