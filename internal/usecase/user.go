package usecase

import (
	"entry_task/internal/config"
	"entry_task/internal/model/user"
	image2 "entry_task/internal/util/image"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
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
	conf := config.Get()

	var mySigningKey = []byte(conf.AuthConfig.JWTSecret)
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
	var err error
	user.Password, err = u.hashPassword(user.Password)
	if err != nil {
		return err
	}

	tx, err := u.ur.CreateTx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = u.ur.UpsertUser(user, tx)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return err
}

func (u *UseCase) AuthenticateUser(username string, password string) (interface{}, error) {
	// connect ke database
	user, err := u.ur.GetUserByName(username)
	if err != nil {
		return nil, err
	}

	if user.UserName == "" {
		return nil, errors.New("username not found")
	}

	match := u.checkPasswordHash(password, user.Password)
	if !match {
		return nil, errors.New("password incorrect")
	}

	validToken, err := u.generateJWT(user.UserName) //JWT , send UserName utk dimasukkan ke claim
	if err != nil {
		return nil, err
	}

	userData := map[string]interface{}{
		"id":          user.Id,
		"username":    user.UserName,
		"nickname":    user.NickName,
		"tokenstring": validToken,
	}
	return userData, nil
}

func (u *UseCase) UpdateUser(user user.User) error {
	tx, err := u.ur.CreateTx()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	err = u.ur.UpsertUser(user, tx)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (u *UseCase) GetUserByID(userID int) (user.User, error) {
	userData, err := u.ur.GetUserByID(userID)
	if err != nil {
		return user.User{}, nil
	}

	return userData, nil
}

func (u *UseCase) UploadUserPic(id int, username string, picData string) error {
	idString := strconv.Itoa(id)
	err := image2.Base64toPng(idString, picData)
	if err != nil {
		return err
	}

	tx, err := u.ur.CreateTx()

	defer tx.Rollback()
	picName := "Pict_" + username
	err = u.ur.UpdateUserPic(picName, id, tx)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (u *UseCase) GetUserPicByID(userID int) (string, error) {
	idString := strconv.Itoa(userID)
	picData, err := image2.Fgetbase64(idString)
	if err != nil {
		return "", err
	}

	return picData, nil
}
