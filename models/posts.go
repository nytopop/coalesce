// coalesce/models/posts.go

package models

import "html/template"

type SQLPost struct {
	Postid     int
	Userid     int
	Title      string
	Body       string
	BodyHTML   string
	RenderHTML template.HTML
	Categoryid int
	Posted     int64
	Updated    int64
}

func QueryPostsPage(page int) ([]SQLPost, error) {
	s := `SELECT * FROM posts WHERE postid > ? ORDER BY postid DESC LIMIT 5`

	rows, err := sqdb.Query(s, page*5)
	if err != nil {
		return []SQLPost{}, err
	}
	defer rows.Close()

	posts := []SQLPost{}
	for rows.Next() {
		p := SQLPost{}
		err = rows.Scan(
			&p.Postid,
			&p.Userid,
			&p.Title,
			&p.Body,
			&p.BodyHTML,
			//&p.Categoryid,
			&p.Posted,
			&p.Updated)

		if err != nil {
			return []SQLPost{}, err
		}

		posts = append(posts, p)
	}

	return posts, nil
}

func QueryPostsUserID(user int) ([]SQLPost, error) {
	s := `SELECT * FROM posts WHERE userid=? ORDER BY postid DESC`

	rows, err := sqdb.Query(s, user)
	if err != nil {
		return []SQLPost{}, err
	}

	posts := []SQLPost{}
	for rows.Next() {
		p := SQLPost{}
		err = rows.Scan(
			&p.Postid,
			&p.Userid,
			&p.Title,
			&p.Body,
			&p.BodyHTML,
			&p.Posted,
			&p.Updated)
		if err != nil {
			return []SQLPost{}, err
		}
		posts = append(posts, p)
	}

	return posts, nil
}

func QueryPost(post int) (SQLPost, error) {
	s := `SELECT * FROM posts WHERE postid = ?`

	row, err := sqdb.Query(s, post)
	if err != nil {
		return SQLPost{}, err
	}

	p := SQLPost{}
	for row.Next() {
		err = row.Scan(
			&p.Postid,
			&p.Userid,
			&p.Title,
			&p.Body,
			&p.BodyHTML,
			//&p.Categoryid,
			&p.Posted,
			&p.Updated)
		if err != nil {
			return SQLPost{}, err
		}
	}

	return p, nil
}

func WritePost(p SQLPost) error {
	//	s := `INSERT INTO posts VALUES (?, ?, ?, ?, ?, ?, ?, ?)` // category
	s := `INSERT INTO posts VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err := sqdb.Exec(s, nil,
		p.Userid,
		p.Title,
		p.Body,
		p.BodyHTML,
		// p.Categoryid,
		p.Posted,
		p.Updated)

	return err
}

func UpdatePost(p SQLPost) error {
	s := `UPDATE posts SET title=?,body=?,bodyHTML=?,updated=? WHERE postid=?`
	_, err := sqdb.Exec(s, p.Title, p.Body, p.BodyHTML, p.Updated, p.Postid)
	return err
}

func DeletePost(post int) error {
	s := `DELETE FROM posts WHERE postid=?`
	_, err := sqdb.Exec(s, post)
	return err
}
