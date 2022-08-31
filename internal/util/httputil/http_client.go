package httputil

import (
	"encoding/json"
	"entry_task/pkg/dto/base"
	"github.com/golang/gddo/httputil/header"
	"net/http"
)

func HttpResponse(w http.ResponseWriter, responseCode int, message string, m interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(responseCode)

	if err := json.NewEncoder(w).Encode(&base.JsonResponse{Status: responseCode, Message: message, Data: m}); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
	}
}

func ErrorResponse(w http.ResponseWriter, errorCode int, errorMsg string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(errorCode)
	json.NewEncoder(w).Encode(&base.JsonErrorResponse{Error: &base.ApiResponse{Status: errorCode, Message: errorMsg}})
}

func CheckPostHeader(r *http.Request) string {
	msg := ""
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			msg = "Content-Type header is not application/json"
			return msg
		}
	}
	return msg
}
