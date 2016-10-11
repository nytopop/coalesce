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
}

type PostEditForm struct {
	PostId string `form:"postid" binding:"required"`
	Title  string `form:"title" binding:"required"`
	Body   string `form:"body" binding:"required"`
}

type Post struct {
	Id        bson.ObjectId `bson:"_id,omitempty"`
	Title     string        `bson:"title"`
	Author    string        `bson:"author"`
	Wordcount int           `bson:"wordcount"`
	Draft     bool          `bson:"draft"` //if true, no publish
	Timestamp time.Time     `bson:"timestamp"`
	Updated   time.Time     `bson:"updated"`
	Body      string        `bson:"body"`
	BodyHTML  template.HTML `bson:"bodyhtml"`
	Tags      []string      `bson:"tags"`
}

// TODO posts should be authed to a specific user

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
		c.Error(err)
		c.Redirect(302, "/error")
	}

	c.HTML(http.StatusOK, "posts/all.html", gin.H{
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
		c.Error(err)
		c.Redirect(302, "/error")
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

		// create tags using cortical.io
		tags, err := GetKeywordsForText(cfg.Server.ApiKey, postform.Body)
		if err != nil {
			c.Error(err)
			c.Redirect(302, "/error")
		}

		// construct post
		post := Post{
			Title:     postform.Title,
			Author:    user.Name,
			Body:      postform.Body,
			BodyHTML:  template.HTML(body),
			Timestamp: time.Now(),
			Updated:   time.Now(),
			Tags:      tags,
		}

		if err := s.Insert(&post); err != nil {
			c.Error(err)
			c.Redirect(302, "/error")
		}

		posturl := "/posts"
		c.Redirect(302, posturl)
	} else {
		c.Error(err)
		c.Redirect(302, "/error")
	}
}

// GET /posts/edit/:id
func PostsEdit(c *gin.Context) {
	session := globalSession.Copy()
	s := session.DB(cfg.Database.Name).C("posts")

	// get obj id from hex
	hexid := c.Param("id")
	id := bson.ObjectIdHex(hexid)

	// get post
	post := Post{}
	if err := s.FindId(id).One(&post); err != nil {
		c.Error(err)
		c.Redirect(302, "/error")
	}

	user := GetUser(c)

	if user.Name == post.Author || user.Name == "admin" {
		c.HTML(http.StatusOK, "posts/edit.html", gin.H{
			"Site": cfg.Site,
			"Post": post,
			"User": user,
		})
	} else {
		c.Redirect(302, "/auth/sign-in")
	}

}

// POST /posts/edit
func PostsTryEdit(c *gin.Context) {
	session := globalSession.Copy()
	s := session.DB(cfg.Database.Name).C("posts")

	user := GetUser(c)

	// validate form
	var postform PostEditForm
	if err := c.Bind(&postform); err == nil {

		// get obj id from hex
		hexid := postform.PostId
		id := bson.ObjectIdHex(hexid)

		// get timestamp from orig post
		oldpost := Post{}
		if err := s.FindId(id).One(&oldpost); err != nil {
			c.Error(err)
			c.Redirect(302, "/error")
		}

		if user.Name == oldpost.Author || user.Name == "admin" {
			// convert markdown
			body := string(blackfriday.MarkdownCommon([]byte(postform.Body)))

			// create tags using cortical.io
			tags, err := GetKeywordsForText(cfg.Server.ApiKey, postform.Body)
			if err != nil {
				c.Error(err)
				c.Redirect(302, "/error")
			}

			// construct updated post
			post := Post{
				Title:     postform.Title,
				Author:    user.Name,
				Body:      postform.Body,
				BodyHTML:  template.HTML(body),
				Timestamp: oldpost.Timestamp,
				Updated:   time.Now(),
				Tags:      tags,
			}

			// update post
			if err := s.UpdateId(id, post); err != nil {
				c.Error(err)
				c.Redirect(302, "/error")
			}

			posturl := "/posts/view/" + hexid
			c.Redirect(302, posturl)
		} else {
			c.Redirect(302, "/auth/sign-in")
		}
	} else {
		c.Error(err)
		c.Redirect(302, "/error")
	}
}

// GET /posts/del/:id
func PostsTryDelete(c *gin.Context) {
	session := globalSession.Copy()
	s := session.DB(cfg.Database.Name).C("posts")

	user := GetUser(c)

	// get obj id from hex
	hexid := c.Param("id")
	id := bson.ObjectIdHex(hexid)

	// get post
	post := Post{}
	if err := s.FindId(id).One(&post); err != nil {
		c.Error(err)
		c.Redirect(302, "/error")
	}

	if user.Name == post.Author || user.Name == "admin" {
		// delete post
		if err := s.RemoveId(id); err != nil {
			c.Error(err)
			c.Redirect(302, "/error")
		}

		c.Redirect(302, "/users/me")
	} else {
		c.Redirect(302, "/auth/sign-in")
	}
}
