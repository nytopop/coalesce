// coalesce/pages.go

package main

import "github.com/gin-gonic/gin"

func PagesHome(c *gin.Context) {
	c.Redirect(302, "/posts")
}
