// coalesce/auth.go

package main

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"net/http"
	"strconv"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

type SignInForm struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type RegisterForm struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type SQLUser struct {
	Userid      int
	Name        string
	Salt        string
	Token       string
	AccessLevel int
}

func GeneratePepper() (string, error) {
	r := make([]byte, 1)
	_, err := rand.Read(r)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(int(r[0])), nil
}

func GenerateSalt() (string, error) {
	r := make([]byte, 32)
	_, err := rand.Read(r)
	if err != nil {
		return "", err
	}

	hash := sha512.Sum512(r)
	return hex.EncodeToString(hash[:]), nil
}

func ComputeToken(salt, pepper, pw string) string {
	chars := salt + pepper + pw
	hash := sha512.Sum512([]byte(chars))
	return hex.EncodeToString(hash[:])
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
		RenderErr(c, err)
		return
	}

	exists, err := queryUserExists(authform.Username)
	if err != nil {
		RenderErr(c, err)
		return
	}

	if !exists {
		c.Redirect(302, "/auth/sign-in")
	}

	user, err := queryUsername(authform.Username)
	if err != nil {
		RenderErr(c, err)
		return
	}

	// iterate through peppers until we get a match
	match := false
	for i := 0; i < 256; i++ {
		token := ComputeToken(user.Salt, strconv.Itoa(i), authform.Password)
		if token == user.Token {
			match = true
			break
		}
	}

	if !match {
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
		RenderErr(c, err)
		return
	}

	// check if user exists
	userExists, err := queryUserExists(authform.Username)
	if err != nil {
		RenderErr(c, err)
		return
	}

	switch userExists {
	case true:
		c.Redirect(302, "/auth/register")
	case false:
		salt, err := GenerateSalt()
		if err != nil {
			RenderErr(c, err)
			return
		}

		pepper, err := GeneratePepper()
		if err != nil {
			RenderErr(c, err)
			return
		}

		token := ComputeToken(salt, pepper, authform.Password)

		// create the user from form
		user := SQLUser{
			Name:        authform.Username,
			Salt:        salt,
			Token:       token,
			AccessLevel: 1,
		}

		err = writeUser(user)
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
