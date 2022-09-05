package handler

import (
	"entry_task/internal/model/user"
	"entry_task/internal/util/testutil"
	"entry_task/pkg/dto/base"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestHandler_EditUserFunc(t *testing.T) {
	type args struct {
		path string
		body interface{}
	}

	//userData := user.User{
	//	Id:       1,
	//	UserName: "test",
	//	NickName: "testName",
	//	Password: "testPass",
	//}

	tests := []struct {
		name string
		mock func(mock *MockUserUseCase)
		args args
		want func(t *testing.T, result *base.JsonResponse, status int)
	}{
		{
			name: "Success",
			mock: func(mock *MockUserUseCase) {
				mock.EXPECT().UpdateUser(gomock.Any()).Return(nil)
			},
			args: args{
				path: "/change-nickname",
				body: &user.User{},
			},
			want: func(t *testing.T, result *base.JsonResponse, status int) {
				assert.Equal(t, http.StatusOK, status)
			},
		},
		{
			name: "Failed",
			mock: func(mock *MockUserUseCase) {
				mock.EXPECT().UpdateUser(gomock.Any()).Return(errors.New("error updating user"))
			},
			args: args{
				path: "/change-nickname",
				body: &user.User{},
			},
			want: func(t *testing.T, result *base.JsonResponse, status int) {
				assert.Equal(t, http.StatusBadRequest, status)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUser := NewMockUserUseCase(ctrl)
			h := &Handler{
				user: mockUser,
			}

			tt.mock(mockUser)

			rr := testutil.NewRequestRecorder(t,
				h.EditUserFunc(), http.MethodPut,
				tt.args.path,
				testutil.WithBody(tt.args.body),
			)

			var resp base.JsonResponse

			testutil.ParseResponse(t, rr, &resp)
			tt.want(t, &resp, rr.Code)
		})
	}
}

func TestHandler_GetProfileFunc(t *testing.T) {
	type args struct {
		path   string
		userID int
	}
	tests := []struct {
		name string
		args args
		mock func(a args, mock *MockUserUseCase)
		want func(t *testing.T, result *base.JsonResponse, status int)
	}{
		{
			name: "Success",
			args: args{
				path:   "/get-profile/%d",
				userID: 1,
			},
			mock: func(a args, mock *MockUserUseCase) {
				mock.EXPECT().GetUserByID(gomock.Any()).Return(user.User{
					Id: 0,
				}, nil)
			},
			want: func(t *testing.T, result *base.JsonResponse, status int) {
				assert.Equal(t, http.StatusOK, status)
			},
		},
		{
			name: "Failed",
			args: args{
				path:   "/get-profile/%d",
				userID: 1,
			},
			mock: func(a args, mock *MockUserUseCase) {
				mock.EXPECT().GetUserByID(gomock.Any()).Return(user.User{}, errors.New("error getting user"))
			},
			want: func(t *testing.T, result *base.JsonResponse, status int) {
				assert.Equal(t, http.StatusBadRequest, status)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUser := NewMockUserUseCase(ctrl)
			h := &Handler{
				user: mockUser,
			}

			tt.mock(tt.args, mockUser)

			fmt.Println(fmt.Sprintf(tt.args.path, tt.args.userID))
			rr := testutil.NewRequestRecorder(t,
				h.GetProfileFunc(), http.MethodGet,
				fmt.Sprintf(tt.args.path, tt.args.userID),
			)

			var resp base.JsonResponse

			testutil.ParseResponse(t, rr, &resp)
			tt.want(t, &resp, rr.Code)
		})
	}
}

//	func TestHandler_GetProfilePictFunc(t *testing.T) {
//		type fields struct {
//			user cmd.UserUseCase
//		}
//		type args struct {
//			w     http.ResponseWriter
//			r     *http.Request
//			param httprouter.Params
//		}
//		tests := []struct {
//			name   string
//			fields fields
//			args   args
//		}{
//			// TODO: Add test cases.
//		}
//		for _, tt := range tests {
//			t.Run(tt.name, func(t *testing.T) {
//				h := &Handler{
//					user: tt.fields.user,
//				}
//				h.GetProfilePictFunc(tt.args.w, tt.args.r, tt.args.param)
//			})
//		}
//	}
func TestHandler_LoginFunc(t *testing.T) {
	type args struct {
		path   string
		body   interface{}
		header map[string]string
	}
	tests := []struct {
		name string
		args args
		mock func(mock *MockUserUseCase)
		want func(t *testing.T, result *base.JsonResponse, status int)
	}{
		{
			name: "Success",
			args: args{
				path: "/login",
				body: &user.User{},
				header: map[string]string{
					"Content-Type": "application/json",
				},
			},
			mock: func(mock *MockUserUseCase) {
				mock.EXPECT().AuthenticateUser(gomock.Any(), gomock.Any()).Return(user.User{
					Id:       1,
					UserName: "test",
					NickName: "testNickname",
					Password: "testPass",
				}, nil)
			},
			want: func(t *testing.T, result *base.JsonResponse, status int) {
				assert.Equal(t, http.StatusOK, status)
			},
		},
		{
			name: "Failed Usecase",
			args: args{
				path: "/login",
				body: &user.User{},
				header: map[string]string{
					"Content-Type": "application/json",
				},
			},
			mock: func(mock *MockUserUseCase) {
				mock.EXPECT().AuthenticateUser(gomock.Any(), gomock.Any()).Return(user.User{}, errors.New("error usecase"))
			},
			want: func(t *testing.T, result *base.JsonResponse, status int) {
				assert.Equal(t, http.StatusBadRequest, status)
			},
		},
		{
			name: "Failed Content Type",
			args: args{
				path: "/login",
				body: &user.User{},
				header: map[string]string{
					"Content-Type": "aa",
				},
			},
			mock: func(mock *MockUserUseCase) {},
			want: func(t *testing.T, result *base.JsonResponse, status int) {
				assert.Equal(t, http.StatusUnsupportedMediaType, status)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUser := NewMockUserUseCase(ctrl)
			h := &Handler{
				user: mockUser,
			}

			tt.mock(mockUser)

			rr := testutil.NewRequestRecorder(t,
				h.LoginFunc(), http.MethodPut,
				tt.args.path,
				testutil.WithBody(tt.args.body),
				testutil.WithRequestHeader(tt.args.header),
			)

			var resp base.JsonResponse

			testutil.ParseResponse(t, rr, &resp)
			tt.want(t, &resp, rr.Code)
		})
	}
}

func TestHandler_RegisterUserFunc(t *testing.T) {
	type args struct {
		path string
		body interface{}
	}
	tests := []struct {
		name string
		args args
		mock func(mock *MockUserUseCase)
		want func(t *testing.T, result *base.JsonResponse, status int)
	}{
		{
			name: "Success",
			args: args{
				path: "/register",
				body: &user.User{},
			},
			mock: func(mock *MockUserUseCase) {
				mock.EXPECT().RegisterUser(gomock.Any()).Return(nil)
			},
			want: func(t *testing.T, result *base.JsonResponse, status int) {
				assert.Equal(t, http.StatusOK, status)
			},
		},
		{
			name: "Failed Usecase",
			args: args{
				path: "/register",
				body: &user.User{},
			},
			mock: func(mock *MockUserUseCase) {
				mock.EXPECT().RegisterUser(gomock.Any()).Return(errors.New("error usecase"))
			},
			want: func(t *testing.T, result *base.JsonResponse, status int) {
				assert.Equal(t, http.StatusBadRequest, status)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUser := NewMockUserUseCase(ctrl)
			h := &Handler{
				user: mockUser,
			}

			tt.mock(mockUser)

			rr := testutil.NewRequestRecorder(t,
				h.RegisterUserFunc(), http.MethodPut,
				tt.args.path,
				testutil.WithBody(tt.args.body),
			)

			var resp base.JsonResponse

			testutil.ParseResponse(t, rr, &resp)
			tt.want(t, &resp, rr.Code)
		})
	}
}

func TestHandler_UploadProfilePictFunc(t *testing.T) {
	type args struct {
		path string
		body interface{}
	}

	tests := []struct {
		name string
		args args
		mock func(mock *MockUserUseCase)
		want func(t *testing.T, result *base.JsonResponse, status int)
	}{
		{
			name: "Success",
			args: args{
				path: "/uploadprofilepict",
				body: &user.UserPicture{},
			},
			mock: func(mock *MockUserUseCase) {
				mock.EXPECT().UploadUserPic(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			want: func(t *testing.T, result *base.JsonResponse, status int) {
				assert.Equal(t, http.StatusOK, status)
			},
		},
		{
			name: "Failed Usecase",
			args: args{
				path: "/uploadprofilepict",
				body: &user.UserPicture{},
			},
			mock: func(mock *MockUserUseCase) {
				mock.EXPECT().UploadUserPic(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("error usecase"))
			},
			want: func(t *testing.T, result *base.JsonResponse, status int) {
				assert.Equal(t, http.StatusBadRequest, status)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUser := NewMockUserUseCase(ctrl)
			h := &Handler{
				user: mockUser,
			}

			tt.mock(mockUser)

			rr := testutil.NewRequestRecorder(t,
				h.UploadProfilePictFunc(), http.MethodPut,
				tt.args.path,
				testutil.WithBody(tt.args.body),
			)

			var resp base.JsonResponse

			testutil.ParseResponse(t, rr, &resp)
			tt.want(t, &resp, rr.Code)
		})
	}
}
