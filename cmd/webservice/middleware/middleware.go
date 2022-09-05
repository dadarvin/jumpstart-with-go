package middleware

import (
	"entry_task/internal/util/httputil"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// IsAuthorized check whether user is authorized or not
func (m *Middleware) IsAuthorized(handler httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
		if r.Header["Token"] == nil {
			httputil.ErrorResponse(w, http.StatusUnauthorized, "No Token Found")
			return
		}

		var mySigningKey = []byte(m.secret)

		token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("there was an error in parsing token")
			}
			return mySigningKey, nil
		})
		if err != nil {
			fmt.Print("this2")
			httputil.ErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			//value := fmt.Sprint(claims["username"])
			handler(w, r, param)
			return
		}

		httputil.ErrorResponse(w, http.StatusUnauthorized, "Not Authorized")
	}
}
