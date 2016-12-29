// coalesce/posts.go

package main

import (
	"errors"
	"html/template"
	"strconv"
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

type SQLPost struct {
	Postid     int
	Userid     int
	Title      string
	Body       string
	BodyHTML   string
	RenderHTML template.HTML
	Categoryid int
	Posted     int64
	Updated    int64
}

// get entire comment hierarchy of post
func (p Post) CommentTree() []Comment {
	/*session := globalSession.Copy()
	s := session.DB(dbname).C("comments")*/

	// get all root level comments
	//id := p.Id
	comments := []*Comment{}
	/*query := bson.M{
		"postid": id,
		"depth":  0,
	}
	if err := s.Find(query).Sort("-timestamp").Iter().All(&comments); err != nil {
		log.Println(err)
	}*/

	// construct tree from root comments
	tree := []Comment{}
	for _, v := range comments {
		tree = append(tree, v.Tree()...)
	}

	return tree
}

// GET /posts[?p=[0,1,2,...]]
func PostsPage(c *gin.Context) {
	// Get page number from arg
	p := c.Query("p")
	var pNum int
	if p != "" {
		var err error
		pNum, err = strconv.Atoi(p)
		if err != nil {
			RenderErr(c, err)
			return
		}
	} else {
		pNum = 0
	}

	// get posts
	posts, err := queryPostsPage(pNum)
	if err != nil {
		RenderErr(c, err)
		return
	}

	// render
	c.HTML(200, "posts/page.html", gin.H{
		"Posts": posts,
		"User":  GetUser(c),
	})
}

// GET /posts/view/:id
func PostsView(c *gin.Context) {
	p := c.Param("id")
	pNum, err := strconv.Atoi(p)
	if err != nil {
		RenderErr(c, err)
		return
	}

	post, err := queryPost(pNum)
	if err != nil {
		RenderErr(c, err)
		return
	}

	if post == (SQLPost{}) {
		RenderErr(c, errors.New("Post not found"))
		return
	}

	// comments!
	post.RenderHTML = template.HTML(post.BodyHTML)
	c.HTML(200, "posts/view.html", gin.H{
		"Post": post,
		"User": GetUser(c),
	})
}

// GET /posts/new
func PostsNew(c *gin.Context) {
	c.HTML(200, "posts/new.html", gin.H{
		//			"Site": GetConf(),
		"User": GetUser(c),
	})
}

// POST /posts/new
func PostsTryNew(c *gin.Context) {
	var postform PostForm
	err := c.Bind(&postform)
	if err != nil {
		RenderErr(c, err)
		return
	}

	user := GetUser(c)

	html := string(blackfriday.MarkdownCommon([]byte(postform.Body)))
	post := SQLPost{
		Userid:     user.Userid,
		Title:      postform.Title,
		Body:       postform.Body,
		BodyHTML:   html,
		Categoryid: 0,
		Posted:     time.Now().Unix(),
		Updated:    time.Now().Unix(),
	}

	err = writePost(post)
	if err != nil {
		RenderErr(c, err)
		return
	}

	c.Redirect(302, "/posts")
}

// GET /posts/edit/:id
func PostsEdit(c *gin.Context) {
	p := c.Param("id")
	pNum, err := strconv.Atoi(p)
	if err != nil {
		RenderErr(c, err)
		return
	}

	post, err := queryPost(pNum)
	if err != nil {
		RenderErr(c, err)
		return
	}

	user := GetUser(c)

	if user.Userid == post.Userid {
		c.HTML(200, "posts/edit.html", gin.H{
			"Post": post,
			"User": user,
		})
	} else {
		c.Redirect(302, "/auth/sign-in")
	}
}

// POST /posts/edit/:id
func PostsTryEdit(c *gin.Context) {
	p := c.Param("id")
	pNum, err := strconv.Atoi(p)
	if err != nil {
		RenderErr(c, err)
		return
	}

	user := GetUser(c)
	oldPost, err := queryPost(pNum)
	if err != nil {
		RenderErr(c, err)
		return
	}

	if user.Userid != oldPost.Userid {
		c.Redirect(302, "/auth/sign-in")
	} else {
		var postform PostForm
		err = c.Bind(&postform)
		if err != nil {
			RenderErr(c, err)
			return
		}

		html := string(blackfriday.MarkdownCommon([]byte(postform.Body)))
		newPost := SQLPost{
			Postid:   pNum,
			Title:    postform.Title,
			Body:     postform.Body,
			BodyHTML: html,
			Updated:  time.Now().Unix(),
		}

		err = updatePost(newPost)
		if err != nil {
			RenderErr(c, err)
			return
		}

		c.Redirect(302, "/posts")
	}
}

// GET /posts/del/:id
func PostsTryDelete(c *gin.Context) {
	p := c.Param("id")
	pNum, err := strconv.Atoi(p)
	if err != nil {
		RenderErr(c, err)
		return
	}

	post, err := queryPost(pNum)
	if err != nil {
		RenderErr(c, err)
		return
	}

	user := GetUser(c)

	if user.Userid != post.Userid {
		c.Redirect(302, "/auth/sign-in")
	} else {
		err = deletePost(pNum)
		if err != nil {
			RenderErr(c, err)
			return
		}

		c.Redirect(302, "/posts")
	}
}
