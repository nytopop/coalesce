// coalesce/pages.go

package main

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday"

	"gopkg.in/mgo.v2/bson"
)

func PagesHome(c *gin.Context) {
	c.Redirect(302, "/posts")
}

func PagesView(c *gin.Context) {
	session := globalSession.Copy()
	s := session.DB(cfg.Database.Name).C("pages")

	// get obj id from hex
	hexid := c.Param("id")
	id := bson.ObjectIdHex(hexid)

	post := Post{}
	if err := s.FindId(id).One(&post); err != nil {
		log.Println(err)
	}

	c.HTML(http.StatusOK, "pages/view.html", gin.H{
		"Site": cfg.Site,
		"Post": post,
	})
}

func PagesNew(c *gin.Context) {
	c.HTML(http.StatusOK, "pages/new.html", gin.H{
		"Site": cfg.Site,
	})
}

func PagesTryNew(c *gin.Context) {
	session := globalSession.Copy()
	s := session.DB(cfg.Database.Name).C("pages")

	// validate form
	var postform PostForm
	if err := c.Bind(&postform); err == nil {
		// convert markdown
		body := string(blackfriday.MarkdownCommon([]byte(postform.Body)))

		post := Post{
			Title:     postform.Title,
			Author:    postform.Author,
			Body:      postform.Body,
			BodyHTML:  template.HTML(body),
			Timestamp: time.Now(),
		}

		if err := s.Insert(&post); err != nil {
			log.Println(err)
		}

		c.Redirect(302, "/pages")
	}
}
