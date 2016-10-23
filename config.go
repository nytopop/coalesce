// coalesce/config.go

package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

type SiteConfigForm struct {
	Name           string `bson:"name" form:"name"`
	Title          string `bson:"title" form:"title" binding:"required"`
	Description    string `bson:"description" form:"description" binding:"required"`
	Owner          string `bson:"owner" form:"owner" binding:"required"`
	Github         string `bson:"github" form:"github"`
	Email          string `bson:"email" form:"email"`
	CorticalApiKey string `bson:"corticalapikey" form:"corticalapikey"`
}

type SiteConfig struct {
	Id             bson.ObjectId `bson:"_id,omitempty"`
	Name           string        `bson:"name"`
	Title          string        `bson:"title"`
	Description    string        `bson:"description"`
	Owner          string        `bson:"owner"`
	Github         string        `bson:"github"`
	Email          string        `bson:"email"`
	CorticalApiKey string        `bson:"corticalapikey"`
}

// Get active site configuration
func GetConf() SiteConfig {
	session := globalSession.Copy()
	s := session.DB(dbname).C("conf")

	query := bson.M{
		"name": "siteconfig",
	}

	conf := SiteConfig{}
	if err := s.Find(query).One(&conf); err != nil {
		return SiteConfig{}
	}

	return conf
}

// GET /config
func ConfigEdit(c *gin.Context) {
	c.HTML(http.StatusOK, "config/edit.html", gin.H{
		"Site": GetConf(),
		"User": GetUser(c),
	})
}

// POST /config/edit
func ConfigTryEdit(c *gin.Context) {
	session := globalSession.Copy()
	s := session.DB(dbname).C("conf")

	// validate
	var confform SiteConfigForm
	if err := c.Bind(&confform); err == nil {
		confform.Name = "siteconfig"
		query := bson.M{
			"name": "siteconfig",
		}

		if _, err := s.Upsert(&query, &confform); err != nil {
			c.Error(err)
			c.Redirect(302, "/error")
		}
	}

	c.Redirect(302, "/config")
}

// POST /config/reset
func ConfigTryReset(c *gin.Context) {
	session := globalSession.Copy()
	s := session.DB(dbname).C("conf")

	conf := SiteConfig{
		Name:        "siteconfig",
		Title:       "coalesce",
		Description: "lightning fast cms",
		Owner:       "nytopop",
		Github:      "nytopop",
		Email:       "ericizoita@gmail.com",
	}

	query := bson.M{
		"name": "siteconfig",
	}

	if _, err := s.Upsert(&query, &conf); err != nil {
		c.Error(err)
		c.Redirect(302, "/error")
	}

	c.Redirect(302, "/config")
}
