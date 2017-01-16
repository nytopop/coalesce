// coalesce/controllers/posts.go

package controllers

import (
	"errors"
	"html/template"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nytopop/coalesce/models"
	"github.com/nytopop/coalesce/util"
	"github.com/russross/blackfriday"
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

// GET /
func Home(c *gin.Context) {
	c.Redirect(302, "/posts")
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
	posts, err := models.QueryPostsPage(pNum, 15)
	if err != nil {
		RenderErr(c, err)
		return
	}

	posts, err = models.ProcessPosts(posts)
	if err != nil {
		RenderErr(c, err)
		return
	}

	site := models.Site
	site.Title = "posts"

	// if this is not page 0, prev should work
	if pNum > 0 {
		site.Prev = pNum - 1
	} else {
		site.Prev = -1
	}

	np, err := models.QueryPostsPage(pNum+1, 15)
	if err != nil {
		RenderErr(c, err)
		return
	}

	// if there is a next page, use it
	if len(np) > 0 {
		site.Next = pNum + 1
	} else {
		site.Next = -1
	}

	// render
	c.HTML(200, "posts/page.html", gin.H{
		"Site":  site,
		"User":  GetUser(c),
		"Posts": posts,
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

	post, err := models.QueryPost(pNum)
	if err != nil {
		RenderErr(c, err)
		return
	}

	if post == (models.SQLPost{}) {
		RenderErr(c, errors.New("Post not found"))
		return
	}

	user, err := models.QueryUserID(post.Userid)
	if err != nil {
		RenderErr(c, err)
		return
	}

	post.Username = user.Name
	post.PostedNice = util.NiceTime(post.Posted)
	post.UpdatedNice = util.NiceTime(post.Updated)
	post.RenderHTML = template.HTML(post.BodyHTML)

	// comments!
	comments, err := CommentsForPost(pNum)
	if err != nil {
		RenderErr(c, err)
		return
	}

	site := models.Site
	site.Title = post.Title

	c.HTML(200, "posts/view.html", gin.H{
		"Site":     site,
		"User":     GetUser(c),
		"Post":     post,
		"Comments": comments,
	})
}

// GET /posts/new
func PostsNew(c *gin.Context) {
	site := models.Site
	site.Title = "write"
	c.HTML(200, "posts/new.html", gin.H{
		"Site": site,
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
	post := models.SQLPost{
		Userid:     user.Userid,
		Title:      postform.Title,
		Body:       postform.Body,
		BodyHTML:   html,
		Categoryid: 0,
		Posted:     time.Now().Unix(),
		Updated:    time.Now().Unix(),
	}

	err = models.WritePost(post)
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

	post, err := models.QueryPost(pNum)
	if err != nil {
		RenderErr(c, err)
		return
	}

	user := GetUser(c)

	if user.Userid == post.Userid {
		site := models.Site
		site.Title = "edit"
		c.HTML(200, "posts/edit.html", gin.H{
			"Site": site,
			"User": user,
			"Post": post,
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
	oldPost, err := models.QueryPost(pNum)
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
		newPost := models.SQLPost{
			Postid:   pNum,
			Title:    postform.Title,
			Body:     postform.Body,
			BodyHTML: html,
			Updated:  time.Now().Unix(),
		}

		err = models.UpdatePost(newPost)
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

	post, err := models.QueryPost(pNum)
	if err != nil {
		RenderErr(c, err)
		return
	}

	user := GetUser(c)

	if user.Userid != post.Userid {
		c.Redirect(302, "/auth/sign-in")
	} else {
		err = models.DeletePost(pNum)
		if err != nil {
			RenderErr(c, err)
			return
		}

		c.Redirect(302, "/posts")
	}
}
