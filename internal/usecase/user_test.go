package usecase

import (
	"entry_task/internal/model/user"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"reflect"
	"testing"
)

func TestUseCase_AuthenticateUser(t *testing.T) {
	type args struct {
		username string
		password string
	}

	userData := user.User{
		Id:       1,
		UserName: "test",
		NickName: "testName",
		Password: "testPass",
	}

	tests := []struct {
		name    string
		args    args
		mock    func(a args, mock *MockuserRepo)
		want    interface{}
		wantErr bool
	}{
		{
			name: "Success Authenticate",
			args: args{
				username: userData.UserName,
				password: userData.Password,
			},
			mock: func(a args, mock *MockuserRepo) {
				mock.EXPECT().GetUserByName(a.username).
					Return(userData, nil)
			},
			want:    userData,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := NewMockuserRepo(ctrl)
			u := &UseCase{
				ur: mockRepo,
			}

			tt.mock(tt.args, mockRepo)
			got, err := u.AuthenticateUser(tt.args.username, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthenticateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AuthenticateUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCase_RegisterUser(t *testing.T) {
	type args struct {
		user user.User
	}

	userData := user.User{
		Id:       1,
		UserName: "test",
		NickName: "testName",
		Password: "testPass",
	}

	tests := []struct {
		name    string
		args    args
		mock    func(a args, mock *MockuserRepo)
		wantErr bool
	}{
		{
			name: "Success Register",
			args: args{
				user: userData,
			},
			mock: func(a args, mock *MockuserRepo) {
				db, dbMock, _ := sqlmock.New()
				dbMock.ExpectBegin()
				tx, _ := db.Begin()
				dbMock.ExpectCommit()

				bytes, _ := bcrypt.GenerateFromPassword([]byte(userData.Password), 14)
				password := string(bytes)

				mock.EXPECT().CreateTx().Return(tx, nil)
				mock.EXPECT().UpsertUser(user.User{
					Id:       userData.Id,
					UserName: userData.UserName,
					NickName: userData.NickName,
					Password: password,
				}, nil).Return(nil)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := NewMockuserRepo(ctrl)
			u := &UseCase{
				ur: mockRepo,
			}

			tt.mock(tt.args, mockRepo)
			if err := u.RegisterUser(tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("RegisterUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCase_GetUserPicByID(t *testing.T) {
	type args struct {
		userID int
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				userID: 1,
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UseCase{}
			got, err := u.GetUserPicByID(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserPicByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetUserPicByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCase_GetUserByID(t *testing.T) {
	type args struct {
		userID int
	}

	userData := user.User{
		Id:       1,
		UserName: "test",
		NickName: "testName",
		Password: "testPass",
	}

	tests := []struct {
		name    string
		args    args
		mock    func(a args, mock *MockuserRepo)
		want    user.User
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				userID: 1,
			},
			mock: func(a args, mock *MockuserRepo) {
				mock.EXPECT().GetUserByID(a.userID).Return(userData, nil)
			},
			want:    userData,
			wantErr: false,
		},
		{
			name: "Failed",
			args: args{
				userID: 0,
			},
			mock: func(a args, mock *MockuserRepo) {
				mock.EXPECT().GetUserByID(a.userID).Return(user.User{}, errors.New("an error occured"))
			},
			want:    user.User{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := NewMockuserRepo(ctrl)
			u := &UseCase{
				ur: mockRepo,
			}

			tt.mock(tt.args, mockRepo)
			got, err := u.GetUserByID(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCase_UpdateUser(t *testing.T) {
	type args struct {
		user user.User
	}

	userData := user.User{
		Id:       1,
		UserName: "test",
		NickName: "testName",
		Password: "testPass",
	}

	tests := []struct {
		name    string
		args    args
		mock    func(a args, mock *MockuserRepo)
		wantErr bool
	}{
		{
			name: "Success UpdateUser",
			args: args{
				user: userData,
			},
			mock: func(a args, mock *MockuserRepo) {
				db, dbMock, _ := sqlmock.New()
				dbMock.ExpectBegin()
				tx, _ := db.Begin()
				dbMock.ExpectCommit()

				mock.EXPECT().CreateTx().Return(tx, nil)
				mock.EXPECT().UpsertUser(gomock.Any(), tx).Return(nil)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := NewMockuserRepo(ctrl)
			u := &UseCase{
				ur: mockRepo,
			}

			tt.mock(tt.args, mockRepo)
			if err := u.UpdateUser(tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCase_UploadUserPic(t *testing.T) {
	type args struct {
		id       int
		username string
		picData  string
	}

	picData := "\\/9j\\/4AAQSkZJRgABAQIAHAAcAAD\\/2wBDABALDA4MChAODQ4SERATGCgaGBYWGDEjJR0oOjM9PDkzODdA\\r\\nSFxOQERXRTc4UG1RV19iZ2hnPk1xeXBkeFxlZ2P\\/2wBDARESEhgVGC8aGi9jQjhCY2NjY2NjY2NjY2Nj\\r\\nY2NjY2NjY2NjY2NjY2NjY2NjY2NjY2NjY2NjY2NjY2NjY2NjY2P\\/wAARCABnAJYDASIAAhEBAxEB\\/8QA\\r\\nHwAAAQUBAQEBAQEAAAAAAAAAAAECAwQFBgcICQoL\\/8QAtRAAAgEDAwIEAwUFBAQAAAF9AQIDAAQRBRIh\\r\\nMUEGE1FhByJxFDKBkaEII0KxwRVS0fAkM2JyggkKFhcYGRolJicoKSo0NTY3ODk6Q0RFRkdISUpTVFVW\\r\\nV1hZWmNkZWZnaGlqc3R1dnd4eXqDhIWGh4iJipKTlJWWl5iZmqKjpKWmp6ipqrKztLW2t7i5usLDxMXG\\r\\nx8jJytLT1NXW19jZ2uHi4+Tl5ufo6erx8vP09fb3+Pn6\\/8QAHwEAAwEBAQEBAQEBAQAAAAAAAAECAwQF\\r\\nBgcICQoL\\/8QAtREAAgECBAQDBAcFBAQAAQJ3AAECAxEEBSExBhJBUQdhcRMiMoEIFEKRobHBCSMzUvAV\\r\\nYnLRChYkNOEl8RcYGRomJygpKjU2Nzg5OkNERUZHSElKU1RVVldYWVpjZGVmZ2hpanN0dXZ3eHl6goOE\\r\\nhYaHiImKkpOUlZaXmJmaoqOkpaanqKmqsrO0tba3uLm6wsPExcbHyMnK0tPU1dbX2Nna4uPk5ebn6Onq\\r\\n8vP09fb3+Pn6\\/9oADAMBAAIRAxEAPwDlwKMD0pwzSiuK57QzGDxS7D6in8Y5ximnAPUfSlcq4m3ilUYp\\r\\n2OKXHvRcVxnTtS7c07HNFK4DQPakC4PNOA+tOx70XAjK\\/So5gBGP94fzqfvUVx\\/qxx\\/EP51UXqRP4WSE\\r\\ncmgjilP3jSEZqS0IO\\/NGDnpUiocDg\\/McDjvV6HTPOdVWYgsM5KcfzzQ2JySM2jp6VYu7SWzmMUwG4cgj\\r\\nkMPUVBjjtTGtRu0Zopw+lFFxhinrGzuqqMsxAA9yaXFSRv5cqSEcIwYj6GpuZ30O30fSLKzhUpbpNMv3\\r\\n5XGTn29BV28jt7pPLuIVljPBBFVreYx+VbqAjycgt3x14zRcNOxGyVFHQkIc\\/wA61exyKLbuzjdZ046d\\r\\nftEuTEw3Rk9SPT8P8Kpbea3tchbyVae4JkjbbGpGdwOM89Af6ViFTWUtGdcXoM2+woK1JtpNtTcoZt+l\\r\\nJt7ZqTbRtouFyPFRXI\\/c9D94fzqzioLsfuD\\/ALw\\/nVReqIn8LJCOTSY+tSMOTmkIpXLRu+F0t5pJxPHG\\r\\nwjjUAuBjJJz1+laD6Pai+WaK9SBX6puzn6ZP+NV\\/Dkdtc6ZNbyAFwxLAHDYPv6VoQ21nPNEEiQGEFRtk\\r\\nGf0NaWTOeW7Of8QwGG4MRZnEbYXPJwRnOR0zWNXW+KrqBLUWi5EjbWCgcAA9c\\/gRXKYqZaGlK\\/LqMH0F\\r\\nFLtHvRSNiYD2pSDTgpp6p0ywUHoTULXYxcktzrdCf7Xo8LP\\/AKyEmMNjJ46dfbFWJ5TDGNwB9lFUvDV9\\r\\nYrbfYGbyrjcWG88S57g+vtV26ZIvMlumKwwjLZ6V0WfU54yTvYwtbubea2WNWbzg4bYQeBgj8OtYeKhj\\r\\nu4y2HQxqxOD1xzxmrWAQCCGB6EGsaikndmsJxeiYzBo280\\/Z7UbayuaXGY5oIp+2lx9KLjIsVDeD\\/Rj\\/\\r\\nALy\\/zq1t96r3y4tT\\/vL\\/ADq4P3kRP4WSleTSFKkkKoCW4GaqNcMxIjXj1pxjKT0FKrGC1Nrw3vGrKkYz\\r\\n5kTAr6455\\/HH510UdwPtRgWCbzF5+YYUf4Vwun39xpmoR3qASMmQUJwGU9Rnt\\/8AWrpbrxhb8\\/ZdOmaQ\\r\\ngAGZwFH5ZJrpVKVlY5ZYhN6kXiu2eO\\/ikZlIljAAB5yM549OawSOOlPuLqe+umuLqTfM4OSOAo7ADsKh\\r\\nhl\\/cRsTuJHPv7mlKi3sVTxNtGP20VJhThgSQaK52mnZnUqsWrpkyeUrr5pABOAPU1AGaXUCWJISHGPfP\\r\\nP8qL7BiKnsMg46H3qrbzupbj5mPTPTpXVSglG551SpzSsXJ4\\/MBUgYIxyKpySyGBYJriV1D7kRpCVH4V\\r\\nbSeNJ4xchni3DeqnBI+td7F4b0mKIRjT45VbktJlzk455+n6VtYzv2PNwFZWBHBGKVJDGVC54\\/nXQeMN\\r\\nNttLNkba1jgWVWDmM8bhg4\\/nzXLSSbXVj6fyNKUdNRp21RtIRJGrjuM0u3FQ2DbodvcEkfQmrW2vLqLl\\r\\nk0ejCXNFMj2\\/jQV9qkxSYNRcsZiq2oI32N2CkhWXJxwOe9XMcVt6hoPn6dFaW0wgRpNzvKDlz6+\\/0rai\\r\\nryv2Jm9LHJai+ZRGCBjnr71ErdAxAY9B611t1Y2cunbbaOQ3FvKZI3UqGlZMbiWwfcfhV231iwvLSM3U\\r\\nlt5Uq52TuZG+hGMA12xXJGxxzjzybOQtNOvb5j9ktZJhnBIHyg+5PFX38JayqK\\/2eLJIBUTgkDA9q7ex\\r\\nitrSHFpGsUbndhRgc+g7VNIyfZJAoJZUbb3I46CtFJMylBo8sdWhmYMuCnylc9wef5VUT7+1chc5NS7h\\r\\nsUZO5RtIPUH3pkBDOxxxmqM9TQtn+WilhHfHaik43KTG3Z4IyPyrNVjGCsZ+dmwv6V3cXhSG8sYpJLud\\r\\nJJIwxChdoJGcYx\\/Wkg8DafA4knvLiQr\\/ALqj+VQpKw3FtnFFfvbiSMgZJ6\\/jXp2n3d9cQRBTFsKD96EP\\r\\noOxPU\\/8A68VVtbbRtMVntbePKDLTSHJH\\/Aj\\/AEqHTvE66rq72VugMMcbSGTnL4wMAfjT5n0HyW3L+s6b\\r\\nbaxaJBdzN+7bcrxkAhun0rz3VNCv7e7lgigknWI43xLu6jjIHTjtXqfkpPGVYsBkghTikgsYIN\\/lhgXb\\r\\ncxLkknp\\/ShczQ7xtY8vtEmhkj8yGRBuCnehUcnHcVtmwfJ\\/fQ8e7f\\/E12txZW91C0U6b42xlST2OR\\/Ko\\r\\nBo1gM\\/uW55\\/1jf41nOipu7LhV5FZHIGzI6zwj\\/vr\\/Ck+yr3uYf8Ax7\\/CutbQdMb71tn\\/ALaN\\/jSf8I\\/p\\r\\nX\\/PoP++2\\/wAan6rAr6wzkWt0II+1Rc\\/7Lf4Vd1eeCSKBbdZDdShYoiZNoyfY10P\\/AAj2lf8APmP++2\\/x\\r\\noPh\\/SjKspsozIuNrZORjp3qo0FHYPb3OZt7ae3SzjuItsiRSAgnccl\\/UA+3Q1yNjKLR4ZZYY5VD7tkv3\\r\\nWwO\\/+e1evPp9nI257aJm6bioz1z1+tY+s6Hplnot9PbWMMcqwOFcLyOO1bJWMZSTOPHi+9w3mosrlyd2\\r\\n9lCj02g9P\\/1e9a3hzxAbl2ikZRcdQueHHt7j864Y8Z4I4oRzG6urFWU5BHBB7HNJxTFGbR6he6Vpmtgm\\r\\neLy5zwZI\\/lb8fX8azIvBUUTHdfSFP4QsYB\\/HNZ+k+KEnRY75hHOvAk6K\\/v7H9K6yyvlnQBmDZ6GsnzR0\\r\\nN0oy1RzOtaN\\/Y1tHNFO06u+zYy4I4Jzx9KKveJblXuordSGES5b6n\\/62PzorKVdp2LjQTVyWz8UWEWlq\\r\\njSgyxfJt6EgdDzWTdeLIZGO7zHI\\/hVajGmWWP+PWL8qwlAIURrhpMAHHJA71pRcZrToZzcoEuo6heakA\\r\\nGHk245CZ6\\/X1qPTLq40q+W5t2QybSpDAkEEc55\\/zilk5k2r91eKhLDzWz2rpsczbbuemeD76fUNG865I\\r\\nMiysmQMZAAwa3a5j4ftu0ByP+fh\\/5CulkLLG7INzhSVHqe1Fh3uOoqn9qQQxyhndmHIxwOmSR2xQ13KD\\r\\nKoiBZOV9JBnt707MVy5RWdNdy7wRGf3bfMinnO1jg+vY03WXLaJO3mhQ20b0zwpYf0qlG7S7icrJs08U\\r\\nVwumgC+YiQyeVtZH567hzj8aSL949oGhE\\/2v5pJCDkksQwBHC4\\/+vXQ8LZ2uYxxCavY7us\\/xCcaBfn0h\\r\\nb+VP0bnSrb94ZMJgOecj1rl\\/GfidUE2k2gy5+SeQjgA\\/wj3rlas2jdao48qrjLAGkSKPk4Gc1WMj92I+\\r\\nlIJnU8OfxPWo5inBokmtQTmM4OOh71b0q6vbFmWCbaxHyqQGAP0PT8KhSTzVyo5ocSKA5VfTOTmqsmRd\\r\\npl99XjPzThzK3zOeOSeveirNmkgg\\/fIpYsTkYORxRXmzlTjJqx6EVUcU7mhkKCzdAK59QI9zYxtG1fYU\\r\\nUVtgtmY4nZEa8Ak9aqFv3rfSiiu1nMeifDv\\/AJF+T\\/r4f+QrqqKKQwzQenNFFMCOKFIgNuThdoJ5OPSk\\r\\nubeK6t3gnXdG4wwziiii\\/UTKMOg6dbzJLFE4dSCP3rEdeOM8805tDsGMvySgSsS6rM6gk9eAcUUVftZt\\r\\n3uyVGNthuq3Eei6DK8H7sRR7YuMgHtXkc8rzTNLM26RyWY+p70UVnLY0iEsUipG7rhZBlDkc1HgYoorM\\r\\n0HwyBXGeRjmrcUhMg2ghezd\\/\\/rUUVcTKW5s2jZtY\\/QDaOKKKK8ip8bPRj8KP\\/9k="

	tests := []struct {
		name    string
		args    args
		mock    func(a args, mock *MockuserRepo)
		wantErr bool
	}{
		{
			name: "Success UploadUserPic",
			args: args{
				id:       1,
				username: "testName",
				picData:  picData,
			},
			mock: func(a args, mock *MockuserRepo) {
				db, dbMock, _ := sqlmock.New()
				dbMock.ExpectBegin()
				tx, _ := db.Begin()
				dbMock.ExpectCommit()

				mock.EXPECT().CreateTx().Return(tx, nil)
				mock.EXPECT().UpdateUserPic("Pict_testUser", a.id, tx).Return(nil)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := NewMockuserRepo(ctrl)
			u := &UseCase{
				ur: mockRepo,
			}

			tt.mock(tt.args, mockRepo)
			if err := u.UploadUserPic(tt.args.id, tt.args.username, tt.args.picData); (err != nil) != tt.wantErr {
				t.Errorf("UploadUserPic() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCase_checkPasswordHash(t *testing.T) {
	type fields struct {
		ur userRepo
	}
	type args struct {
		password string
		hash     string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UseCase{
				ur: tt.fields.ur,
			}
			if got := u.checkPasswordHash(tt.args.password, tt.args.hash); got != tt.want {
				t.Errorf("checkPasswordHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCase_generateJWT(t *testing.T) {
	type fields struct {
		ur userRepo
	}
	type args struct {
		username string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UseCase{
				ur: tt.fields.ur,
			}
			got, err := u.generateJWT(tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("generateJWT() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCase_hashPassword(t *testing.T) {
	type fields struct {
		ur userRepo
	}
	type args struct {
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UseCase{
				ur: tt.fields.ur,
			}
			got, err := u.hashPassword(tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("hashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("hashPassword() got = %v, want %v", got, tt.want)
			}
		})
	}
}
