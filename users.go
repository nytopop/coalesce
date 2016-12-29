// coalesce/users.go

package main

import "github.com/gin-gonic/gin"

// GET /users/all
func UsersAll(c *gin.Context) {
	users, err := queryUsersAll()
	if err != nil {
		c.Error(err)
		c.Redirect(302, "/error")
	}

	c.HTML(200, "users/all.html", gin.H{
		"Users": users,
		"User":  GetUser(c),
	})

	/*session := globalSession.Copy()
	s := session.DB(dbname).C("users")

	// query for all users
	users := []*User{}
	if err := s.Find(nil).Sort("name").Iter().All(&users); err != nil {
		c.Error(err)
		c.Redirect(302, "/error")
	}

	c.HTML(http.StatusOK, "users/all.html", gin.H{
		"Site":  GetConf(),
		"Users": users,
		"User":  GetUser(c),
	})*/
}

// GET /users/promote/:name
func UsersTryPromote(c *gin.Context) {
	/*session := globalSession.Copy()
	s := session.DB(dbname).C("users")

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

	c.Redirect(302, "/users/all")*/
}

// GET /users/demote/:name
func UsersTryDemote(c *gin.Context) {
	/*session := globalSession.Copy()
	s := session.DB(dbname).C("users")

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

	c.Redirect(302, "/users/all")*/
}

// GET /users/del/:name
func UsersTryDelete(c *gin.Context) {
	/*session := globalSession.Copy()
	s := session.DB(dbname).C("users")

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

	c.Redirect(302, "/users/all")*/
}
