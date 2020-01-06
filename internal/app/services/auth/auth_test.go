package authsrv

import (
	"context"
	"github.com/golang/protobuf/ptypes/any"
	"goim-pro/api/protos"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want protos.SMSServiceServer
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

func Test_smsServer_ObtainSMSCode(t *testing.T) {
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
			name:    "testing_for_obtain_sms_code_register",
			args:    args{
				ctx: context.Background(),
				req: &protos.BaseClientRequest{
					Code: 0,
					Data: &any.Any{
						TypeUrl: "",
						Value: []byte("register"),
					},
					Message: "",
				},
			},
			wantRes: &protos.BaseServerResponse{
				Code:                 200,
				Data:                 &any.Any{
					TypeUrl:              "",
					Value:                []byte("123456"),
				},
				Message:              "sending sms code success",
			},
			wantErr: false,
		},
		{
			name:    "testing_for_obtain_sms_code_by_error_type",
			args:    args{
				ctx: context.Background(),
				req: &protos.BaseClientRequest{
					Code: 0,
					Data: &any.Any{
						TypeUrl: "",
						Value: []byte("CHECK"),
					},
					Message: "",
				},
			},
			wantRes: &protos.BaseServerResponse{
				Code:                 400,
				Data:                 &any.Any{
					TypeUrl:              "",
					Value:                []byte("123456"),
				},
				Message:              "invalid request code type",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &smsServer{}
			gotRes, err := s.ObtainSMSCode(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ObtainSMSCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes.GetData(), tt.wantRes.GetData()) {
				t.Errorf("ObtainSMSCode() gotRes = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
