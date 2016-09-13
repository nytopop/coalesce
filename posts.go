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
	Tags      []string      `bson:"tags"`
}

type CommentForm struct {
	PostId string `form:"postid" binding:"required"`
	Author string `form:"author" binding:"required"`
	Body   string `form:"body" binding:"required"`
}

type Comment struct {
	Id        bson.ObjectId   `bson:"_id,omitempty"`
	PostId    bson.ObjectId   `bson:"postid"`
	Author    string          `bson:"author"`
	Timestamp time.Time       `bson:"timestamp"`
	Body      string          `bson:"body"`
	Depth     int             `bson:"depth"`
	Replies   []bson.ObjectId `bson:"replies"`
}

// make a comment reply
func (c Comment) Reply(doc Comment) {
	//doc.Depth = c.Depth + 1

	// write the new doc, get its ID

	//c.Replies = append(c.Replies, docid)

	// update the new comment reply in orig comment
}

// recursive comment chain flattener
func (c Comment) Tree() []Comment {
	// load db session
	session := globalSession.Copy()
	s := session.DB(cfg.Database.Name).C("comments")

	// add self to result set
	tree := []Comment{}
	tree = append(tree, c)

	if len(c.Replies) > 0 {
		// if there is a deeper level, keep going
		for _, v := range c.Replies {
			// go deeper
			next := Comment{}

			// get the referenced comment
			if err := s.FindId(v).One(&next); err != nil {
			}

			// iterate through next comment
			for _, vv := range next.Tree() {
				tree = append(tree, vv)
			}
		}
	}
	return tree
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
	coms := session.DB(cfg.Database.Name).C("comments")

	// get obj id from hex
	hexid := c.Param("id")
	id := bson.ObjectIdHex(hexid)

	post := Post{}
	if err := s.FindId(id).One(&post); err != nil {
		log.Println(err)
	}

	// get all root level comments
	comments := []*Comment{}
	query := bson.M{
		"postid": id,
		"depth":  0,
	}
	if err := coms.Find(query).Sort("-timestamp").Iter().All(&comments); err != nil {
		log.Println(err)
	}

	tree := []Comment{}
	for _, v := range comments {
		for _, vv := range v.Tree() {
			tree = append(tree, vv)
		}
	}

	c.HTML(http.StatusOK, "posts/view.html", gin.H{
		"Site":     cfg.Site,
		"Post":     post,
		"Comments": tree,
		// need comments as well
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

// POST /posts/comment
func PostsTryComment(c *gin.Context) {
	session := globalSession.Copy()
	s := session.DB(cfg.Database.Name).C("comments")

	// validate
	var cform CommentForm
	if err := c.Bind(&cform); err == nil {
		// get obj id from hex
		hexid := cform.PostId
		id := bson.ObjectIdHex(hexid)

		comment := Comment{
			PostId: id,
			Author: cform.Author,
			Body:   cform.Body,
			Depth:  0,
		}

		if err := s.Insert(&comment); err != nil {

		}

		// redirect to orig post
		redir := "/posts/view/" + hexid
		c.Redirect(302, redir)
	}
}

// POST /posts/comment/reply
func PostsTryCommentReply(c *gin.Context) {
	// same as above EXCEPT
	// create a comment
	// find a specific comment by id, append new comment to replies
}
