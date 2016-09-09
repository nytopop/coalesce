// coalesce/auth.go

package main

import (
	"crypto/sha512"
	"encoding/hex"
	"net/http"
	"strconv"

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
	Id    bson.ObjectId `bson:"_id,omitempty"`
	Name  string        `bson:"name"`
	Token string        `bson:"token"`
}

// modified auth middleware
func BasicAuthForRealm(realm string) gin.HandlerFunc {
	//session := globalSession.Copy()
	//s := session.DB(cfg.Database.Name).C("auth")

	if realm == "" {
		realm = "Authorization Required"
	}

	realm = "Basic realm=" + strconv.Quote(realm)

	return func(c *gin.Context) {
		// Search user in the slice of allowed credentials
		//c.Request.Header.Get("Authorization")
		/*
			if !found {
				// Credentials doesn't match, we return 401 and abort handlers chain.
				c.Header("WWW-Authenticate", realm)
				c.AbortWithStatus(401)
			} else {
				// The user credentials was found, set user's id to key AuthUserKey in this context, the userId can be read later using
				// c.MustGet(gin.AuthUserKey)
				c.Set("user", user)
			}
		*/
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
	session := globalSession.Copy()
	s := session.DB(cfg.Database.Name).C("users")

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
			// set auth header
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

			user := User{
				Name:  authform.Username,
				Token: token,
			}

			// write db
			if err := s.Insert(&user); err != nil {
				// do stuff
			}

			c.Redirect(302, "/auth/register")
		}
	}

}
