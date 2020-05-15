package contactsrv

import (
	"context"
	"fmt"
	"github.com/jinzhu/gorm"
	protos "goim-pro/api/protos/salty"
	. "goim-pro/internal/app/constants"
	"goim-pro/internal/app/models"
	. "goim-pro/internal/app/repos/contact"
	. "goim-pro/internal/app/repos/user"
	"goim-pro/internal/app/services/converters"
	mysqlsrv "goim-pro/pkg/db/mysql"
	redsrv "goim-pro/pkg/db/redis"
	"goim-pro/pkg/errors"
	"goim-pro/pkg/http"
	"goim-pro/pkg/logs"
	"goim-pro/pkg/utils"
	"strings"
)

var (
	logger  = logs.GetLogger("INFO")
	myRedis *redsrv.BaseClient
	mysqlDB *gorm.DB

	userRepo    IUserRepo
	contactRepo IContactRepo
)

type ContactService struct {
}

func New() *ContactService {
	myRedis = redsrv.NewRedisConnection().GetRedisClient()
	mysqlDB = mysqlsrv.NewMysqlConnection().GetMysqlInstance()

	userRepo = NewUserRepo(mysqlDB)
	contactRepo = NewContactRepo(mysqlDB)

	return &ContactService{}
}

// RequestContact: request add contact
func (cs *ContactService) RequestContact(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
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

	_, err = userRepo.FindByUserId(contactId)
	if err != nil {
		if err == errmsg.ErrInvalidUserId {
			resp.Code = http.StatusBadRequest
			resp.Message = errmsg.ErrInvalidContact.Error()
			return
		}
		logger.Errorf("find contact by userId error: %s", err.Error())
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		return
	}

	isExists, err := contactRepo.IsContactExists(userId, contactId)
	if err != nil {
		logger.Errorf("checking contact error: %s", err.Error())
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		return
	}
	if isExists {
		resp.Code = http.StatusBadRequest
		resp.Message = errmsg.ErrContactAlreadyExists.Error()
		return
	}
	// TODO: cache in redis, should replace to Push notification server
	ctKey := fmt.Sprintf("CT-REQ-%s-%s", userId, contactId)
	cacheContent := make(map[string]interface{})
	cacheContent["contactId"] = contactId
	cacheContent["reason"] = requestReason
	jsonStr, err := utils.TransformMapToJSONString(cacheContent)
	if err != nil {
		logger.Errorf("redis operations error, transform map error: %s", err.Error())
	}
	err = myRedis.Set(ctKey, jsonStr, DefaultExpiresTime)

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
func (cs *ContactService) RefusedContact(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
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

	_, err = userRepo.FindByUserId(contactId)
	if err != nil {
		if err == errmsg.ErrInvalidUserId {
			resp.Code = http.StatusBadRequest
			resp.Message = errmsg.ErrInvalidContact.Error()
			return
		}
		logger.Errorf("find contact by userId error: %s", err.Error())
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		return
	}

	isExists, err := contactRepo.IsContactExists(userId, contactId)
	if err != nil {
		logger.Errorf("checking contact error: %s", err.Error())
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		return
	}
	if isExists {
		resp.Code = http.StatusBadRequest
		resp.Message = errmsg.ErrContactAlreadyExists.Error()
		return
	}

	// TODO: cache in redis, should replace to Push notification server
	ctKey := fmt.Sprintf("CT-REF-%s-%s", userId, contactId)
	cacheContent := make(map[string]interface{})
	cacheContent["contactId"] = contactId
	cacheContent["reason"] = refusedReason

	jsonStr, err := utils.TransformMapToJSONString(cacheContent)
	if err != nil {
		logger.Errorf("redis operations error, transform map error: %s", err.Error())
	}
	err = myRedis.Set(ctKey, jsonStr, DefaultExpiresTime)

	if err != nil {
		logger.Errorf("redis cache error: %s", err.Error())
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		return
	}
	return
}

// accept contact request
func (cs *ContactService) AcceptContact(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
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
		resp.Message = errmsg.ErrInvalidParameters.Error()
		return
	}
	if strings.EqualFold(userId, contactId) {
		resp.Code = http.StatusBadRequest
		resp.Message = errmsg.ErrIllegalOperation.Error()
		return
	}

	_, err = userRepo.FindByUserId(contactId)
	if err != nil {
		if err == errmsg.ErrInvalidUserId {
			resp.Code = http.StatusBadRequest
			resp.Message = errmsg.ErrInvalidContact.Error()
			return
		}
		logger.Errorf("find contact by userId error: %s", err.Error())
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		return
	}

	isExists, err := contactRepo.IsContactExists(userId, contactId)
	if err != nil {
		logger.Errorf("checking contact error: %s", err.Error())
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		return
	}
	if isExists {
		resp.Code = http.StatusBadRequest
		resp.Message = errmsg.ErrContactAlreadyExists.Error()
		return
	}

	// TODO: cache in redis, should replace to Push notification server
	ctKey := fmt.Sprintf("CT-ACP-%s-%s", userId, contactId)
	err = myRedis.Set(ctKey, contactId, DefaultExpiresTime)
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

// delete contact
func (cs *ContactService) DeleteContact(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	resp, _ = utils.NewGRPCResp(http.StatusOK, nil, "")

	var err error
	var delContactReq protos.DeleteContactReq
	if err = utils.UnmarshalGRPCReq(req, &delContactReq); err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = err.Error()
		logger.Errorf(`unmarshal error: %v`, err)
		return
	}

	userId := req.GetToken()
	contactId := strings.Trim(delContactReq.GetUserId(), "")

	if contactId == "" {
		resp.Code = http.StatusBadRequest
		resp.Message = errmsg.ErrInvalidParameters.Error()
		return
	}
	if strings.EqualFold(userId, contactId) {
		resp.Code = http.StatusBadRequest
		resp.Message = errmsg.ErrIllegalOperation.Error()
		return
	}

	isExists, err := contactRepo.IsContactExists(userId, contactId)
	if err != nil {
		logger.Errorf("checking contact error: %s", err.Error())
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		return
	}
	if !isExists {
		resp.Code = http.StatusBadRequest
		resp.Message = errmsg.ErrContactNotExists.Error()
		return
	}

	go func() {
		// TODO: cache in redis, should replace to Push notification server
		ctKey := fmt.Sprintf("CT-DEL-%s-%s", userId, contactId)
		err = myRedis.Set(ctKey, contactId, DefaultExpiresTime)
		if err != nil {
			// TODO: should log down exception information
			logger.Errorf("redis cache error: %s", err.Error())
		}
	}()

	if err = handleDeleteContact(cs, userId, contactId); err != nil {
		logger.Errorf("remove contact error: %s", err.Error())
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		return
	}
	resp.Message = "contact deleted successfully"
	return
}

