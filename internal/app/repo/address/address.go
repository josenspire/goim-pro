package address

import "github.com/jinzhu/gorm"

type IAddress interface {
	QueryUserAddressList(userID string) []interface{}
}

type AddressImpl struct {

}

func New(*gorm.DB) IAddress {
	return &AddressImpl{}
}

func (ad *AddressImpl) QueryUserAddressList(userID string) []interface{} {
	return make([]interface{}, 4)
}
