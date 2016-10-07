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
	Title string `form:"title" binding:"required"`
	Body  string `form:"body" binding:"required"`
	Tags  string `form:"tags" binding:"required"`
}

type Post struct {
	Id        bson.ObjectId `bson:"_id,omitempty"`
	Title     string        `bson:"title"`
	Author    string        `bson:"author"`
	Wordcount int           `bson:"wordcount"`
	Timestamp time.Time     `bson:"timestamp"`
	Body      string        `bson:"body"`
	BodyHTML  template.HTML `bson:"bodyhtml"`
	Tags      []string      `bson:"tags"`
}

// get entire comment hierarchy of post
func (p Post) CommentTree() []Comment {
	session := globalSession.Copy()
	s := session.DB(cfg.Database.Name).C("comments")

	// get all root level comments
	id := p.Id
	comments := []*Comment{}
	query := bson.M{
		"postid": id,
		"depth":  0,
	}
	if err := s.Find(query).Sort("-timestamp").Iter().All(&comments); err != nil {
		log.Println(err)
	}

	// construct tree from root comments
	tree := []Comment{}
	for _, v := range comments {
		tree = append(tree, v.Tree()...)
	}

	return tree
}

// GET /posts
func PostsHome(c *gin.Context) {
	session := globalSession.Copy()
	s := session.DB(cfg.Database.Name).C("posts")

	// get posts
	posts := []*Post{}
	if err := s.Find(nil).Sort("-timestamp").Iter().All(&posts); err != nil {
		log.Println(err)
		return
	}

	c.HTML(http.StatusOK, "posts/home.html", gin.H{
		"Site": cfg.Site,
		"List": posts,
		"User": GetUser(c),
	})
}

// GET /posts/view/:id
func PostsView(c *gin.Context) {
	session := globalSession.Copy()
	s := session.DB(cfg.Database.Name).C("posts")

	// get obj id from hex
	hexid := c.Param("id")
	id := bson.ObjectIdHex(hexid)

	// get post
	post := Post{}
	if err := s.FindId(id).One(&post); err != nil {
		log.Println(err)
	}

	// get comments
	tree := post.CommentTree()

	c.HTML(http.StatusOK, "posts/view.html", gin.H{
		"Site":     cfg.Site,
		"Post":     post,
		"Comments": tree,
		"User":     GetUser(c),
	})
}

// GET /posts/new
func PostsNew(c *gin.Context) {
	c.HTML(http.StatusOK, "posts/new.html", gin.H{
		"Site": cfg.Site,
		"User": GetUser(c),
	})
}

// POST /posts/new
func PostsTryNew(c *gin.Context) {
	session := globalSession.Copy()
	s := session.DB(cfg.Database.Name).C("posts")

	user := GetUser(c)

	// validate form
	var postform PostForm
	if err := c.Bind(&postform); err == nil {
		// convert markdown
		body := string(blackfriday.MarkdownCommon([]byte(postform.Body)))

		post := Post{
			Title:     postform.Title,
			Author:    user.Name,
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
