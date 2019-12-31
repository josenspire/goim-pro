package userrepo

import (
	"github.com/jinzhu/gorm"
	"goim-pro/internal/app/repo"
	"goim-pro/pkg/db"
	"reflect"
	"testing"
)

func TestNewUserModel(t *testing.T) {
	tests := []struct {
		name string
		want *User
	}{
		{
			name: "test_for_NewUserModel_method",
			want: &User{},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserModel(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserModel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewUserRepo(t *testing.T) {
	_ = db.GetMysqlConnection().InitConnectionPool()
	mysqlDB := db.GetMysqlConnection().GetMysqlDBInstance()

	type args struct {
		db *gorm.DB
	}
	tests := []struct {
		name string
		args args
		want IUserRepo
	}{
		// TODO: Add test cases.
		{
			name: "test_for_NewUserRepo_method",
			args: args{
				mysqlDB,
			},
			want: &User{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserRepo(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserRepo() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TODO: should verify the result, have not got all pass
func TestUser_IsTelephoneRegistered(t *testing.T) {
	type fields struct {
		UserID      uint64
		Password    string
		Role        string
		Status      string
		UserProfile UserProfile
		BaseModel   repo.BaseModel
	}
	type args struct {
		telephone string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test_for_IsTelephoneExist_returns_True",
			fields: fields{
				UserID:   0,
				Password: "uiowpqejfmdlvm",
				Role:     "0",
				Status:   "1",
				UserProfile: UserProfile{
					Telephone: "13631210000",
				},
				BaseModel: repo.BaseModel{},
			},
			args: args{
				telephone: "13631210000",
			},
			want:    true,
			wantErr: true,
		},
		{
			name: "test_for_IsTelephoneExist_returns_False",
			fields: fields{
				UserID:   1,
				Password: "zxcvdfgreqgrewqg",
				Role:     "0",
				Status:   "1",
				UserProfile: UserProfile{
					Telephone: "13631210001",
				},
				BaseModel: repo.BaseModel{},
			},
			args: args{
				telephone: "13631210001",
			},
			want:    false,
			wantErr: false,
		},
	}
	db.GetMysqlConnection().InitConnectionPool()
	mysqlDB := db.GetMysqlConnection().GetMysqlDBInstance()
	NewUserRepo(mysqlDB)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				UserID:      0,
				Password:    "",
				Role:        "",
				Status:      "",
				UserProfile: UserProfile{},
				BaseModel:   repo.BaseModel{},
			}
			got, _ := u.IsTelephoneRegistered(tt.args.telephone)
			if got != tt.want {
				t.Errorf("IsTelephoneRegistered() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_Register(t *testing.T) {
	type fields struct {
		UserID      uint64
		Password    string
		Role        string
		Status      string
		UserProfile UserProfile
		BaseModel   repo.BaseModel
	}
	type args struct {
		newUser *User
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantUser *User
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				UserID:      tt.fields.UserID,
				Password:    tt.fields.Password,
				Role:        tt.fields.Role,
				Status:      tt.fields.Status,
				UserProfile: tt.fields.UserProfile,
				BaseModel:   tt.fields.BaseModel,
			}
			gotUser, err := u.Register(tt.args.newUser)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotUser, tt.wantUser) {
				t.Errorf("Register() gotUser = %v, want %v", gotUser, tt.wantUser)
			}
		})
	}
}
