package grpc

import (
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc"
	"reflect"
	"testing"
)

func TestGRPCServer_ConnectGRPCServer(t *testing.T) {
	type fields struct {
		grpcServer *grpc.Server
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &GRPCServer{
				grpcServer: tt.fields.grpcServer,
			}
		})
	}
}

func TestGRPCServer_ForceStopGRPCServer(t *testing.T) {
	type fields struct {
		grpcServer *grpc.Server
	}
	var tests []struct {
		name   string
		fields fields
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &GRPCServer{
				grpcServer: tt.fields.grpcServer,
			}
		})
	}
}

func TestGRPCServer_GracefulStopGRPCServer(t *testing.T) {
	type fields struct {
		grpcServer *grpc.Server
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &GRPCServer{
				grpcServer: tt.fields.grpcServer,
			}
		})
	}
}

func TestGRPCServer_InitServer(t *testing.T) {
	type fields struct {
		grpcServer *grpc.Server
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &GRPCServer{
				grpcServer: tt.fields.grpcServer,
			}
		})
	}
}

func TestNewServer(t *testing.T) {
	tests := []struct {
		name string
		want *GRPCServer
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewServer(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handleServiceRegister(t *testing.T) {
	type args struct {
		srv *grpc.Server
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func Test_initialMysqlTables(t *testing.T) {
	type args struct {
		db *gorm.DB
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := initialMysqlTables(tt.args.db); (err != nil) != tt.wantErr {
				t.Errorf("initialMysqlTables() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
