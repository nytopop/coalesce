// coalesce/models/comments.go

package models

import (
	"database/sql"
	"fmt"
)

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

func QueryCommentsPost(post int) ([]SQLComment, error) {
	s := `SELECT * FROM comments WHERE postid=? ORDER BY commentid DESC`

	rows, err := sqdb.Query(s, post)
	if err != nil {
		return []SQLComment{}, err
	}

	comments := []SQLComment{}
	for rows.Next() {
		c := SQLComment{}
		err = rows.Scan(
			&c.Commentid,
			&c.Postid,
			&c.Parentid,
			&c.Userid,
			&c.Body,
			&c.Posted,
			&c.Updated)
		if err != nil {
			return []SQLComment{}, err
		}
		comments = append(comments, c)
	}

	return comments, nil
}

func WriteComment(c SQLComment) error {
	s := `INSERT INTO comments VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := sqdb.Exec(s, nil,
		c.Postid,
		nil,
		c.Userid,
		c.Body,
		c.Posted,
		c.Updated)
	return err
}

func WriteCommentReply(c SQLComment) error {
	s := `INSERT INTO comments VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := sqdb.Exec(s, nil,
		c.Postid,
		c.Parentid.Int64,
		c.Userid,
		c.Body,
		c.Posted,
		c.Updated)
	return err
}

func DeleteCommentID(c SQLComment) error {
	s := `DELETE FROM comments WHERE commentid=?`
	fmt.Println(s)
	return nil
}