// update contact remark profile
func (cs *ContactService) UpdateRemarkInfo(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	resp, _ = utils.NewGRPCResp(http.StatusOK, nil, "")

	var err error
	var updateRemarkInfoReq protos.UpdateRemarkInfoReq
	if err = utils.UnmarshalGRPCReq(req, &updateRemarkInfoReq); err != nil {
		resp.Code = http.StatusBadRequest
		resp.Message = err.Error()
		logger.Errorf(`unmarshal error: %v`, err)
		return
	}

	userId := req.GetToken()
	contactId := strings.Trim(updateRemarkInfoReq.GetUserId(), "")
	contactRemark := updateRemarkInfoReq.GetRemarkInfo()

	if contactId == "" {
		resp.Code = http.StatusBadRequest
		resp.Message = errmsg.ErrInvalidParameters.Error()
		return
	}

	isExists, err := contactRepo.IsContactExists(userId, contactId)
	if err != nil {
		logger.Errorf("checking contact error: %s", err.Error())
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		return
	}
	if !isExists {
		resp.Code = http.StatusBadRequest
		resp.Message = errmsg.ErrContactNotExists.Error()
		return
	}

	go func() {
		// TODO: cache in redis, should replace to Push notification server
		ctKey := fmt.Sprintf("CT-UDT-%s-%s", userId, contactId)
		err = myRedis.Set(ctKey, contactId, DefaultExpiresTime)
		if err != nil {
			// TODO: should log down exception information
			logger.Errorf("redis cache error: %s", err.Error())
		}
	}()

	if err = handleUpdateContactRemark(cs, userId, contactId, contactRemark); err != nil {
		logger.Errorf("update contact remark error: %s", err.Error())
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		return
	}
	resp.Message = "contact remark profile updated successfully"
	return
}

// query user's contact list
func (cs *ContactService) GetContacts(ctx context.Context, req *protos.GrpcReq) (resp *protos.GrpcResp, gRPCErr error) {
	resp, _ = utils.NewGRPCResp(http.StatusOK, nil, "")
	userId := req.GetToken()

	// TODO: should consider the black list function
	criteria := map[string]interface{}{
		"UserId": userId,
	}
	contacts, err := contactRepo.FindAll(criteria)
	if err != nil {
		logger.Errorf("query user contacts error: %s", err.Error())
		resp.Code = http.StatusInternalServerError
		resp.Message = err.Error()
		return
	}

	getContactsResp := &protos.GetContactsResp{
		Contacts: converters.ConvertEntity2ProtoForContacts(contacts),
	}

	resp.Data, err = utils.MarshalMessageToAny(getContactsResp)
	if err != nil {
		logger.Errorf("[get contacts] response marshal message error: %s", err.Error())
	}

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

func handleAcceptContact(cs *ContactService, userId string, contactId string) (err error) {
	newContact1 := &models.Contact{
		UserId:    userId,
		ContactId: contactId,
	}
	newContact2 := &models.Contact{
		UserId:    contactId,
		ContactId: userId,
	}

	err = contactRepo.InsertContacts(newContact1, newContact2)
	return err
}

func handleDeleteContact(cs *ContactService, userId string, contactId string) (err error) {
	tx := mysqlDB.Begin()
	if err = contactRepo.RemoveContactsByIds(userId, contactId); err != nil {
		logger.Errorf("remove contact err: %s", err.Error())
		tx.Rollback()
		return
	}
	if err = contactRepo.RemoveContactsByIds(contactId, userId); err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
	return err
}

func handleUpdateContactRemark(cs *ContactService, userId string, contactId string, pbProfile *protos.ContactRemark) (err error) {
	criteria := map[string]interface{}{
		"UserId":    userId,
		"ContactId": contactId,
	}

	remarkProfile := converters.ConvertProto2EntityForRemarkProfile(pbProfile)
	updateMap := utils.TransformStructToMap(remarkProfile)

	err = contactRepo.FindOneAndUpdateRemark(criteria, updateMap)

	return
}
