// coalesce/users.go

package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

// TODO: should return comments and posts
// GET /users/me
func UsersMe(c *gin.Context) {
	session := globalSession.Copy()
	s := session.DB(cfg.Database.Name).C("posts")

	user := GetUser(c)

	// query for user
	query := bson.M{
		"author": user.Name,
	}

	// get posts
	posts := []*Post{}
	if err := s.Find(query).Sort("-timestamp").Iter().All(&posts); err != nil {
		c.Error(err)
		c.Redirect(302, "/error")
	}

	c.HTML(http.StatusOK, "users/me.html", gin.H{
		"Site": cfg.Site,
		"List": posts,
		"User": user,
	})
}

// GET /users/all
func UsersAll(c *gin.Context) {
	session := globalSession.Copy()
	s := session.DB(cfg.Database.Name).C("users")

	// query for all users
	users := []*User{}
	if err := s.Find(nil).Sort("name").Iter().All(&users); err != nil {
		c.Error(err)
		c.Redirect(302, "/error")
	}

	c.HTML(http.StatusOK, "users/all.html", gin.H{
		"Site":  cfg.Site,
		"Users": users,
		"User":  GetUser(c),
	})
}

// TODO promote/demote should ask to reauthenticate

// GET /users/promote/:name
func UsersTryPromote(c *gin.Context) {
	// change accesslevel to 2
}

// GET /users/demote/:name
func UsersTryDemote(c *gin.Context) {

}
