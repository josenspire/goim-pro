package contactsrv

import (
	"github.com/jinzhu/gorm"
	protos "goim-pro/api/protos/salty"
	"goim-pro/internal/app/models"
	. "goim-pro/internal/app/models/errors"
	. "goim-pro/internal/app/repos/contact"
	. "goim-pro/internal/app/repos/notification"
	. "goim-pro/internal/app/repos/user"
	"goim-pro/internal/app/services/converters"
	mysqlsrv "goim-pro/pkg/db/mysql"
	"goim-pro/pkg/errors"
	"goim-pro/pkg/logs"
	"goim-pro/pkg/utils"
)

var (
	logger  = logs.GetLogger("INFO")
	mysqlDB *gorm.DB

	userRepo         IUserRepo
	contactRepo      IContactRepo
	notificationRepo INotificationRepo
)

type IContactService interface {
	RequestContact(userId, contactId, reqReason string) (tErr *TError)
	RefusedContact(userId, contactId, refusedReason string) (tErr *TError)
	AcceptContact(userId, contactId string) (tErr *TError)
	DeleteContact(userId, contactId string) (tErr *TError)
	UpdateRemarkInfo(userId, contactId string, contactRemark *protos.ContactRemark) (tErr *TError)
	GetContacts(userId string) (contacts []models.Contact, tErr *TError)
	GetContactOperationMessageList(userId string, maxMessageTime int64) (contactOptsList []models.Notification, tErr *TError)
}

type ContactService struct {
}

func New() IContactService {
	mysqlDB = mysqlsrv.NewMysql()

	userRepo = NewUserRepo(mysqlDB)
	contactRepo = NewContactRepo(mysqlDB)
	notificationRepo = NewNotificationRepo(mysqlDB)

	return &ContactService{}
}

// RequestContact: request add contact
func (cs *ContactService) RequestContact(userId, contactId, reqReason string) (tErr *TError) {
	contact, err := userRepo.FindByUserId(contactId)
	if err != nil {
		logger.Errorf("find contact by userId error: %s", err.Error())
		return NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, err)
	}
	if contact == nil {
		return NewTError(protos.StatusCode_STATUS_BAD_REQUEST, errmsg.ErrInvalidContact)
	}

	isExists, err := contactRepo.IsContactExists(userId, contactId)
	if err != nil {
		logger.Errorf("checking contact error: %s", err.Error())
		return NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, err)
	}
	if isExists {
		return NewTError(protos.StatusCode_STATUS_BAD_REQUEST, errmsg.ErrContactAlreadyExists)
	}

	msgContent := make(map[string]interface{})
	msgContent["addReason"] = reqReason
	msgContent["rejectReason"] = ""
	msgContent["isNeedRemind"] = true
	msgContent["operationType"] = protos.ContactOperationMessage_ACCEPT_ACTIVE
	jsonStr, _ := utils.TransformMapToJSONString(msgContent)
	err = dispatchNotificationMessage(userId, contactId, protos.ContactOperationMessage_ACCEPT_ACTIVE, jsonStr, 0)
	if err != nil {
		return NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, err)
	}
	return nil
}

// refused request contact
func (cs *ContactService) RefusedContact(userId, contactId, refusedReason string) (tErr *TError) {
	contact, err := userRepo.FindByUserId(contactId)
	if err != nil {
		logger.Errorf("find contact by userId error: %s", err.Error())
		return NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, err)
	}
	if contact == nil {
		return NewTError(protos.StatusCode_STATUS_BAD_REQUEST, errmsg.ErrInvalidContact)
	}

	isExists, err := contactRepo.IsContactExists(userId, contactId)
	if err != nil {
		logger.Errorf("checking contact error: %s", err.Error())
		return NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, err)
	}
	if isExists {
		return NewTError(protos.StatusCode_STATUS_BAD_REQUEST, errmsg.ErrContactAlreadyExists)
	}

	msgContent := make(map[string]interface{})
	msgContent["addReason"] = ""
	msgContent["rejectReason"] = refusedReason
	msgContent["isNeedRemind"] = true
	msgContent["operationType"] = protos.ContactOperationMessage_REJECT_PASSIVE // 推送给被拒绝用户
	jsonStr, _ := utils.TransformMapToJSONString(msgContent)
	err = dispatchNotificationMessage(userId, contactId, protos.ContactOperationMessage_REJECT_PASSIVE, jsonStr, 0)
	if err != nil {
		return NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, err)
	}
	return
}

