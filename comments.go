// coalesce/comments.go

package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"

	"gopkg.in/mgo.v2/bson"
)

type CommentForm struct {
	PostId string `form:"postid" binding:"required"`
	Body   string `form:"body" binding:"required"`
}

type CommentReplyForm struct {
	PostId    string `form:"postid" binding:"required"`
	CommentId string `form:"commentid" binding:"required"`
	Body      string `form:"body" binding:"required"`
}

type Comment struct {
	Id        bson.ObjectId   `bson:"_id,omitempty"`
	PostId    bson.ObjectId   `bson:"postid"`
	Author    string          `bson:"author"`
	Timestamp time.Time       `bson:"timestamp"`
	Body      string          `bson:"body"`
	Depth     int             `bson:"depth"`
	Replies   []bson.ObjectId `bson:"replies"`
}

// returns length of indentation
func (c Comment) Indent() int {
	return c.Depth * 10
}

// recursive comment chain flattener, single tree
func (c Comment) Tree() []Comment {
	// load db session
	session := globalSession.Copy()
	s := session.DB(dbname).C("comments")

	// add self to result set
	tree := []Comment{}
	tree = append(tree, c)

	if len(c.Replies) > 0 {
		// if there is a deeper level, keep going
		for _, v := range c.Replies {
			// go deeper
			next := Comment{}

			// BUG: does this do 1 query for _every_ comment????
			// get the referenced comment
			if err := s.FindId(v).One(&next); err != nil {
				log.Println(err)
			}

			// recurse!
			tree = append(tree, next.Tree()...)
		}
	}
	return tree
}

// POST /comments/new
func CommentsTryNew(c *gin.Context) {
	session := globalSession.Copy()
	s := session.DB(dbname).C("comments")

	// validate
	var cform CommentForm
	if err := c.Bind(&cform); err == nil {
		// get obj id from hex
		hexid := cform.PostId
		id := bson.ObjectIdHex(hexid)

		comment := Comment{
			PostId:    id,
			Author:    c.MustGet("name").(string),
			Body:      cform.Body,
			Timestamp: time.Now(),
			Depth:     0,
		}
		if err := s.Insert(&comment); err != nil {
			c.Error(err)
			c.Redirect(302, "/error")
		}

		// redirect to orig post
		redir := "/posts/view/" + hexid
		c.Redirect(302, redir)
	} else {
		c.Error(err)
		c.Redirect(302, "/error")
	}
}

// TODO non-working for now, don't touch
// POST /comments/reply
func CommentsTryReply(c *gin.Context) {
	session := globalSession.Copy()
	s := session.DB(dbname).C("comments")

	// validate
	var rform CommentReplyForm
	if err := c.Bind(&rform); err == nil {
		// get parent id
		parhexid := rform.CommentId
		parid := bson.ObjectIdHex(parhexid)

		// get the parent, update replies list, write parent
		// move this to after comment exists
		parent := Comment{}
		if err := s.FindId(parid).One(&parent); err != nil {
			log.Println(err)
		}

		// get post id
		posthexid := rform.PostId
		postid := bson.ObjectIdHex(posthexid)

		comment := Comment{
			//Author:    rform.Author,
			PostId:    postid,
			Author:    c.MustGet("name").(string),
			Body:      rform.Body,
			Timestamp: time.Now(),
			Depth:     parent.Depth + 1,
		}
		if info, err := s.Upsert(comment, comment); err == nil {
			parent.Replies = append(parent.Replies, info.UpsertedId.(bson.ObjectId))
		} else {
			// error
		}

		// write updated parent
		if err := s.UpdateId(parid, parent); err != nil {
			// err
		}
	}
}
