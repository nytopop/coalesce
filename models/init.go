// coalesce/models/init.go

package models

import (
	"database/sql"
	"io/ioutil"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nytopop/coalesce/util"
)

var sqdb *sql.DB

func init() {
	var err error
	sqdb, err = sql.Open("sqlite3", "coalesce.db")
	if err != nil {
		log.Fatal(err)
	}
	//defer sqdb.Close()

	s, err := ioutil.ReadFile("./resources/init.sql")
	if err != nil {
		log.Fatal(err)
	}

	if _, err := sqdb.Exec(string(s)); err != nil {
		log.Fatal(err)
	}

	pass := os.Getenv("ADMIN_PASS")
	if pass == "" {
		pass = "password"
	}

	salt, err := util.GenerateSalt()
	if err != nil {
		log.Fatal(err)
	}

	token, err := util.ComputeToken(salt, pass)
	if err != nil {
		log.Fatal(err)
	}

	adm := SQLUser{
		Name:        "admin",
		Salt:        salt,
		Token:       token,
		AccessLevel: 3,
	}

	admInDB, err := QueryUserExists(adm.Name)
	if err != nil {
		log.Fatal(err)
	}

	switch admInDB {
	case true:
		err = UpdateUser(adm)
	case false:
		err = WriteUser(adm)
	}

	if err != nil {
		log.Fatal(err)
	}
}
