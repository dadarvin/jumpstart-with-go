package handler

import (
	"entry_task/cmd"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"testing"
)

func TestHandler_EditUserFunc(t *testing.T) {
	type fields struct {
		user cmd.UserUseCase
	}
	type args struct {
		w   http.ResponseWriter
		r   *http.Request
		in2 httprouter.Params
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				user: tt.fields.user,
			}
			h.EditUserFunc(tt.args.w, tt.args.r, tt.args.in2)
		})
	}
}

func TestHandler_GetProfileFunc(t *testing.T) {
	type fields struct {
		user cmd.UserUseCase
	}
	type args struct {
		w     http.ResponseWriter
		in1   *http.Request
		param httprouter.Params
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				user: tt.fields.user,
			}
			h.GetProfileFunc(tt.args.w, tt.args.in1, tt.args.param)
		})
	}
}

func TestHandler_GetProfilePictFunc(t *testing.T) {
	type fields struct {
		user cmd.UserUseCase
	}
	type args struct {
		w     http.ResponseWriter
		r     *http.Request
		param httprouter.Params
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				user: tt.fields.user,
			}
			h.GetProfilePictFunc(tt.args.w, tt.args.r, tt.args.param)
		})
	}
}

func TestHandler_LoginFunc(t *testing.T) {
	type fields struct {
		user cmd.UserUseCase
	}
	type args struct {
		w   http.ResponseWriter
		r   *http.Request
		in2 httprouter.Params
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				user: tt.fields.user,
			}
			h.LoginFunc(tt.args.w, tt.args.r, tt.args.in2)
		})
	}
}

func TestHandler_RegisterUserFunc(t *testing.T) {
	type fields struct {
		user cmd.UserUseCase
	}
	type args struct {
		w   http.ResponseWriter
		r   *http.Request
		in2 httprouter.Params
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				user: tt.fields.user,
			}
			h.RegisterUserFunc(tt.args.w, tt.args.r, tt.args.in2)
		})
	}
}

func TestHandler_UploadProfilePictFunc(t *testing.T) {
	type fields struct {
		user cmd.UserUseCase
	}
	type args struct {
		w   http.ResponseWriter
		r   *http.Request
		in2 httprouter.Params
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				user: tt.fields.user,
			}
			h.UploadProfilePictFunc(tt.args.w, tt.args.r, tt.args.in2)
		})
	}
}
