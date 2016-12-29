// coalesce/db.go

package main

import (
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
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

	hash := sha512.Sum512([]byte(pass))
	token := hex.EncodeToString(hash[:])
	adm := SQLUser{
		Name:        "admin",
		Token:       token,
		AccessLevel: 3,
	}

	// check if an 'admin' user exists
	// if it does, update it to the right password
	// if not, write a new user

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
	//s := `SELECT * FROM posts`

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

func queryPostsUser(user string) ([]SQLPost, error) {
	return []SQLPost{}, nil
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
		p.Updated,
	)

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
	s := `SELECT * FROM users WHERE username = ?`

	row, err := sqdb.Query(s, name)
	if err != nil {
		return false, err
	}

	u := SQLUser{}
	for row.Next() {
		err = row.Scan(
			&u.Userid,
			&u.Name,
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

func queryUserID(name string) (int, error) {
	s := `SELECT userid FROM users WHERE username = ?`

	row, err := sqdb.Query(s, name)
	if err != nil {
		return 0, err
	}

	var id int
	for row.Next() {
		err = row.Scan(&id)
		if err != nil {
			return 0, err
		}
	}

	return id, nil
}

func writeUser(user SQLUser) error {
	s := `INSERT INTO users (userid, username, token, privlevel) VALUES (?, ?, ?, ?)`
	_, err := sqdb.Exec(s, nil, user.Name, user.Token, user.AccessLevel)
	return err
}

func updateUser(user SQLUser) error {
	s := `UPDATE users SET username=?,token=?,privlevel=? WHERE username=?`
	_, err := sqdb.Exec(s, user.Name, user.Token, user.AccessLevel, user.Name)
	return err
}
