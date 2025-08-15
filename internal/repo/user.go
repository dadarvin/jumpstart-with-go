package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"entry_task/internal/model/user"
	"entry_task/pkg/client/token"
	"fmt"
	"github.com/go-redis/redis/v8"
	"image"
	"image/png"
	"os"
	"time"
)

func (r *Repo) CreateTx() (*sql.Tx, error) {
	tx, err := r.db.Master.Begin()
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (r *Repo) GetJWT(username string) (string, error) {
	validToken, err := token.GenerateJWT(username)

	return validToken, err
}

// UpsertUser Update or insert user data
func (r *Repo) UpsertUser(user user.User, tx *sql.Tx) error {
	var err error
	if user.Id > 0 {
		if tx != nil {
			_, err = tx.Exec("UPDATE user SET nickname=? WHERE id=?", user.NickName, user.Id)
		} else {
			_, err = r.db.Master.Exec("UPDATE user SET nickname=? WHERE id=?", user.NickName, user.Id)
		}

	} else {
		if tx != nil {
			_, err = tx.Exec("INSERT INTO user(username, nickname, password) VALUES (?, ?, ?) ", user.UserName, user.NickName, user.Password)
		} else {
			_, err = r.db.Master.Exec("INSERT INTO user(username, nickname, password) VALUES (?, ?, ?) ", user.UserName, user.NickName, user.Password)
		}
	}

	if err != nil {
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
	var (
		err  error
		each user.User
	)

	// Get Data from Redis if exist on cache
	cacheData, err := r.cache.Redis.Get(context.Background(), "userData_"+username).Result()
	if err == nil {
		err = json.Unmarshal([]byte(cacheData), &each)
		if err != nil {
			return user.User{}, err
		}

		return each, nil
	} else if err != redis.Nil {
		// Return error when Redis operation fails for reasons other than missing key
		return user.User{}, err
	}

	// If data is not exist in cache, get from database and set cache
	rows, err := r.db.Master.Query("SELECT id, username, nickname, password FROM user WHERE username=? LIMIT 1", username)
	if err != nil {
		return user.User{}, err
	}

	for rows.Next() {
		var err = rows.Scan(&each.Id, &each.UserName, &each.NickName, &each.Password)
		if err != nil {
			return user.User{}, err
		}
	}

	json, _ := json.Marshal(each)
	ttl := time.Duration(7) * time.Second
	setStatus := r.cache.Redis.Set(context.Background(), "userData_"+username, string(json), ttl)
	if setStatus.Err() != nil {
		fmt.Errorf("errRedis %+v", err)
	}

	return each, nil
}

// UpdateUserPic update user picture
func (r *Repo) UpdateUserPic(picName string, userID int, tx *sql.Tx) error {
	var err error

	if tx != nil {
		_, err = tx.Exec("UPDATE user SET profile_picture = ? WHERE id = ?", picName, userID)
	} else {
		_, err = r.db.Master.Exec("UPDATE user SET profile_picture = ? WHERE id = ?", picName, userID)
	}

	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) UploadUserPic(image image.Image, fileName string) error {
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return err
	}

	err = png.Encode(f, image)
	if err != nil {
		return err
	}

	return nil
}
