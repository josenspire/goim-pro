package grpc

import (
	"google.golang.org/grpc"
	"testing"
)

func TestGRPCServer_InitServer(t *testing.T) {
	type fields struct {
		grpcServer *grpc.Server
	}
	var tests []struct {
		name   string
		fields fields
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = &GRPCServer{
				grpcServer: tt.fields.grpcServer,
			}
		})
	}
}
