// coalesce/db.go

package main

import (
	"database/sql"
	"io/ioutil"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var sqdb *sql.DB

func initDB() error {
	s, err := ioutil.ReadFile("./resources/init.sql")
	if err != nil {
		return err
	}

	if _, err := sqdb.Exec(string(s)); err != nil {
		return err
	}

	pass := os.Getenv("ADMIN_PASS")
	if pass == "" {
		pass = "password"
	}

	salt, err := GenerateSalt()
	if err != nil {
		return err
	}

	pepper, err := GeneratePepper()
	if err != nil {
		return err
	}

	token := ComputeToken(salt, pepper, pass)
	adm := SQLUser{
		Name:        "admin",
		Salt:        salt,
		Token:       token,
		AccessLevel: 3,
	}

	admInDB, err := queryUserExists(adm.Name)
	if err != nil {
		return err
	}

	switch admInDB {
	case true:
		err = updateUser(adm)
	case false:
		err = writeUser(adm)
	}

	return err
}

func queryPostsPage(page int) ([]SQLPost, error) {
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

func queryPostsUserID(user int) ([]SQLPost, error) {
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

func queryPost(post int) (SQLPost, error) {
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

func writePost(p SQLPost) error {
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

func updatePost(p SQLPost) error {
	s := `UPDATE posts SET title=?,body=?,bodyHTML=?,updated=? WHERE postid=?`
	_, err := sqdb.Exec(s, p.Title, p.Body, p.BodyHTML, p.Updated, p.Postid)
	return err
}

func deletePost(post int) error {
	s := `DELETE FROM posts WHERE postid=?`
	_, err := sqdb.Exec(s, post)
	return err
}

func queryUsersAll() ([]SQLUser, error) {
	s := `SELECT * FROM users WHERE username!="admin" ORDER BY userid`

	rows, err := sqdb.Query(s)
	if err != nil {
		return []SQLUser{}, err
	}

	users := []SQLUser{}
	for rows.Next() {
		u := SQLUser{}
		err = rows.Scan(
			&u.Userid,
			&u.Name,
			&u.Salt,
			&u.Token,
			&u.AccessLevel,
		)
		if err != nil {
			return []SQLUser{}, err
		}
		users = append(users, u)
	}

	return users, nil
}

func queryUser(name, token string) (SQLUser, error) {
	s := `SELECT * FROM users WHERE username = ? AND token = ?`

	row, err := sqdb.Query(s, name, token)
	if err != nil {
		return SQLUser{}, err
	}

	u := SQLUser{}
	for row.Next() {
		err = row.Scan(
			&u.Userid,
			&u.Name,
			&u.Salt,
			&u.Token,
			&u.AccessLevel,
		)
		if err != nil {
			return SQLUser{}, err
		}
	}

	return u, nil
}

func queryUserExists(name string) (bool, error) {
	s := `SELECT * FROM users WHERE username=?`

	row, err := sqdb.Query(s, name)
	if err != nil {
		return false, err
	}

	u := SQLUser{}
	for row.Next() {
		err = row.Scan(
			&u.Userid,
			&u.Name,
			&u.Salt,
			&u.Token,
			&u.AccessLevel,
		)
		if err != nil {
			return false, err
		}
	}

	if u == (SQLUser{}) {
		return false, nil
	} else {
		return true, nil
	}
}

func queryUsername(name string) (SQLUser, error) {
	s := `SELECT * FROM users WHERE username=?`

	row, err := sqdb.Query(s, name)
	if err != nil {
		return SQLUser{}, err
	}

	u := SQLUser{}
	for row.Next() {
		err = row.Scan(
			&u.Userid,
			&u.Name,
			&u.Salt,
			&u.Token,
			&u.AccessLevel)
		if err != nil {
			return SQLUser{}, err
		}
	}

	return u, nil
}

func queryUserID(user int) (SQLUser, error) {
	s := `SELECT * FROM users WHERE userid=?`

	row, err := sqdb.Query(s, user)
	if err != nil {
		return SQLUser{}, err
	}

	u := SQLUser{}
	for row.Next() {
		err = row.Scan(
			&u.Userid,
			&u.Name,
			&u.Salt,
			&u.Token,
			&u.AccessLevel,
		)
		if err != nil {
			return SQLUser{}, err
		}
	}

	return u, nil
}

func writeUser(user SQLUser) error {
	s := `INSERT INTO users (userid, username, salt, token, privlevel) VALUES (?, ?, ?, ?, ?)`
	_, err := sqdb.Exec(s, nil, user.Name, user.Salt, user.Token, user.AccessLevel)
	return err
}

func updateUser(user SQLUser) error {
	s := `UPDATE users SET username=?,salt=?,token=?,privlevel=? WHERE userid=?`
	_, err := sqdb.Exec(s, user.Name, user.Salt, user.Token, user.AccessLevel, user.Userid)
	return err
}

func deleteUser(user int) error {
	s := `DELETE FROM users WHERE userid=?`
	_, err := sqdb.Exec(s, user)
	return err
}
