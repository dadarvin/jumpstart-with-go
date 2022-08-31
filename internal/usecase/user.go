package usecase

import (
	"entry_task/internal/model/user"
	"entry_task/internal/util/httputil"
	image2 "entry_task/internal/util/image"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"time"
)

func (u *UseCase) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (u *UseCase) checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (u *UseCase) generateJWT(username string) (string, error) {
	var mySigningKey = []byte(secretkey)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Minute * 3000000).Unix()

	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		fmt.Println("Something went Wrong: %s" + err.Error())
		return "", err
	}

	return tokenString, nil //JWT
}

func (u *UseCase) RegisterUser(user user.User) error {
	// hash password with bcrypt
	hashedPass, err := u.hashPassword(user.Password)
	if err != nil {
		return err
	}

	// connecting to database
	db, err := DBConn()
	if err != nil {
		return err
	}
	defer db.Close()

	// insert ke database (username uniq key, id sequance )
	_, err = db.Exec("INSERT INTO user(username, nickname, password) VALUES (?, ?, ?) ", user.UserName, user.NickName, hashedPass)
	if err != nil {
		return err
	}

	return err
}

func (u *UseCase) AuthenticateUser(username string, password string) (interface{}, error) {
	// connect ke database
	db, err := DBConn()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, username, nickname, password FROM user WHERE username=? LIMIT 1", username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var each user.User
	for rows.Next() {
		var err = rows.Scan(&each.Id, &each.UserName, &each.NickName, &each.Password)
		if err != nil {
			return nil, err
		}

		// result = append(result, each)
	}

	if each.UserName == "" {
		return nil, errors.New("username not found")
	}

	match := u.checkPasswordHash(password, each.Password)
	if !match {
		return nil, errors.New("password incorrect")
	}

	validToken, err := u.generateJWT(each.UserName) //JWT , send UserName utk dimasukkan ke claim
	if err != nil {
		return nil, err
	}

	userData := map[string]interface{}{
		"id":          each.Id,
		"username":    each.UserName,
		"nickname":    each.NickName,
		"tokenstring": validToken,
	}
	return userData, nil
}

func (u *UseCase) UpdateUser(user user.User) error {
	db, err := DBConn()
	if err != nil {
		httputil.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return err
	}
	defer db.Close()

	_, err = db.Exec("UPDATE user SET nickname=? WHERE id=?", user.NickName, user.Id)
	if err != nil {
		return err
	}

	return err
}

func (u *UseCase) GetUserByID(userID int) (user.User, error) {
	db, err := DBConn()
	if err != nil {
		return user.User{}, err
	}
	defer db.Close()

	//????		rows, err := db.Query("SELECT * FROM Employee WHERE id=? LIMIT 1", "a")
	rows, err := db.Query("SELECT id, username, nickname FROM user WHERE id=? LIMIT 1", userID)
	if err != nil {
		return user.User{}, err
	}
	defer rows.Close()

	var userData user.User

	for rows.Next() {
		var err = rows.Scan(&userData.Id, &userData.UserName, &userData.NickName)
		if err != nil {
			return user.User{}, err
		}
	}

	return userData, nil
}

func (u *UseCase) UploadUserPic(id int, username string, picData string) error {
	idString := strconv.Itoa(id)
	err := image2.Base64toPng(idString, picData)
	if err != nil {
		return err
	}

	db, err := DBConn()
	if err != nil {
		return err
	}
	defer db.Close()

	picName := "Pict_" + username
	rows, err := db.Query("UPDATE user SET profile_picture = ? WHERE id = ?", picName, id)
	if err != nil {
		return err
	}
	defer rows.Close()

	return err
}

func (u *UseCase) GetUserPicByID(userID int) (string, error) {
	idString := strconv.Itoa(userID)
	picData, err := image2.Fgetbase64(idString)
	if err != nil {
		return "", err
	}

	return picData, nil
}
