package usersrv

import (
	"context"
	any "github.com/golang/protobuf/ptypes/any"
	"goim-pro/api/protos"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want protos.UserServiceServer
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userServer_Login(t *testing.T) {
	type args struct {
		ctx context.Context
		req *protos.BaseClientRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *protos.BaseServerResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := &userServer{}
			got, err := us.Login(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Login() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userServer_Register(t *testing.T) {
	type args struct {
		ctx context.Context
		req *protos.BaseClientRequest
	}
	tests := []struct {
		name    string
		args    args
		wantRes *protos.BaseServerResponse
		wantErr bool
	}{
		{
			name: "testing_for_grpc_register",
			args: args{
				ctx: context.Background(),
				req: &protos.BaseClientRequest{
					Code: 0,
					Data: &any.Any{
						TypeUrl: "",
						Value: []byte{123, 34, 97, 103, 101, 34, 58, 50, 52, 44, 34, 112, 97, 115, 115, 119, 111, 114, 100, 34, 58, 34, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 34, 44, 34, 117, 115, 101, 114, 110, 97, 109,
							101, 34, 58, 34, 74, 65, 77, 69, 83, 34, 125},
					},
					Message: "",
				},
			},
			wantRes: &protos.BaseServerResponse{
				Code: 200,
				Data: &any.Any{
					Value: []byte("user regist successful.."),
				},
				Message: "user regist successful..",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := &userServer{}
			gotRes, err := us.Register(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes.GetData(), tt.wantRes.GetData()) {
				t.Errorf("Register() gotRes = %v, want %v", gotRes.GetData(), tt.wantRes.GetData())
			}
		})
	}
}
