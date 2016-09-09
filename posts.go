// coalesce/posts.go

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

type PostForm struct {
	Title  string `form:"title" binding:"required"`
	Author string `form:"author" binding:"required"`
	Body   string `form:"body" binding:"required"`
	Tags   string `form:"tags" binding:"required"`
}

type Post struct {
	Id        bson.ObjectId `bson:"_id,omitempty"`
	Title     string        `bson:"title"`
	Author    string        `bson:"author"`
	Wordcount int           `bson:"wordcount"`
	Timestamp time.Time     `bson:"timestamp"`
	Body      string        `bson:"body"`
	BodyHTML  template.HTML `bson:"bodyhtml"`
	Comments  []Comment     `bson:"comments"`
	Tags      []string      `bson:"tags"`
}

type Comment struct {
	Author    string    `bson:"author"`
	Timestamp time.Time `bson:"timestamp"`
	Body      string    `bson:"body"`
}

// GET /posts
func PostsHome(c *gin.Context) {
	session := globalSession.Copy()
	s := session.DB(cfg.Database.Name).C("posts")

	posts := []*Post{}
	if err := s.Find(nil).Sort("-timestamp").Iter().All(&posts); err != nil {
		log.Println(err)
		return
	}

	c.HTML(http.StatusOK, "posts/home.html", gin.H{
		"Site": cfg.Site,
		"List": posts,
	})
}

// GET /posts/view/:id
func PostsView(c *gin.Context) {
	session := globalSession.Copy()
	s := session.DB(cfg.Database.Name).C("posts")

	// get obj id from hex
	hexid := c.Param("id")
	id := bson.ObjectIdHex(hexid)

	post := Post{}
	if err := s.FindId(id).One(&post); err != nil {
		log.Println(err)
	}

	c.HTML(http.StatusOK, "posts/view.html", gin.H{
		"Site": cfg.Site,
		"Post": post,
	})
}

// GET /posts/new
func PostsNew(c *gin.Context) {
	c.HTML(http.StatusOK, "posts/new.html", gin.H{
		"Site": cfg.Site,
	})
}

// POST /posts/new
func PostsTryNew(c *gin.Context) {
	session := globalSession.Copy()
	s := session.DB(cfg.Database.Name).C("posts")

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

		c.Redirect(302, "/posts")
	}
}
