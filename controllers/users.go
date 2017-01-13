// coalesce/controllers/users.go

package controllers

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nytopop/coalesce/models"
	"github.com/nytopop/coalesce/util"
)

// GET /users/all
func UsersAll(c *gin.Context) {
	users, err := models.QueryUsersAll()
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

	user, err := models.QueryUserID(pNum)
	if err != nil {
		RenderErr(c, err)
		return
	}

	user.AccessLevel = 2
	err = models.UpdateUser(user)
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

	user, err := models.QueryUserID(pNum)
	if err != nil {
		RenderErr(c, err)
		return
	}

	user.AccessLevel = 1
	err = models.UpdateUser(user)
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

	err = models.DeleteUser(pNum)
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

// PassChangeForm
type PassChangeForm struct {
	OldPassword  string `form:"oldpassword" binding:"required"`
	NewPassword1 string `form:"newpassword1" binding:"required"`
	NewPassword2 string `form:"newpassword2" binding:"required"`
}

// POST /users/passchange
func UsersTryPassChange(c *gin.Context) {
	var pform PassChangeForm
	err := c.Bind(&pform)
	if err != nil {
		RenderErr(c, err)
		return
	}

	if pform.NewPassword1 != pform.NewPassword2 {
		RenderErr(c, errors.New("Passwords do not match!"))
		return
	}

	cuser := GetUser(c)
	user, err := models.QueryUserID(cuser.Userid)
	if err != nil {
		RenderErr(c, err)
		return
	}

	err = util.CheckToken(user.Salt, pform.OldPassword, user.Token)
	if err != nil {
		RenderErr(c, err)
		return
	} else {
		// we update the user with new password
		newtoken, err := util.ComputeToken(user.Salt, pform.NewPassword1)
		if err != nil {
			RenderErr(c, err)
			return
		}

		user.Token = newtoken
		err = models.UpdateUser(user)
		if err != nil {
			RenderErr(c, err)
			return
		}

		c.Redirect(302, "/users/me")
	}
}

// GET /users/myposts
func UsersMyPosts(c *gin.Context) {
	user := GetUser(c)
	posts, err := models.QueryPostsUserID(user.Userid)
	if err != nil {
		RenderErr(c, err)
		return
	}

	posts, err = models.ProcessPosts(posts)
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
