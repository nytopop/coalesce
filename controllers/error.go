// coalesce/controllers/error.go

package controllers

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nytopop/coalesce/models"
)

type Logs struct {
	Error  *log.Logger
	Access *log.Logger
}

func Logger(logs Logs) gin.HandlerFunc {
	return func(c *gin.Context) {
		// start timer
		start := time.Now()
		path := c.Request.URL.Path

		// process request
		c.Next()

		// collect data for log
		end := time.Now()
		latency := end.Sub(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		status := c.Writer.Status()
		errors := c.Errors.String()

		// write log
		switch {
		case status >= 200 && status < 400: // All good signals
			logs.Access.Println(
				status, method, latency,
				clientIP, path, errors)
		case status >= 400 && status < 600: // Errors
			logs.Error.Println(
				status, method, latency,
				clientIP, path, errors)
		}
	}
}

// Render Error
func RenderErr(c *gin.Context, err error) {
	c.Error(err)

	site := models.Site
	site.Title = "Error"
	c.HTML(500, "misc/error.html", gin.H{
		"Site":  site,
		"User":  GetUser(c),
		"Error": err,
	})
}
