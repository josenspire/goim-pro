package contactsrv

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v7"
	protos "goim-pro/api/protos/salty"
	"goim-pro/internal/app/repos"
	"goim-pro/internal/app/repos/contact"
	"goim-pro/internal/app/repos/user"
	redsrv "goim-pro/pkg/db/redis"
	"goim-pro/pkg/logs"
	"goim-pro/pkg/utils"
	"net/http"
	"strings"
	"time"
)

var (
	logger      = logs.GetLogger("INFO")
	expiresTime = time.Hour * time.Duration(24*7) // 7 days
	myRedis     *redis.Client
)

type contactService struct {
	userRepo    user.IUserRepo
	contactRepo contact.IContactRepo
}

func New() protos.ContactServiceServer {
	myRedis = redsrv.NewRedisConnection().GetRedisClient()

	repoServer := repos.New()
	return &contactService{
		userRepo:    repoServer.UserRepo,
		contactRepo: repoServer.ContactRepo,
	}
}

// request add contact
func (cs *contactService) RequestContact(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	resp, _ = utils.NewGRPCResp(http.StatusOK, nil, "")

	var err error
	var reqContactReq protos.RequestContactReq
	if err = utils.UnmarshalGRPCReq(req, &reqContactReq); err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = err.Error()
		logger.Errorf(`unmarshal error: %v`, err)
		return
	}

	userId := req.GetToken()

	if err = requestContactParameterCalibration(userId, &reqContactReq); err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = err.Error()
		return
	}

	contactId := reqContactReq.GetUserId()
	requestReason := reqContactReq.GetReason()

	_, err = cs.userRepo.FindByUserId(contactId)
	if err != nil {
		if err == utils.ErrInvalidUserId {
			resp.Code = http.StatusBadRequest
			resp.Message = utils.ErrInvalidContact.Error()
			return
		}
		logger.Errorf("find contact by userId error: %s", err.Error())
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		return
	}

	isExists, err := cs.contactRepo.IsExistContact(userId, contactId)
	if err != nil {
		logger.Errorf("checking contact error: %s", err.Error())
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		return
	}
	if isExists {
		resp.Code = http.StatusBadRequest
		resp.Message = utils.ErrContactAlreadyExists.Error()
		return
	}
	// TODO: cache in redis, should replace to Push notification server
	ctKey := fmt.Sprintf("CT-REQ-%s-%s", userId, contactId)
	cacheContent := make(map[string]interface{})
	cacheContent["contactId"] = contactId
	cacheContent["reason"] = requestReason

	for field, value := range cacheContent {
		err = myRedis.HSet(ctKey, field, value).Err()
	}
	if err != nil {
		logger.Errorf("redis cache error: %s", err.Error())
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		return
	}
	resp.Message = "the request to add contact succeeded"
	return
}

// refused request contact
func (cs *contactService) RefusedContact(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	resp, _ = utils.NewGRPCResp(http.StatusOK, nil, "")

	var err error
	var refContactReq protos.RefusedContactReq
	if err = utils.UnmarshalGRPCReq(req, &refContactReq); err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = err.Error()
		logger.Errorf(`unmarshal error: %v`, err)
		return
	}

	userId := req.GetToken()

	if err = refusedContactParameterCalibration(userId, &refContactReq); err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = err.Error()
		return
	}

	contactId := refContactReq.GetUserId()
	refusedReason := refContactReq.GetReason()

	_, err = cs.userRepo.FindByUserId(contactId)
	if err != nil {
		if err == utils.ErrInvalidUserId {
			resp.Code = http.StatusBadRequest
			resp.Message = utils.ErrInvalidContact.Error()
			return
		}
		logger.Errorf("find contact by userId error: %s", err.Error())
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		return
	}

	isExists, err := cs.contactRepo.IsExistContact(userId, contactId)
	if err != nil {
		logger.Errorf("checking contact error: %s", err.Error())
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		return
	}
	if isExists {
		resp.Code = http.StatusBadRequest
		resp.Message = utils.ErrContactAlreadyExists.Error()
		return
	}

	// TODO: cache in redis, should replace to Push notification server
	ctKey := fmt.Sprintf("CT-REF-%s-%s", userId, contactId)
	cacheContent := make(map[string]interface{})
	cacheContent["contactId"] = contactId
	cacheContent["reason"] = refusedReason

	for field, value := range cacheContent {
		err = myRedis.HSet(ctKey, field, value).Err()
	}
	if err != nil {
		logger.Errorf("redis cache error: %s", err.Error())
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		return
	}
	return
}

