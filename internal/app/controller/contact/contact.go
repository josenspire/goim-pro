package contact

import (
	"context"
	protos "goim-pro/api/protos/salty"
	contactsrv "goim-pro/internal/app/services/contact"
	"goim-pro/internal/app/services/converters"
	errmsg "goim-pro/pkg/errors"
	"goim-pro/pkg/logs"
	"goim-pro/pkg/utils"
	"strings"
)

var (
	logger         = logs.GetLogger("INFO")
	contactService contactsrv.IContactService
)

type contactServer struct {
}

func New() protos.ContactServiceServer {
	contactService = contactsrv.New()
	return &contactServer{}
}

func (s *contactServer) RequestContact(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	resp, _ = utils.NewGRPCResp(protos.StatusCode_STATUS_OK, nil, "")

	var err error
	var reqContactReq protos.RequestContactReq
	if err = utils.UnmarshalGRPCReq(req, &reqContactReq); err != nil {
		resp.Code = protos.StatusCode_STATUS_BAD_REQUEST
		resp.Message = err.Error()
		logger.Errorf(`unmarshal error: %v`, err)
		return
	}

	userId := req.GetToken()

	if err = requestContactParameterCalibration(userId, &reqContactReq); err != nil {
		resp.Code = protos.StatusCode_STATUS_BAD_REQUEST
		resp.Message = err.Error()
		return
	}

	contactId := reqContactReq.GetUserId()
	requestReason := reqContactReq.GetReason()

	tErr := contactService.RequestContact(userId, contactId, requestReason)
	if tErr != nil {
		resp.Code = tErr.Code
		resp.Message = tErr.Detail
		return
	}
	resp.Message = "the request to add contact succeeded"
	return
}

func (s *contactServer) RefusedContact(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	resp, _ = utils.NewGRPCResp(protos.StatusCode_STATUS_OK, nil, "")

	var err error
	var refContactReq protos.RefusedContactReq
	if err = utils.UnmarshalGRPCReq(req, &refContactReq); err != nil {
		resp.Code = protos.StatusCode_STATUS_BAD_REQUEST
		resp.Message = err.Error()
		logger.Errorf(`unmarshal error: %v`, err)
		return
	}

	userId := req.GetToken()

	if err = refusedContactParameterCalibration(userId, &refContactReq); err != nil {
		resp.Code = protos.StatusCode_STATUS_BAD_REQUEST
		resp.Message = err.Error()
		return
	}

	contactId := refContactReq.GetUserId()
	refusedReason := refContactReq.GetReason()

	tErr := contactService.RefusedContact(userId, contactId, refusedReason)
	if tErr != nil {
		resp.Code = tErr.Code
		resp.Message = tErr.Detail
		return
	}
	return
}

func (s *contactServer) AcceptContact(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	resp, _ = utils.NewGRPCResp(protos.StatusCode_STATUS_OK, nil, "")

	var err error
	var acpContactReq protos.AcceptContactReq
	if err = utils.UnmarshalGRPCReq(req, &acpContactReq); err != nil {
		resp.Code = protos.StatusCode_STATUS_BAD_REQUEST
		resp.Message = err.Error()
		logger.Errorf(`unmarshal error: %v`, err)
		return
	}

	userId := req.GetToken()
	contactId := strings.Trim(acpContactReq.GetUserId(), "")

	if contactId == "" {
		resp.Code = protos.StatusCode_STATUS_BAD_REQUEST
		resp.Message = errmsg.ErrInvalidParameters.Error()
		return
	}
	if strings.EqualFold(userId, contactId) {
		resp.Code = protos.StatusCode_STATUS_BAD_REQUEST
		resp.Message = errmsg.ErrIllegalOperation.Error()
		return
	}

	tErr := contactService.AcceptContact(userId, contactId)
	if tErr != nil {
		resp.Code = tErr.Code
		resp.Message = tErr.Detail
		return
	}

	resp.Message = "successfully accepted"
	return
}

func (s *contactServer) DeleteContact(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	resp, _ = utils.NewGRPCResp(protos.StatusCode_STATUS_OK, nil, "")

	var err error
	var delContactReq protos.DeleteContactReq
	if err = utils.UnmarshalGRPCReq(req, &delContactReq); err != nil {
		resp.Code = protos.StatusCode_STATUS_BAD_REQUEST
		resp.Message = err.Error()
		logger.Errorf(`unmarshal error: %v`, err)
		return
	}

	userId := req.GetToken()
	contactId := strings.Trim(delContactReq.GetUserId(), "")

	if contactId == "" {
		resp.Code = protos.StatusCode_STATUS_BAD_REQUEST
		resp.Message = errmsg.ErrInvalidParameters.Error()
		return
	}
	if strings.EqualFold(userId, contactId) {
		resp.Code = protos.StatusCode_STATUS_BAD_REQUEST
		resp.Message = errmsg.ErrIllegalOperation.Error()
		return
	}

	tErr := contactService.DeleteContact(userId, contactId)
	if tErr != nil {
		resp.Code = tErr.Code
		resp.Message = tErr.Detail
		return
	}

	resp.Message = "contact deleted successfully"
	return
}

