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
	tbl "goim-pro/internal/db"
	mysqlsrv "goim-pro/internal/db/mysql"
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
	GetContactList(userId string) (contacts []models.Contact, tErr *TError)
	GetContactOperationList(userId string, startDateTime, endDateTime int64) (contactOptsList []models.Notification, tErr *TError)
}

type ContactService struct {
}

func New() IContactService {
	mysqlDB = mysqlsrv.GetMysql()

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
	msgContent["operationType"] = protos.ContactOperationMessage_REQUEST
	jsonStr, _ := utils.TransformMapToJSONString(msgContent)
	err = dispatchNotificationMessage(userId, contactId, models.MsgTypeContactRequest, protos.ContactOperationMessage_REQUEST, jsonStr, "", 0)
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

	rErr := refreshNotificationMessageStatus(userId, contactId, models.MsgTypeContactRequest, protos.ContactOperationMessage_REJECT, models.StsReceived)
	if rErr != nil {
		return NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, rErr)
	}

	msgContent := make(map[string]interface{})
	msgContent["addReason"] = ""
	msgContent["rejectReason"] = refusedReason
	msgContent["isNeedRemind"] = true
	msgContent["operationType"] = protos.ContactOperationMessage_REJECT // 推送给被拒绝用户
	jsonStr, _ := utils.TransformMapToJSONString(msgContent)
	ntlErr := dispatchNotificationMessage(userId, contactId, models.MsgTypeContactRefuse, protos.ContactOperationMessage_REJECT, jsonStr, "", 0)
	if ntlErr != nil {
		return NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, ntlErr)
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

	rErr := refreshNotificationMessageStatus(userId, contactId, models.MsgTypeContactRequest, protos.ContactOperationMessage_ACCEPT, models.StsReceived)
	if rErr != nil {
		return NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, rErr)
	}

	msgContent := make(map[string]interface{})
	msgContent["addReason"] = ""
	msgContent["rejectReason"] = ""
	msgContent["isNeedRemind"] = true
	msgContent["operationType"] = protos.ContactOperationMessage_ACCEPT
	jsonStr, _ := utils.TransformMapToJSONString(msgContent)
	err = dispatchNotificationMessage(userId, contactId, models.MsgTypeContactAccept, protos.ContactOperationMessage_ACCEPT, jsonStr, "", 0) // 推送给被接受的用户
	err = dispatchNotificationMessage(contactId, userId, models.MsgTypeContactAccept, protos.ContactOperationMessage_ACCEPT, jsonStr, "", 1) // 推送给主动接受的用户
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

	//msgContent := make(map[string]interface{})
	//msgContent["addReason"] = ""
	//msgContent["rejectReason"] = ""
	//msgContent["isNeedRemind"] = true
	//msgContent["operationType"] = protos.ContactOperationMessage_DELETE_ACTIVE
	//jsonStr, _ := utils.TransformMapToJSONString(msgContent)
	//err = dispatchNotificationMessage(userId, contactId, models.MsgTypeContactDelete, protos.ContactOperationMessage_DELETE_ACTIVE, jsonStr, "", 0) // 推送给当前用户
	//if err != nil {
	//	return NewTError(protos.StatusCode_STATUS_INTERNAL_SERVER_ERROR, err)
	//}

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
func (cs *ContactService) GetContactList(userId string) (contacts []models.Contact, tErr *TError) {
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

func (cs *ContactService) GetContactOperationList(userId string, startDateTime, endDateTime int64) (notifications []models.Notification, tErr *TError) {
	var err error
	if startDateTime == endDateTime && startDateTime == -1 {
		notifications, err = notificationRepo.FindAllOperations(userId)
	} else {
		startDateStr := utils.ParseTimestampToDateTimeStr(startDateTime, utils.MysqlDateTimeFormat)
		endDateStr := utils.ParseTimestampToDateTimeStr(endDateTime, utils.MysqlDateTimeFormat)
		notifications, err = notificationRepo.FindAllByTimeRange(userId, startDateStr, endDateStr)
	}
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
	txContactRepo := NewContactRepo(tx)
	if err = txContactRepo.RemoveContactsByIds(userId, contactId); err != nil {
		logger.Errorf("remove contact err: %s", err.Error())
		tx.Rollback()
		return
	}
	if err = txContactRepo.RemoveContactsByIds(contactId, userId); err != nil {
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

func dispatchNotificationMessage(senderId, targetId, msgType string, operation protos.ContactOperationMessage_OperationType, strContent, strExtra string, randomSeed int64) (err error) {
	ntfCondition := make(map[string]interface{})
	ntfCondition["senderId"] = senderId
	ntfCondition["targetId"] = targetId
	ntfCondition["msgType"] = msgType
	msgUpdated := &models.NotificationMessage{
		MsgOperation: int32(operation),
		MsgContent:   strContent,
		MsgExtra:     strExtra,
	}
	ntf, err := refreshNotificationMessage(ntfCondition, msgUpdated)
	if err != nil {
		return err
	}

	messageId := utils.NewULID(randomSeed)
	notification := &models.Notification{
		NtfId:            utils.NewULID(randomSeed),
		MessageId:        messageId,
		SenderId:         senderId,
		TargetId:         targetId,
		ObjectName:       models.ObjectNameCustomContact, // 单用户通知
		NotificationType: models.NotifyTypeSystem,
		MsgType:          msgType,
		IsNeedRemind:     true,
		Status:           models.StsSent,
		Message: models.NotificationMessage{
			MessageId:    messageId,
			MsgOperation: int32(operation),
			MsgContent:   strContent,
			MsgExtra:     strExtra,
		},
	}

	if ntf == nil {
		notifyMsg, err := notificationRepo.InsertOne(notification)
		if err != nil {
			logger.Errorf("insert new notification error: %s", err.Error())
			return err
		}

		// TODO: send notification
		logger.Info(notifyMsg)
	}
	return nil
}

func refreshNotificationMessageStatus(targetId, senderId string, msgType string, operation protos.ContactOperationMessage_OperationType, status string) (err error) {
	ntfCondition := make(map[string]interface{})
	ntfCondition["senderId"] = senderId
	ntfCondition["targetId"] = targetId
	ntfCondition["msgType"] = msgType

	ntfUpdated := &models.Notification{
		IsNeedRemind: false,
		Status:       status,
	}
	var ntf = new(models.Notification)
	tx := mysqlDB.Begin()
	if err := tx.Table(tbl.TableNotifications).Where(ntfCondition).Update(ntfUpdated).First(ntf).Error; err != nil {
		tx.Rollback()
		logger.Errorf("error happened to update notification: %v", err)
		return err
	}
	msgCondition := make(map[string]interface{})
	msgCondition["messageId"] = ntf.MessageId
	msgUpdated := &models.NotificationMessage{
		MsgOperation: int32(operation),
	}
	if err := tx.Table(tbl.TableNotificationMsgs).Where(msgCondition).Update(msgUpdated).Error; err != nil {
		tx.Rollback()
		logger.Errorf("error happened to update notification msg: %v", err)
		return err
	}
	return tx.Commit().Error
}

func refreshNotificationMessage(ntfCondition interface{}, msgUpdated interface{}) (ntf *models.Notification, err error) {
	ntf, err = notificationRepo.FindOne(ntfCondition)
	if err != nil {
		logger.Errorf("error happened to update notification message: %v", err)
		return nil, err
	}
	if ntf == nil {
		return nil, nil
	}

	msgCondition := make(map[string]interface{})
	msgCondition["messageId"] = ntf.MessageId
	if err := mysqlDB.Table(tbl.TableNotificationMsgs).Where(msgCondition).Update(msgUpdated).Error; err != nil {
		logger.Errorf("error happened to update notification msg: %v", err)
		return ntf, err
	}
	return ntf, err
}
