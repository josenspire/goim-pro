package contact

import (
	"context"
	protos "goim-pro/api/protos/salty"
	contactsrv "goim-pro/internal/app/services/contact"
	"goim-pro/pkg/logs"
)

var (
	logger         = logs.GetLogger("INFO")
	contactService *contactsrv.ContactService
)

type contactServer struct {
}

func New() protos.ContactServiceServer {
	contactService = contactsrv.New()
	return &contactServer{}
}

func (s *contactServer) RequestContact(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	panic("implement me")
}

func (s *contactServer) RefusedContact(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	panic("implement me")
}

func (s *contactServer) AcceptContact(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	panic("implement me")
}

func (s *contactServer) DeleteContact(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	panic("implement me")
}

func (s *contactServer) UpdateRemarkInfo(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	panic("implement me")
}

func (s *contactServer) GetContacts(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	panic("implement me")
}