// accept contact request
func (cs *ContactService) AcceptContact(userId, contactId string) (tErr *TError) {
	contact, err := userRepo.FindByUserId(contactId)
	if err != nil {
		logger.Errorf("find contact by userId error: %s", err.Error())
		return NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, err)
	}
	if contact == nil {
		return NewTError(protos.StatusCode_STATUS_BAD_REQUEST, errmsg.ErrInvalidContact)
	}

	isExists, err := contactRepo.IsContactExists(userId, contactId)
	if err != nil {
		logger.Errorf("checking contact error: %s", err.Error())
		return NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, err)
	}
	if isExists {
		return NewTError(protos.StatusCode_STATUS_BAD_REQUEST, errmsg.ErrContactAlreadyExists)
	}

	if err = handleAcceptContact(userId, contactId); err != nil {
		logger.Errorf("insert contact error: %s", err.Error())
		return NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, err)
	}

	msgContent := make(map[string]interface{})
	msgContent["addReason"] = ""
	msgContent["rejectReason"] = ""
	msgContent["isNeedRemind"] = true
	msgContent["operationType"] = protos.ContactOperationMessage_ACCEPT_PASSIVE
	jsonStr, _ := utils.TransformMapToJSONString(msgContent)
	err = dispatchNotificationMessage(userId, contactId, protos.ContactOperationMessage_REJECT_PASSIVE, jsonStr, 0) // 推送给被接受的用户
	err = dispatchNotificationMessage(contactId, userId, protos.ContactOperationMessage_REJECT_ACTIVE, jsonStr, 1)  // 推送给主动接受的用户
	if err != nil {
		return NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, err)
	}

	return
}

// delete contact
func (cs *ContactService) DeleteContact(userId, contactId string) (tErr *TError) {
	isExists, err := contactRepo.IsContactExists(userId, contactId)
	if err != nil {
		logger.Errorf("checking contact error: %s", err.Error())
		return NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, err)
	}
	if !isExists {
		return NewTError(protos.StatusCode_STATUS_BAD_REQUEST, errmsg.ErrContactNotExists)
	}

	if err = handleDeleteContact(userId, contactId); err != nil {
		logger.Errorf("remove contact error: %s", err.Error())
		return NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, err)
	}

	msgContent := make(map[string]interface{})
	msgContent["addReason"] = ""
	msgContent["rejectReason"] = ""
	msgContent["isNeedRemind"] = true
	msgContent["operationType"] = protos.ContactOperationMessage_DELETE_ACTIVE
	jsonStr, _ := utils.TransformMapToJSONString(msgContent)
	err = dispatchNotificationMessage(userId, contactId, protos.ContactOperationMessage_DELETE_ACTIVE, jsonStr, 0) // 推送给当前用户
	if err != nil {
		return NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, err)
	}

	return
}

// update contact remark profile
func (cs *ContactService) UpdateRemarkInfo(userId, contactId string, contactRemark *protos.ContactRemark) (tErr *TError) {
	isExists, err := contactRepo.IsContactExists(userId, contactId)
	if err != nil {
		logger.Errorf("checking contact error: %s", err.Error())
		return NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, err)
	}
	if !isExists {
		return NewTError(protos.StatusCode_STATUS_BAD_REQUEST, errmsg.ErrContactNotExists)
	}

	if err = handleUpdateContactRemark(userId, contactId, contactRemark); err != nil {
		logger.Errorf("update contact remark error: %s", err.Error())
		return NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, err)
	}
	return
}

// query user's contact list
func (cs *ContactService) GetContacts(userId string) (contacts []models.Contact, tErr *TError) {
	// TODO: should consider the black list function
	criteria := map[string]interface{}{
		"userId": userId,
	}
	contacts, err := contactRepo.FindAll(criteria)
	if err != nil {
		logger.Errorf("query user contacts error: %s", err.Error())
		return nil, NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, err)
	}

	return contacts, nil
}

func (cs *ContactService) GetContactOperationMessageList(userId string, maxMessageTime int64) (notifications []models.Notification, tErr *TError) {
	fromDateStr := utils.ParseTimestampToDateTimeStr(maxMessageTime, utils.MysqlDateTimeFormat)
	notifications, err := notificationRepo.FindAll(userId, fromDateStr)
	if err != nil {
		logger.Errorf("query user contacts error: %s", err.Error())
		return nil, NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, err)
	}
	return notifications, nil
}

func handleAcceptContact(userId string, contactId string) (err error) {
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

func handleDeleteContact(userId string, contactId string) (err error) {
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

func handleUpdateContactRemark(userId string, contactId string, pbProfile *protos.ContactRemark) (err error) {
	criteria := map[string]interface{}{
		"UserId":    userId,
		"ContactId": contactId,
	}

	remarkProfile := converters.ConvertProto2EntityForRemarkProfile(pbProfile)
	updateMap := utils.TransformStructToMap(remarkProfile)

	err = contactRepo.FindOneAndUpdateRemark(criteria, updateMap)

	return
}

func dispatchNotificationMessage(userId, contactId string, operation protos.ContactOperationMessage_OperationType, strContent string, randomSeed int64) (err error) {
	messageId := utils.NewULID(randomSeed)
	notification := &models.Notification{
		NotifyId:         utils.NewULID(randomSeed),
		MessageId:        messageId,
		NotificationType: models.NotifyTypeSystem,
		FromUserId:       userId,
		ToUserId:         contactId,
		IsNeedRemind:     true,
		Status:           models.StsPending,
		Message: models.NotificationMessage{
			MessageId:    messageId,
			MsgType:      models.MsgTypeContactRequest,
			MsgOperation: int32(operation),
			MsgContent:   strContent,
		},
	}

	notifyMsg, err := notificationRepo.InsertOne(notification)
	if err != nil {
		logger.Errorf("insert new notification error: %s", err.Error())
		return err
	}

	// TODO: send notification
	logger.Info(notifyMsg)

	return nil
}
