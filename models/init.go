// coalesce/models/init.go

package models

import (
	"database/sql"
	"io/ioutil"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nytopop/coalesce/util"
)

var sqdb *sql.DB

func CloseDB() {
	sqdb.Close()
}

func InitDB(dbfile, initfile string) error {
	var err error
	sqdb, err = sql.Open("sqlite3", dbfile)
	if err != nil {
		return err
	}

	s, err := ioutil.ReadFile(initfile)
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

	salt, err := util.GenerateSalt()
	if err != nil {
		return err
	}

	token, err := util.ComputeToken(salt, pass)
	if err != nil {
		return err
	}

	adm := SQLUser{
		Name:        "admin",
		Salt:        salt,
		Token:       token,
		AccessLevel: 3,
	}

	admInDB, err := QueryUserExists(adm.Name)
	if err != nil {
		return err
	}

	switch admInDB {
	case true:
		err = UpdateUser(adm)
	case false:
		err = WriteUser(adm)
	}

	return err
}
