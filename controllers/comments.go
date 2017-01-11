// coalesce/controllers/comments.go

package controllers

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nytopop/coalesce/models"
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

func CommentsForPost(postid int) ([]models.SQLComment, error) {
	raw, err := models.QueryCommentsPost(postid)
	if err != nil {
		return []models.SQLComment{}, err
	}

	userCache := map[int]string{}
	for i, c := range raw {
		// query db and add to userCache if we haven't
		if _, ok := userCache[c.Userid]; !ok {
			user, err := models.QueryUserID(c.Userid)
			if err != nil {
				return []models.SQLComment{}, err
			}
			userCache[user.Userid] = user.Name
		}
		raw[i].Username = userCache[c.Userid]
		raw[i].Nicetime = NiceTime(c.Posted)
	}

	tree := []models.SQLComment{}
	for _, c := range raw {
		if !c.Parentid.Valid {
			branch := CommentTree(c, raw)
			tree = append(tree, branch...)
		}
	}

	return tree, nil
}

func CommentTree(root models.SQLComment, comments []models.SQLComment) []models.SQLComment {
	root.Separator += "|"
	out := []models.SQLComment{root}
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

	comment := models.SQLComment{
		Postid:  pNum,
		Userid:  user.Userid,
		Body:    cform.Body,
		Posted:  time.Now().Unix(),
		Updated: time.Now().Unix(),
	}
	err = models.WriteComment(comment)
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
	reply := models.SQLComment{
		Postid:   pNum,
		Parentid: parent,
		Userid:   user.Userid,
		Body:     rform.Body,
		Posted:   time.Now().Unix(),
		Updated:  time.Now().Unix(),
	}
	err = models.WriteCommentReply(reply)
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
