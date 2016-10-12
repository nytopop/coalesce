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

// Get active user
func GetUser(c *gin.Context) User {
	var name string
	var alevel int

	if x := c.MustGet("accesslevel"); x != nil {
		alevel = x.(int)
	}

	if x := c.MustGet("name"); x != nil {
		name = x.(string)
	}

	user := User{
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
			c.Set("name", cookies.Get("name"))
			c.Set("accesslevel", cookies.Get("accesslevel"))
		} else {
			c.Set("name", "guest")
			c.Set("accesslevel", 0)
		}
	}
}

// ensure user is auth'd at access level
func AccessLevelAuth(level int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// BUG: mustget borked
		if c.MustGet("accesslevel").(int) < level {
			c.Redirect(302, "/auth/sign-in")
		}
	}
}

// GET /auth/sign-in
func AuthSignIn(c *gin.Context) {
	c.HTML(http.StatusOK, "auth/sign-in.html", gin.H{
		"Site": cfg.Site,
		"User": GetUser(c),
	})
}

// POST /auth/sign-in
func AuthTrySignIn(c *gin.Context) {
	// db
	session := globalSession.Copy()
	s := session.DB(dbname).C("users")

	// cookies
	cookies := sessions.Default(c)

	// validate auth form
	var authform SignInForm
	if err := c.Bind(&authform); err == nil {
		// create token; hash password
		hash := sha512.Sum512([]byte(authform.Password))
		token := hex.EncodeToString(hash[:])

		// construct query
		query := bson.M{
			"name":  authform.Username,
			"token": token,
		}

		// query user
		user := User{}
		if err := s.Find(query).One(&user); err != nil {
			c.Error(err)
			c.Redirect(302, "/error")
		}

		// set cookies
		if user != (User{}) {
			// success
			cookies.Set("name", user.Name)
			cookies.Set("accesslevel", user.AccessLevel)
			cookies.Save()

			c.Redirect(302, "/posts")
		} else {
			// fail
			c.Redirect(302, "/auth/sign-in")
		}
	} else {
		c.Error(err)
		c.Redirect(302, "/error")
	}

}

// GET /auth/register
func AuthRegister(c *gin.Context) {
	c.HTML(http.StatusOK, "auth/register.html", gin.H{
		"Site": cfg.Site,
		"User": GetUser(c),
	})
}

// POST /auth/register
func AuthTryRegister(c *gin.Context) {
	session := globalSession.Copy()
	s := session.DB(dbname).C("users")

	// validate form
	var authform RegisterForm
	if err := c.Bind(&authform); err == nil {
		// check if user exists already
		query := bson.M{
			"name": authform.Username,
		}

		if n, _ := s.Find(query).Count(); n > 0 {
			// user exists
			c.Redirect(302, "/auth/register")
		} else {
			// no existing user, good 2 go
			hash := sha512.Sum512([]byte(authform.Password))
			token := hex.EncodeToString(hash[:])

			// create user
			user := User{
				Name:        authform.Username,
				Token:       token,
				AccessLevel: 1,
			}

			// write db
			if err := s.Insert(&user); err != nil {
				c.Error(err)
				c.Redirect(302, "/error")
			}

			c.Redirect(302, "/posts")
		}
	} else {
		c.Error(err)
		c.Redirect(302, "/error")
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

// Create the admin user
func CreateAdmin() {
	session := globalSession.Copy()
	s := session.DB(dbname).C("users")

	// no existing user, good 2 go
	hash := sha512.Sum512([]byte(cfg.Server.AdminPassword))
	token := hex.EncodeToString(hash[:])

	admin := User{
		Name:        "admin",
		Token:       token,
		AccessLevel: 3,
	}

	query := bson.M{
		"name": "admin",
	}

	// write db
	if _, err := s.Upsert(&query, &admin); err != nil {
		// handle error
	}
}
