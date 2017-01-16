// coalesce/controllers/auth.go

package controllers

import (
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/nytopop/coalesce/models"
	"github.com/nytopop/coalesce/util"
)

type SignInForm struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type RegisterForm struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

// Get active user
func GetUser(c *gin.Context) models.SQLUser {
	var userid int
	var name string
	var alevel int

	if x := c.MustGet("userid"); x != nil {
		userid = x.(int)
	}
	if x := c.MustGet("name"); x != nil {
		name = x.(string)
	}
	if x := c.MustGet("accesslevel"); x != nil {
		alevel = x.(int)
	}

	user := models.SQLUser{
		Userid:      userid,
		Name:        name,
		AccessLevel: alevel,
	}

	return user
}

// set name, authlevel
func AuthCheckpoint() gin.HandlerFunc {
	return func(c *gin.Context) {
		// cookies
		cookies := sessions.Default(c)

		// check if sign in cookie exists
		if cookies.Get("name") != nil {
			c.Set("userid", cookies.Get("userid"))
			c.Set("name", cookies.Get("name"))
			c.Set("accesslevel", cookies.Get("accesslevel"))
		} else {
			c.Set("userid", -1)
			c.Set("name", "guest")
			c.Set("accesslevel", 0)
		}
	}
}

// ensure user is auth'd at access level
func AccessLevelAuth(level int) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.MustGet("accesslevel").(int) < level {
			c.Redirect(302, "/auth/sign-in")
		}
	}
}

// GET /auth/sign-in
func AuthSignIn(c *gin.Context) {
	site := models.Site
	site.Title = "Sign in"
	c.HTML(http.StatusOK, "auth/sign-in.html", gin.H{
		"Site": site,
		"User": GetUser(c),
	})
}

// POST /auth/sign-in
func AuthTrySignIn(c *gin.Context) {
	// cookies
	cookies := sessions.Default(c)

	// validate auth form
	var authform SignInForm
	err := c.Bind(&authform)
	if err != nil {
		RenderErr(c, err)
		return
	}

	exists, err := models.QueryUserExists(authform.Username)
	if err != nil {
		RenderErr(c, err)
		return
	}

	if !exists {
		c.Redirect(302, "/auth/sign-in")
	}

	user, err := models.QueryUsername(authform.Username)
	if err != nil {
		RenderErr(c, err)
		return
	}

	// Password check : bcrypt
	err = util.CheckToken(user.Salt, authform.Password, user.Token)
	if err != nil {
		c.Redirect(302, "/auth/sign-in")
	} else {
		cookies.Set("userid", user.Userid)
		cookies.Set("name", user.Name)
		cookies.Set("accesslevel", user.AccessLevel)
		cookies.Save()

		c.Redirect(302, "/posts")
	}
}

// GET /auth/register
func AuthRegister(c *gin.Context) {
	site := models.Site
	site.Title = "Register"
	c.HTML(http.StatusOK, "auth/register.html", gin.H{
		"Site": site,
		"User": GetUser(c),
	})
}

// POST /auth/register
func AuthTryRegister(c *gin.Context) {
	// validate form
	var authform RegisterForm
	err := c.Bind(&authform)
	if err != nil {
		RenderErr(c, err)
		return
	}

	// check if user exists
	userExists, err := models.QueryUserExists(authform.Username)
	if err != nil {
		RenderErr(c, err)
		return
	}

	switch userExists {
	case true:
		c.Redirect(302, "/auth/register")
	case false:
		salt, err := util.GenerateSalt()
		if err != nil {
			RenderErr(c, err)
			return
		}

		token, err := util.ComputeToken(salt, authform.Password)
		if err != nil {
			RenderErr(c, err)
			return
		}

		// create the user from form
		user := models.SQLUser{
			Name:        authform.Username,
			Salt:        salt,
			Token:       token,
			AccessLevel: 1,
		}

		err = models.WriteUser(user)
		if err != nil {
			RenderErr(c, err)
			return
		}

		c.Redirect(302, "/posts")
	}
}

// GET /auth/sign-out
func AuthSignOut(c *gin.Context) {
	// cookies
	cookies := sessions.Default(c)

	cookies.Clear()
	cookies.Save()

	c.Redirect(302, "/posts")
}
