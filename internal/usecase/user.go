package usecase

import (
	"encoding/base64"
	"entry_task/internal/model/user"
	image2 "entry_task/internal/util/image"
	"entry_task/pkg/client/encrypt"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"strconv"
	"strings"
)

// Base64toPng convert base 64 to png format
func (u *UseCase) base64toPng(fIdUser string, fPicture string) error {

	//fPicture adalah base64cie yg dikirim dari clien utk diubah jadi png
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(fPicture))

	m, formatString, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
		return err
	}
	bounds := m.Bounds()
	fmt.Println(bounds, formatString)

	//Encode from image format to writer
	//fUser=nama user yg dijadinakan file name (nama user unique)
	pngFilename := "assets/Pict_" + fIdUser + ".png"

	err = u.ur.UploadUserPic(m, pngFilename)
	if err != nil {
		return err
	}
	fmt.Println("Png file", pngFilename, "created")

	return nil
}

func (u *UseCase) RegisterUser(user user.User) error {
	// hash password with bcrypt
	var err error
	user.Password, err = encrypt.HashPassword(user.Password)
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

	match := encrypt.CheckPasswordHash(password, user.Password)
	if !match {
		return nil, errors.New("password incorrect")
	}

	validToken, err := u.ur.GetJWT(user.UserName) //JWT , send UserName utk dimasukkan ke claim
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
	err := u.base64toPng(idString, picData)
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