func (cs *contactService) AcceptContact(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	resp, _ = utils.NewGRPCResp(http.StatusOK, nil, "")

	var err error
	var acpContactReq protos.AcceptContactReq
	if err = utils.UnmarshalGRPCReq(req, &acpContactReq); err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = err.Error()
		logger.Errorf(`unmarshal error: %v`, err)
		return
	}

	userId := req.GetToken()
	contactId := strings.Trim(acpContactReq.GetUserId(), "")

	if contactId == "" {
		resp.Code = http.StatusBadRequest
		resp.Message = utils.ErrInvalidParameters.Error()
		return
	}
	if strings.EqualFold(userId, contactId) {
		resp.Code = http.StatusBadRequest
		resp.Message = utils.ErrIllegalOperation.Error()
		return
	}

	_, err = cs.userRepo.FindByUserId(contactId)
	if err != nil {
		if err == utils.ErrInvalidUserId {
			resp.Code = http.StatusBadRequest
			resp.Message = utils.ErrInvalidContact.Error()
			return
		}
		logger.Errorf("find contact by userId error: %s", err.Error())
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		return
	}

	isExists, err := cs.contactRepo.IsExistContact(userId, contactId)
	if err != nil {
		logger.Errorf("checking contact error: %s", err.Error())
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		return
	}
	if isExists {
		resp.Code = http.StatusBadRequest
		resp.Message = utils.ErrContactAlreadyExists.Error()
		return
	}

	// TODO: cache in redis, should replace to Push notification server
	ctKey := fmt.Sprintf("CT-ACP-%s-%s", userId, contactId)
	err = myRedis.HSet(ctKey, "contactId", contactId).Err()
	if err != nil {
		logger.Errorf("redis cache error: %s", err.Error())
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		return
	}

	if err = handleAcceptContact(cs, userId, contactId); err != nil {
		logger.Errorf("insert contact error: %s", err.Error())
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		return
	}
	resp.Message = "successfully accepted"
	return
}

func (cs *contactService) DeleteContact(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	panic("implement me")
}

func (cs *contactService) UpdateRemarkInfo(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	panic("implement me")
}

func requestContactParameterCalibration(userId string, req *protos.RequestContactReq) (err error) {
	csErr := utils.ErrInvalidParameters

	contactId := strings.Trim(req.UserId, "")
	requestReason := strings.Trim(req.Reason, "")
	if utils.IsEmptyStrings(contactId) {
		err = csErr
	} else if strings.EqualFold(userId, contactId) {
		err = utils.ErrIllegalOperation
	} else {
		req.UserId = contactId
		req.Reason = requestReason
	}
	return
}

func refusedContactParameterCalibration(userId string, req *protos.RefusedContactReq) (err error) {
	csErr := utils.ErrInvalidParameters

	contactId := strings.Trim(req.UserId, "")
	requestReason := strings.Trim(req.Reason, "")
	if utils.IsEmptyStrings(contactId) {
		err = csErr
	} else if strings.EqualFold(userId, contactId) {
		err = utils.ErrIllegalOperation
	} else {
		req.UserId = contactId
		req.Reason = requestReason
	}
	return
}


func handleAcceptContact(cs *contactService, userId string, contactId string) (err error) {
	newContact1 := &contact.Contact{
		UserId:    userId,
		ContactId: contactId,
	}
	newContact2 := &contact.Contact{
		UserId:    contactId,
		ContactId: userId,
	}

	err = cs.contactRepo.InsertContacts(newContact1, newContact2)
	return err
}
