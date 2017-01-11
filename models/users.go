// coalesce/models/users.go

package models

type SQLUser struct {
	Userid      int
	Name        string
	Salt        string
	Token       string
	AccessLevel int
}

// GetUsersAll
// GetUserByID
// GetUserByName
// GetUserExists
// WriteUser
// UpdateUser
// DeleteUser

func QueryUsersAll() ([]SQLUser, error) {
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

func QueryUser(name, token string) (SQLUser, error) {
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

func QueryUserExists(name string) (bool, error) {
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

func QueryUsername(name string) (SQLUser, error) {
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

func QueryUserID(user int) (SQLUser, error) {
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

func WriteUser(user SQLUser) error {
	s := `INSERT INTO users (userid, username, salt, token, privlevel) VALUES (?, ?, ?, ?, ?)`
	_, err := sqdb.Exec(s, nil, user.Name, user.Salt, user.Token, user.AccessLevel)
	return err
}

func UpdateUser(user SQLUser) error {
	s := `UPDATE users SET username=?,salt=?,token=?,privlevel=? WHERE userid=?`
	_, err := sqdb.Exec(s, user.Name, user.Salt, user.Token, user.AccessLevel, user.Userid)
	return err
}

func DeleteUser(user int) error {
	s := `DELETE FROM users WHERE userid=?`
	_, err := sqdb.Exec(s, user)
	return err
}
