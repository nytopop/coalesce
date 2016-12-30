// coalesce/comments.go

package main

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"gopkg.in/mgo.v2/bson"
)

type CommentForm struct {
	Postid string `form:"postid" binding:"required"`
	Body   string `form:"body" binding:"required"`
}

type CommentReplyForm struct {
	Postid    string `form:"postid" binding:"required"`
	Commentid string `form:"commentid" binding:"required"`
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

type SQLComment struct {
	Commentid int64
	Postid    int
	Parentid  sql.NullInt64
	Userid    int
	Body      string
	Posted    int64
	Updated   int64
	// Not in SQL below this comment
	Separator string
	Indent    int
	Username  string
	Nicetime  string
}

func CommentsForPost(postid int) ([]SQLComment, error) {
	raw, err := queryCommentsPost(postid)
	if err != nil {
		return []SQLComment{}, err
	}

	userCache := map[int]string{}
	for i, c := range raw {
		// query db and add to userCache if we haven't
		if _, ok := userCache[c.Userid]; !ok {
			user, err := queryUserID(c.Userid)
			if err != nil {
				return []SQLComment{}, err
			}
			userCache[user.Userid] = user.Name
		}
		raw[i].Username = userCache[c.Userid]
		raw[i].Nicetime = NiceTime(c.Posted)
	}

	tree := []SQLComment{}
	for _, c := range raw {
		if !c.Parentid.Valid {
			branch := CommentTree(c, raw)
			tree = append(tree, branch...)
		}
	}

	return tree, nil
}

func CommentTree(root SQLComment, comments []SQLComment) []SQLComment {
	root.Separator += "|"
	out := []SQLComment{root}
	for _, c := range comments {
		if c.Parentid.Int64 == root.Commentid {
			replies := CommentTree(c, comments)
			for i, _ := range replies {
				replies[i].Separator += "|"
				replies[i].Indent += 1
			}
			out = append(out, replies...)
		}
	}
	return out
}

// POST /comments/new
func CommentsTryNew(c *gin.Context) {
	var cform CommentForm
	err := c.Bind(&cform)
	if err != nil {
		RenderErr(c, err)
		return
	}

	pNum, err := strconv.Atoi(cform.Postid)
	if err != nil {
		RenderErr(c, err)
		return
	}

	user := GetUser(c)

	comment := SQLComment{
		Postid:  pNum,
		Userid:  user.Userid,
		Body:    cform.Body,
		Posted:  time.Now().Unix(),
		Updated: time.Now().Unix(),
	}
	err = writeComment(comment)
	if err != nil {
		RenderErr(c, err)
		return
	}

	posturl := "/posts/view/" + cform.Postid
	c.Redirect(302, posturl)
}

// POST /comments/reply
func CommentsTryReply(c *gin.Context) {
	var rform CommentReplyForm
	err := c.Bind(&rform)
	if err != nil {
		RenderErr(c, err)
		return
	}

	pNum, err := strconv.Atoi(rform.Postid)
	if err != nil {
		RenderErr(c, err)
		return
	}

	par, err := strconv.Atoi(rform.Commentid)
	if err != nil {
		RenderErr(c, err)
		return
	}

	var parent sql.NullInt64
	err = parent.Scan(par)
	if err != nil {
		RenderErr(c, err)
		return
	}

	user := GetUser(c)
	reply := SQLComment{
		Postid:   pNum,
		Parentid: parent,
		Userid:   user.Userid,
		Body:     rform.Body,
		Posted:   time.Now().Unix(),
		Updated:  time.Now().Unix(),
	}
	err = writeCommentReply(reply)
	if err != nil {
		RenderErr(c, err)
		return
	}

	posturl := "/posts/view/" + rform.Postid
	c.Redirect(302, posturl)
}

// TODO for tomorrow
func CommentsTryDelete(c *gin.Context) {
	// if comment has replies, we set message to <deleted>
	// if comment has no replies, we delete from db
}
