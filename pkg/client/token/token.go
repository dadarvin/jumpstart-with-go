package token

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

func GenerateJWT(username string) (string, error) {
	//conf := config.Get()

	var mySigningKey = []byte("secret")
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
