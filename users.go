// coalesce/users.go

package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// GET /users/all
func UsersAll(c *gin.Context) {
	users, err := queryUsersAll()
	if err != nil {
		RenderErr(c, err)
		return
	}

	c.HTML(200, "users/all.html", gin.H{
		"Users": users,
		"User":  GetUser(c),
	})
}

// GET /users/promote/:id
func UsersTryPromote(c *gin.Context) {
	p := c.Param("id")
	pNum, err := strconv.Atoi(p)
	if err != nil {
		RenderErr(c, err)
		return
	}

	user, err := queryUserID(pNum)
	if err != nil {
		RenderErr(c, err)
		return
	}

	user.AccessLevel = 2
	err = updateUser(user)
	if err != nil {
		RenderErr(c, err)
		return
	}

	c.Redirect(302, "/users/all")
}

// GET /users/demote/:id
func UsersTryDemote(c *gin.Context) {
	p := c.Param("id")
	pNum, err := strconv.Atoi(p)
	if err != nil {
		RenderErr(c, err)
		return
	}

	user, err := queryUserID(pNum)
	if err != nil {
		RenderErr(c, err)
		return
	}

	user.AccessLevel = 1
	err = updateUser(user)
	if err != nil {
		RenderErr(c, err)
		return
	}

	c.Redirect(302, "/users/all")
}

// GET /users/del/:id
func UsersTryDelete(c *gin.Context) {
	p := c.Param("id")
	pNum, err := strconv.Atoi(p)
	if err != nil {
		RenderErr(c, err)
		return
	}

	err = deleteUser(pNum)
	if err != nil {
		RenderErr(c, err)
		return
	}

	c.Redirect(302, "/users/all")
}

// GET /users/me
func UsersMe(c *gin.Context) {
	c.HTML(200, "users/me.html", gin.H{
		"User": GetUser(c),
	})
}

// GET /users/myposts
func UsersMyPosts(c *gin.Context) {
	user := GetUser(c)
	posts, err := queryPostsUserID(user.Userid)
	if err != nil {
		RenderErr(c, err)
		return
	}

	c.HTML(200, "users/myposts.html", gin.H{
		"Posts": posts,
		"User":  user,
	})
}

// GET /users/mycomments
func UsersMyComments(c *gin.Context) {
	c.HTML(200, "users/mycomments.html", gin.H{
		"User": GetUser(c),
	})
}