func (s *contactServer) UpdateRemarkInfo(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	resp, _ = utils.NewGRPCResp(protos.StatusCode_STATUS_OK, nil, "")

	var err error
	var updateRemarkInfoReq protos.UpdateRemarkInfoReq
	if err = utils.UnmarshalGRPCReq(req, &updateRemarkInfoReq); err != nil {
		resp.Code = protos.StatusCode_STATUS_BAD_REQUEST
		resp.Message = err.Error()
		logger.Errorf(`unmarshal error: %v`, err)
		return
	}

	userId := req.GetToken()
	contactId := strings.Trim(updateRemarkInfoReq.GetUserId(), "")
	contactRemark := updateRemarkInfoReq.GetRemarkInfo()

	if contactId == "" {
		resp.Code = protos.StatusCode_STATUS_BAD_REQUEST
		resp.Message = errmsg.ErrInvalidParameters.Error()
		return
	}

	tErr := contactService.UpdateRemarkInfo(userId, contactId, contactRemark)
	if tErr != nil {
		resp.Code = tErr.Code
		resp.Message = tErr.Detail
		return
	}

	resp.Message = "contact remark profile updated successfully"
	return
}

func (s *contactServer) GetContactList(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	resp, _ = utils.NewGRPCResp(protos.StatusCode_STATUS_OK, nil, "")
	userId := req.GetToken()

	contacts, tErr := contactService.GetContactList(userId)
	if tErr != nil {
		resp.Code = tErr.Code
		resp.Message = tErr.Detail
		return
	}
	if contacts == nil || len(contacts) == 0 {
		return
	}

	getContactListResp := &protos.GetContactListResp{
		ContactList: converters.ConvertEntity2ProtoForContacts(contacts),
	}
	anyData, err := utils.MarshalMessageToAny(getContactListResp)
	if err != nil {
		logger.Errorf("[get contacts] response marshal message error: %s", err.Error())
		return
	}
	resp.Data = anyData
	return
}

func (s *contactServer) GetContactOperationList(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	resp, _ = utils.NewGRPCResp(protos.StatusCode_STATUS_OK, nil, "")

	var err error
	var optsMessageReq protos.GetContactOperationListReq
	if err = utils.UnmarshalGRPCReq(req, &optsMessageReq); err != nil {
		resp.Code = protos.StatusCode_STATUS_BAD_REQUEST
		resp.Message = err.Error()
		logger.Errorf(`unmarshal error: %v`, err)
		return
	}

	userId := req.GetToken()

	startDateTime := optsMessageReq.StartDateTime
	endDateTime := optsMessageReq.EndDateTime

	notifications, tErr := contactService.GetContactOperationList(userId, startDateTime, endDateTime)
	if tErr != nil {
		resp.Code = tErr.Code
		resp.Message = tErr.Detail
		return
	}
	if notifications == nil || len(notifications) == 0 {
		return
	}

	getContactNotificationMsg := &protos.GetContactOperationListResp{
		MessageList:          converters.ConvertEntity2ProtoForNotificationMsg(notifications),
	}
	anyData, err := utils.MarshalMessageToAny(getContactNotificationMsg)
	if err != nil {
		logger.Errorf("[get contacts operations message] response marshal message error: %s", err.Error())
		return
	}
	resp.Data = anyData
	return
}

func requestContactParameterCalibration(userId string, req *protos.RequestContactReq) (err error) {
	csErr := errmsg.ErrInvalidParameters

	contactId := strings.Trim(req.UserId, "")
	requestReason := strings.Trim(req.Reason, "")
	if utils.IsEmptyStrings(contactId) {
		err = csErr
	} else if strings.EqualFold(userId, contactId) {
		err = errmsg.ErrIllegalOperation
	} else {
		req.UserId = contactId
		req.Reason = requestReason
	}
	return
}

func refusedContactParameterCalibration(userId string, req *protos.RefusedContactReq) (err error) {
	csErr := errmsg.ErrInvalidParameters

	contactId := strings.Trim(req.UserId, "")
	requestReason := strings.Trim(req.Reason, "")
	if utils.IsEmptyStrings(contactId) {
		err = csErr
	} else if strings.EqualFold(userId, contactId) {
		err = errmsg.ErrIllegalOperation
	} else {
		req.UserId = contactId
		req.Reason = requestReason
	}
	return
}
