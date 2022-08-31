package middleware

import (
	"entry_task/internal/config"
	"entry_task/internal/util/httputil"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// IsAuthorized check whether user is authorized or not
func IsAuthorized(handler httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
		conf := config.Get()

		if r.Header["Token"] == nil {
			httputil.ErrorResponse(w, http.StatusUnauthorized, "No Token Found")
			return
		}

		var mySigningKey = []byte(conf.AuthConfig.JWTSecret)

		token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("there was an error in parsing token")
			}
			return mySigningKey, nil
		})
		if err != nil {
			httputil.ErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			value := fmt.Sprint(claims["username"])
			fmt.Print(value)
			handler(w, r, param)
			return
		}

		httputil.ErrorResponse(w, http.StatusUnauthorized, "Not Authorized")
	}
}
