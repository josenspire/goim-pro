package address

import "goim-pro/internal/app/models"

type Address models.Address

type IAddress interface {
	QueryUserAddressList(userId string) []interface{}
}

func New() IAddress {
	return &Address{}
}

func (ad *Address) QueryUserAddressList(userId string) []interface{} {
	return make([]interface{}, 4)
}
