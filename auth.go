// coalesce/auth.go

package main

import (
	"crypto/sha512"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"

	"gopkg.in/mgo.v2/bson"
)

type SignInForm struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type RegisterForm struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type User struct {
	Id          bson.ObjectId `bson:"_id,omitempty"`
	Name        string        `bson:"name"`
	Token       string        `bson:"token"`
	AccessLevel int           `bson:"accesslevel"`
}

type SQLUser struct {
	Userid      int
	Name        string
	Token       string
	AccessLevel int
}

func GenerateSalt() string {
	return ""
}

// Get active user
func GetUser(c *gin.Context) SQLUser {
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

	user := SQLUser{
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
	c.HTML(http.StatusOK, "auth/sign-in.html", gin.H{
		//"Site": GetConf(),
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
		c.Error(err)
		c.Redirect(302, "/error")
	}

	/*
		Check if user exists, if not DENY
		Get user's salt, hash password with salt
		If not matched, DENY
	*/

	// create token; hash password
	hash := sha512.Sum512([]byte(authform.Password))
	token := hex.EncodeToString(hash[:])
	user, err := queryUser(authform.Username, token)
	if err != nil {
		c.Error(err)
		c.Redirect(302, "/error")
	}

	// set cookies
	switch user {
	case SQLUser{}: // no match
		c.Redirect(302, "/auth/sign-in")
	default: // user matched
		cookies.Set("userid", user.Userid)
		cookies.Set("name", user.Name)
		cookies.Set("accesslevel", user.AccessLevel)
		cookies.Save()

		c.Redirect(302, "/posts")
	}

}

// GET /auth/register
func AuthRegister(c *gin.Context) {
	c.HTML(http.StatusOK, "auth/register.html", gin.H{
		//"Site": GetConf(),
		"User": GetUser(c),
	})
}

// POST /auth/register
func AuthTryRegister(c *gin.Context) {
	// validate form
	var authform RegisterForm
	err := c.Bind(&authform)
	if err != nil {
		c.Error(err)
		c.Redirect(302, "/error")
	}

	// check if user exists
	userExists, err := queryUserExists(authform.Username)
	if err != nil {
		c.Error(err)
		c.Redirect(302, "/error")
	}

	switch userExists {
	case true:
		c.Redirect(302, "/auth/register")
	case false:
		hash := sha512.Sum512([]byte(authform.Password))
		token := hex.EncodeToString(hash[:])

		// create the user from form
		user := SQLUser{
			Name:        authform.Username,
			Token:       token,
			AccessLevel: 2,
		}

		err = writeUser(user)
		if err != nil {
			c.Error(err)
			c.Redirect(302, "/error")
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
