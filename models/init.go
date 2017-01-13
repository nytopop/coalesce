// coalesce/models/init.go

package models

import (
	"database/sql"
	"io/ioutil"

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

	// Create admin user with default credentials if it doesn't exist
	admInDB, err := QueryUserExists("admin")
	if err != nil {
		return err
	}

	if !admInDB {
		salt, err := util.GenerateSalt()
		if err != nil {
			return err
		}

		token, err := util.ComputeToken(salt, "coalesce")
		if err != nil {
			return err
		}

		adm := SQLUser{
			Name:        "admin",
			Salt:        salt,
			Token:       token,
			AccessLevel: 3,
		}

		err = WriteUser(adm)
	}

	// we return err whether or not nil
	return err
}
