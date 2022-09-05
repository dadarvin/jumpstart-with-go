package handler

import (
	"encoding/json"
	"entry_task/internal/model/user"
	"entry_task/internal/util/httputil"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"strconv"
)

// RegisterUserFunc register new user
// @Router /register [POST]
// @Param username body string
// @Param nickname body string
func (h *Handler) RegisterUserFunc() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var newUser user.User

		checkPost := httputil.CheckPostHeader(r)
		if checkPost != "" {
			httputil.ErrorResponse(w, http.StatusUnsupportedMediaType, checkPost)
			return
		}

		bodyVal, err := io.ReadAll(r.Body)
		if err != nil {
			httputil.ErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		err = json.Unmarshal(bodyVal, &newUser)
		if err != nil {
			httputil.ErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		err = h.user.RegisterUser(newUser)
		if err != nil {
			httputil.ErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		httputil.HttpResponse(w, http.StatusOK, "Register Sukses", nil)
	}
}

// LoginFunc endpoint for user login
// @Routes /login
// @Param id body int
// @Param password body string
// jika user ada, return json (id,namauser dan nickname) sebaliknya null
func (h *Handler) LoginFunc() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var user user.User

		checkPost := httputil.CheckPostHeader(r)
		if checkPost != "" {
			httputil.ErrorResponse(w, http.StatusUnsupportedMediaType, " Harus menggunakan POST")
			return
		}

		bodyVal, err := io.ReadAll(r.Body)
		if err != nil {
			httputil.ErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		err = json.Unmarshal(bodyVal, &user)
		if err != nil {
			httputil.ErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		userData, err := h.user.AuthenticateUser(user.UserName, user.Password)
		if err != nil {
			httputil.ErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		httputil.HttpResponse(w, http.StatusOK, "Login Sukses", userData)
	}
}

// EditUserFunc Edit user nickname
// @Router /edit-user
func (h *Handler) EditUserFunc() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var user user.User

		bodyVal, err := io.ReadAll(r.Body)
		if err != nil {
			httputil.ErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		err = json.Unmarshal(bodyVal, &user)
		if err != nil {
			httputil.ErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		err = h.user.UpdateUser(user)
		if err != nil {
			httputil.ErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		httputil.HttpResponse(w, http.StatusOK, "Edit User Sukses", user)
	}
}

// GetProfileFunc Get Profile info
// @Routes /user-profile/:id
// @Params id queryParam int
// jika username tdk ditemukan, return null
func (h *Handler) GetProfileFunc() httprouter.Handle {
	return func(w http.ResponseWriter, _ *http.Request, param httprouter.Params) {
		var (
			userID int
			err    error
		)

		idString := param.ByName("id")
		userID, _ = strconv.Atoi(idString)

		userData, err := h.user.GetUserByID(userID)
		if err != nil {
			httputil.ErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		userProfile := map[string]interface{}{
			"id":       userData.Id,
			"username": userData.UserName,
			"nickname": userData.NickName,
		}

		httputil.HttpResponse(w, http.StatusOK, "Sukses", userProfile)
	}
}

// ----------------------------------
// Picture Logic
// ----------------------------------

// UploadProfilePictFunc pict stlh di decode di write/dikirim ke user dgn json
// file di simpan di server dgn nama file sesuai nama user
// Converts  base64 data to   png
func (h *Handler) UploadProfilePictFunc() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var (
			userPic user.UserPicture
			err     error
		)

		checkPost := httputil.CheckPostHeader(r)
		if checkPost != "" {
			httputil.ErrorResponse(w, http.StatusUnsupportedMediaType, " Harus menggunakan POST")
			return
		}
		bodyVal, err := io.ReadAll(r.Body)
		if err != nil {
			httputil.ErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		err = json.Unmarshal(bodyVal, &userPic)
		if err != nil {
			httputil.ErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		err = h.user.UploadUserPic(userPic.Id, userPic.UserName, userPic.ProfilePicture)
		if err != nil {
			httputil.ErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		httputil.HttpResponse(w, http.StatusOK, "image uploaded", nil)
	}

}

// GetProfilePictFunc pict stlh di decode di write/dikirim ke user dgn json
// @Routes /get-profile-pict/:id
// @Params id parameter id user int
// jika username tdk ditemukan, return null
func (h *Handler) GetProfilePictFunc() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
		var (
			userID int
			err    error
		)

		idString := param.ByName("id")
		userID, _ = strconv.Atoi(idString)

		pic, err := h.user.GetUserPicByID(userID)
		if err != nil {
			httputil.ErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		httputil.HttpResponse(w, http.StatusOK, "success get user profile picture", pic)
	}
}
