package address

type IAddress interface {
	QueryUserAddressList(userId string) []interface{}
}

type AddressImpl struct {

}

func New() IAddress {
	return &AddressImpl{}
}

func (ad *AddressImpl) QueryUserAddressList(userId string) []interface{} {
	return make([]interface{}, 4)
}
