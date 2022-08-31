package repo

import "entry_task/internal/model/user"

// UpsertUser Update or insert user data
func (r *Repo) UpsertUser(user user.User) error {
	var err error
	if user.Id > 0 {
		_, err = r.db.Master.Query("UPDATE user SET nickname=? WHERE id=?", user.NickName, user.Id)
		_, err = r.db.Slave.Query("UPDATE user SET nickname=? WHERE id=?", user.NickName, user.Id)
	} else {
		_, err = r.db.Master.Query("INSERT INTO user(username, nickname, password) VALUES (?, ?, ?) ", user.UserName, user.NickName, user.Password)
		_, err = r.db.Slave.Query("INSERT INTO user(username, nickname, password) VALUES (?, ?, ?) ", user.UserName, user.NickName, user.Password)
	}

	if err != nil {
		return err
	}

	return nil
}

// GetUser Get user data by nickname
func (r *Repo) GetUser(userID int) (user.User, error) {
	var err error
	rows, err := r.db.Slave.Query("SELECT id, username, nickname, password FROM user WHERE id=? LIMIT 1", userID)
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

func (r *Repo) UpdateUserPic(picName string, userID int) error {
	var err error
	_, err = r.db.Master.Query("UPDATE user SET profile_picture = ? WHERE id = ?", picName, userID)
	_, err = r.db.Slave.Query("UPDATE user SET profile_picture = ? WHERE id = ?", picName, userID)

	if err != nil {
		return err
	}

	return nil
}
