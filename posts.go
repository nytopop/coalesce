// coalesce/posts.go

package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

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
	Comments  []Comment     `bson:"comments"`
	Tags      []string      `bson:"tags"`
}

type Comment struct {
	Author    string    `bson:"author"`
	Timestamp time.Time `bson:"timestamp"`
	Body      string    `bson:"body"`
}

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

func PostsNew(c *gin.Context) {
	//	session := globalSession.Copy()
	//	s := session.DB(cfg.Database.Name).C("posts")

	// create a page with a form, defaults filled out
	// body should be markdown acceptable

	c.HTML(http.StatusOK, "posts/new.html", gin.H{
		"Site": cfg.Site,
	})
}

func PostsCreate(c *gin.Context) {
	session := globalSession.Copy()
	s := session.DB(cfg.Database.Name).C("posts")

	// validate form
	var postform PostForm
	if err := c.Bind(&postform); err == nil {
		post := Post{
			Title:     postform.Title,
			Author:    postform.Author,
			Body:      postform.Body,
			Timestamp: time.Now(),
		}

		if err := s.Insert(&post); err != nil {
			log.Println(err)
		}

		c.Redirect(302, "/posts")
	}
}
