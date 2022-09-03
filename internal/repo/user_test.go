package repo

import (
	"database/sql"
	"entry_task/internal/component"
	"entry_task/internal/model/user"
	"github.com/DATA-DOG/go-sqlmock"
	"reflect"
	"testing"
)

type dummy struct {
	mockDB sqlmock.Sqlmock
	fields
	transaction *sql.Tx
}

type fields struct {
	db *sql.DB
}

//func newFromDB(masterDB *sql.DB) *component.DB{
//
//}

func setDBDummy() (dummy, error) {
	db, mockDB, err := sqlmock.New()
	if err != nil {
		return dummy{}, err
	}

	//dataStore := component.InitDatabase()
	//if err != nil {
	//	return dummy{}, err
	//}
	mockDB.ExpectBegin()

	tx, err := db.Begin()
	if err != nil {
		return dummy{}, err
	}

	return dummy{
		mockDB: mockDB,
		fields: fields{
			db: db,
		},
		transaction: tx,
	}, nil
}

func TestRepo_CreateTx(t *testing.T) {

	dummyEnv, err := setDBDummy()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}

	defer dummyEnv.fields.db.Close()

	tests := []struct {
		name    string
		mock    func()
		wantErr bool
	}{
		{
			name: "Success",
			mock: func() {
				dummyEnv.mockDB.ExpectBegin()
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repo{
				db: &component.DB{
					Master: dummyEnv.db,
				},
			}

			tt.mock()
			_, err := r.CreateTx()
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestRepo_GetUserByID(t *testing.T) {
	type args struct {
		userID int
	}

	dummyEnv, err := setDBDummy()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}

	defer dummyEnv.fields.db.Close()

	returnData := user.User{
		Id:       1,
		UserName: "testName",
		NickName: "test1",
		Password: "testPass",
	}

	tests := []struct {
		name    string
		args    args
		mock    func()
		want    user.User
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				userID: 1,
			},
			mock: func() {
				dummyEnv.mockDB.ExpectQuery("SELECT id, username, nickname, password FROM").
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "username", "nickname", "password",
					}).AddRow(
						returnData.Id,
						returnData.UserName,
						returnData.NickName,
						returnData.Password),
					)
			},
			want:    returnData,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repo{
				db: &component.DB{
					Master: dummyEnv.db,
				},
			}

			tt.mock()
			got, err := r.GetUserByID(tt.args.userID)
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

func TestRepo_GetUserByName(t *testing.T) {
	type args struct {
		username string
	}

	dummyEnv, err := setDBDummy()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}

	defer dummyEnv.fields.db.Close()

	returnData := user.User{
		Id:       1,
		UserName: "testName",
		NickName: "test1",
		Password: "testPass",
	}

	tests := []struct {
		name    string
		args    args
		mock    func()
		want    user.User
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				username: "testName",
			},
			mock: func() {
				dummyEnv.mockDB.ExpectQuery("SELECT id, username, nickname, password FROM").
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "username", "nickname", "password",
					}).AddRow(
						returnData.Id,
						returnData.UserName,
						returnData.NickName,
						returnData.Password),
					)
			},
			want:    returnData,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repo{
				db: &component.DB{
					Master: dummyEnv.db,
				},
			}

			tt.mock()
			got, err := r.GetUserByName(tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserByName() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepo_UpdateUserPic(t *testing.T) {
	type args struct {
		picName string
		userID  int
		tx      *sql.Tx
	}

	dummyEnv, err := setDBDummy()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}

	defer dummyEnv.fields.db.Close()

	tests := []struct {
		name    string
		args    args
		mock    func()
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				picName: "user1_png",
				userID:  1,
				tx:      nil,
			},
			mock: func() {
				dummyEnv.mockDB.ExpectExec("UPDATE user").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repo{
				db: &component.DB{
					Master: dummyEnv.db,
				},
			}

			tt.mock()
			if err := r.UpdateUserPic(tt.args.picName, tt.args.userID, tt.args.tx); (err != nil) != tt.wantErr {
				t.Errorf("UpdateUserPic() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepo_UpsertUser(t *testing.T) {
	type args struct {
		user user.User
		tx   *sql.Tx
	}

	dummyEnv, err := setDBDummy()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}

	defer dummyEnv.fields.db.Close()

	insertData := user.User{
		Id:       0,
		UserName: "testName",
		NickName: "test1",
		Password: "testPass",
	}

	updateData := insertData
	updateData.Id = 1

	tests := []struct {
		name    string
		args    args
		mock    func()
		wantErr bool
	}{
		{
			name: "Success Insert",
			args: args{
				user: insertData,
				tx:   nil,
			},
			mock: func() {
				dummyEnv.mockDB.ExpectExec("INSERT INTO user").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "Success Update",
			args: args{
				user: updateData,
				tx:   nil,
			},
			mock: func() {
				dummyEnv.mockDB.ExpectExec("UPDATE user").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repo{
				db: &component.DB{
					Master: dummyEnv.db,
				},
			}

			tt.mock()
			if err := r.UpsertUser(tt.args.user, tt.args.tx); (err != nil) != tt.wantErr {
				t.Errorf("UpsertUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
