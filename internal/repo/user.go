package repo

import (
	"database/sql"
	"entry_task/internal/model/user"
	"fmt"
)

func (r *Repo) CreateTx() (*sql.Tx, error) {
	tx, err := r.db.Master.Begin()
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// UpsertUser Update or insert user data
func (r *Repo) UpsertUser(user user.User, tx *sql.Tx) error {
	var err error
	if user.Id > 0 {
		_, err = tx.Query("UPDATE user SET nickname=? WHERE id=?", user.NickName, user.Id)
	} else {
		_, err = tx.Query("INSERT INTO user(username, nickname, password) VALUES (?, ?, ?) ", user.UserName, user.NickName, user.Password)
	}

	if err != nil {
		fmt.Println("error disini")
		return err
	}

	return nil
}

// GetUserByID Get user data by ID
func (r *Repo) GetUserByID(userID int) (user.User, error) {
	var err error
	rows, err := r.db.Master.Query("SELECT id, username, nickname, password FROM user WHERE id=? LIMIT 1", userID)
	if err != nil {
		return user.User{}, err
	}

	var each user.User
	for rows.Next() {
		var err = rows.Scan(&each.Id, &each.UserName, &each.NickName, &each.Password)
		if err != nil {
			return user.User{}, err
		}
	}

	return each, nil
}

// GetUserByName Get user data by nickname
func (r *Repo) GetUserByName(username string) (user.User, error) {
	var err error
	rows, err := r.db.Master.Query("SELECT id, username, nickname, password FROM user WHERE username=? LIMIT 1", username)
	if err != nil {
		return user.User{}, err
	}

	var each user.User
	for rows.Next() {
		var err = rows.Scan(&each.Id, &each.UserName, &each.NickName, &each.Password)
		if err != nil {
			return user.User{}, err
		}
	}

	return each, nil
}

// UpdateUserPic update user picture
func (r *Repo) UpdateUserPic(picName string, userID int, tx *sql.Tx) error {
	var err error
	_, err = tx.Query("UPDATE user SET profile_picture = ? WHERE id = ?", picName, userID)

	if err != nil {
		return err
	}

	return nil
}
