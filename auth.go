// coalesce/auth.go

package main

import (
	"crypto/sha512"
	"encoding/hex"
	"log"
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

// check if user is auth'd at the supplied access level
func AccessLevelAuth(level int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// db
		session := globalSession.Copy()
		s := session.DB(cfg.Database.Name).C("users")

		// cookies
		cookies := sessions.Default(c)

		// check if user exists, and is val by token
		query := bson.M{
			"name":  cookies.Get("name"),
			"token": cookies.Get("token"),
			//"accesslevel": level,
		}

		user := User{}
		if err := s.Find(query).One(&user); err != nil {
			log.Println(err)
		}

		if user != (User{}) {
			if user.AccessLevel >= level {
				c.Set("name", cookies.Get("name"))
			} else {
				c.Redirect(302, "/auth/sign-in")
			}
		} else {
			c.Redirect(302, "/auth/sign-in")
		}
		// c.MustGet("name") to check
	}
}

// GET /auth/sign-in
func AuthSignIn(c *gin.Context) {
	c.HTML(http.StatusOK, "auth/sign-in.html", gin.H{
		"Site": cfg.Site,
	})
}

// POST /auth/sign-in
func AuthTrySignIn(c *gin.Context) {
	// db
	session := globalSession.Copy()
	s := session.DB(cfg.Database.Name).C("users")

	// cookies
	cookies := sessions.Default(c)

	// validate auth form
	var authform SignInForm
	if err := c.Bind(&authform); err == nil {
		// create token; hash password
		hash := sha512.Sum512([]byte(authform.Password))
		token := hex.EncodeToString(hash[:])

		// check if user exists, and is val by token
		query := bson.M{
			"name":  authform.Username,
			"token": token,
		}

		if n, _ := s.Find(query).Count(); n > 0 {
			// success logged in
			// set auth cookie
			cookies.Set("name", authform.Username)
			cookies.Set("token", token)
			cookies.Save()

			c.Redirect(302, "/posts")
		} else {
			// aww no logged in
			// redirect sign in
			c.Redirect(302, "/auth/sign-in")
		}
	}

}

// GET /auth/register
func AuthRegister(c *gin.Context) {
	c.HTML(http.StatusOK, "auth/register.html", gin.H{
		"Site": cfg.Site,
	})
}

// POST /auth/register
func AuthTryRegister(c *gin.Context) {
	session := globalSession.Copy()
	s := session.DB(cfg.Database.Name).C("users")

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
				// do stuff
			}

			c.Redirect(302, "/posts")
		}
	}
}

// Create the admin user
func CreateAdmin() {
	session := globalSession.Copy()
	s := session.DB(cfg.Database.Name).C("users")

	// no existing user, good 2 go
	hash := sha512.Sum512([]byte(cfg.Server.AdminPassword))
	token := hex.EncodeToString(hash[:])

	admin := User{
		Name:        "admin",
		Token:       token,
		AccessLevel: 2,
	}

	query := bson.M{
		"name": "admin",
	}

	// write db
	if _, err := s.Upsert(&query, &admin); err != nil {
		// do stuff
	}
}
