// coalesce/models/posts.go

package models

import (
	"html/template"

	"github.com/nytopop/coalesce/util"
)

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
	// Not in SQL below this line
	Username    string
	PostedNice  string
	UpdatedNice string
}

// ProcessPosts adds the following dynamic information to a slice of SQLPost objects: the username that posted them and the 'nice time' of their post and last updated date.
func ProcessPosts(posts []SQLPost) ([]SQLPost, error) {
	userCache := map[int]string{}
	for i, post := range posts {
		if _, ok := userCache[post.Userid]; !ok {
			user, err := QueryUserID(post.Userid)
			if err != nil {
				return []SQLPost{}, err
			}
			userCache[user.Userid] = user.Name
		}
		posts[i].Username = userCache[post.Userid]
		posts[i].PostedNice = util.NiceTime(post.Posted)
		posts[i].UpdatedNice = util.NiceTime(post.Updated)
	}
	return posts, nil
}

// QueryPostsPage returns posts in page page, of size size.
func QueryPostsPage(page, size int) ([]SQLPost, error) {
	s := `SELECT *
	FROM posts
	ORDER BY postid DESC
	LIMIT ?
	OFFSET ?`

	rows, err := sqdb.Query(s, size, page*size)
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
	s := `SELECT * 
	FROM posts 
	WHERE userid=? 
	ORDER BY postid DESC`

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
	s := `SELECT * 
	FROM posts 
	WHERE postid = ?`

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

func QueryPrevPost(post int) (int, error) {
	return 0, nil
}

func QueryNextPost(post int) (int, error) {
	return 0, nil
}

func WritePost(p SQLPost) error {
	//	s := `INSERT INTO posts VALUES (?, ?, ?, ?, ?, ?, ?, ?)` // category
	s := `INSERT INTO posts 
	VALUES (?, ?, ?, ?, ?, ?, ?)`

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
	s := `UPDATE posts 
	SET title=?,body=?,bodyHTML=?,updated=? 
	WHERE postid=?`
	_, err := sqdb.Exec(s, p.Title, p.Body, p.BodyHTML, p.Updated, p.Postid)
	return err
}

func DeletePost(post int) error {
	s := `DELETE FROM posts 
	WHERE postid=?`
	_, err := sqdb.Exec(s, post)
	return err
}
