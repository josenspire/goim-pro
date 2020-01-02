package user

import (
	"goim-pro/internal/app/repos/base"
	"goim-pro/pkg/db"
	"reflect"
	"testing"

	"github.com/jinzhu/gorm"
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

func TestUser_IsTelephoneRegistered(t *testing.T) {
	type fields struct {
		UserID      uint64
		Password    string
		Role        string
		Status      string
		UserProfile UserProfile
		BaseModel   base.BaseModel
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
				BaseModel: base.BaseModel{},
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
				BaseModel: base.BaseModel{},
			},
			args: args{
				telephone: "13631210001",
			},
			want:    false,
			wantErr: false,
		},
	}
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
				BaseModel:   base.BaseModel{},
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
		BaseModel   base.BaseModel
	}
	type args struct {
		newUser *User
	}
	mysqlDB := db.GetMysqlConnection().GetMysqlDBInstance()
	NewUserRepo(mysqlDB)
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "testing_register_with_new_record",
			fields: fields{
				UserID:   2,
				Password: "1234567890",
				UserProfile: UserProfile{
					Telephone: "13631210010",
					Email:     "294001@qq.com",
					Username:  "TEST02",
					Nickname:  "TEST02",
					Signature: "Never Settle",
				},
			},
			args: args{
				&User{
					UserID:   2,
					Password: "1234567890",
					UserProfile: UserProfile{
						Telephone: "13631210010",
						Email:     "294001@qq.com",
						Username:  "TEST02",
						Nickname:  "TEST02",
						Signature: "Never Settle",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "testing_register_with_exist_record",
			fields: fields{
				UserID:   2,
				Password: "1234567890",
				UserProfile: UserProfile{
					Telephone: "13631210010",
					Email:     "294001@qq.com",
					Username:  "TEST02",
					Nickname:  "TEST02",
					Signature: "Never Settle",
				},
			},
			args: args{
				&User{
					UserID:   2,
					Password: "1234567890",
					UserProfile: UserProfile{
						Telephone: "13631210010",
						Email:     "294001@qq.com",
						Username:  "TEST02",
						Nickname:  "TEST02",
						Signature: "Never Settle",
					},
				},
			},
			wantErr: true,
		},
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
			if err := u.Register(tt.args.newUser); (err != nil) != tt.wantErr {
				t.Errorf("User.Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUser_RemoveUserByUserId(t *testing.T) {
	type fields struct {
		UserID      uint64
		Password    string
		Role        string
		Status      string
		UserProfile UserProfile
		BaseModel   base.BaseModel
	}
	type args struct {
		userID uint64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "testing_for_remove_user_by_userID",
			fields:  fields{
				UserID:      2,
			},
			args:    args{
				userID: 2,
			},
			wantErr: false,
		},
	}
	mysqlDB := db.GetMysqlConnection().GetMysqlDBInstance()
	NewUserRepo(mysqlDB)
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
			if err := u.RemoveUserByUserId(tt.args.userID); (err != nil) != tt.wantErr {
				t.Errorf("User.RemoveUserByUserId() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
