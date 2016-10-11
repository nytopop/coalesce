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
// TODO generic user priv change function

// GET /users/promote/:name
func UsersTryPromote(c *gin.Context) {
	session := globalSession.Copy()
	s := session.DB(cfg.Database.Name).C("users")

	// query for user
	name := c.Param("name")
	query := bson.M{
		"name": name,
	}

	// get user
	user := User{}
	if err := s.Find(query).One(&user); err != nil {
		c.Error(err)
		c.Redirect(302, "/error")
	}

	// update user
	user.AccessLevel = 2
	if err := s.Update(query, user); err != nil {
		c.Error(err)
		c.Redirect(302, "/error")
	}

	c.Redirect(302, "/users/all")
}

// GET /users/demote/:name
func UsersTryDemote(c *gin.Context) {
	session := globalSession.Copy()
	s := session.DB(cfg.Database.Name).C("users")

	// query for user
	name := c.Param("name")
	query := bson.M{
		"name": name,
	}

	// get user
	user := User{}
	if err := s.Find(query).One(&user); err != nil {
		c.Error(err)
		c.Redirect(302, "/error")
	}

	// update user
	user.AccessLevel = 1
	if err := s.Update(query, user); err != nil {
		c.Error(err)
		c.Redirect(302, "/error")
	}

	c.Redirect(302, "/users/all")
}

// GET /users/del/:name
func UsersTryDelete(c *gin.Context) {
	session := globalSession.Copy()
	s := session.DB(cfg.Database.Name).C("users")

	name := c.Param("name")
	if name != "admin" {
		// query
		query := bson.M{
			"name": name,
		}

		// delete user
		if err := s.Remove(query); err != nil {
			c.Error(err)
			c.Redirect(302, "/error")
		}
	}

	c.Redirect(302, "/users/all")
}
