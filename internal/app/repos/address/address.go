package address

type IAddress interface {
	QueryUserAddressList(userID string) []interface{}
}

type AddressImpl struct {

}

func New() IAddress {
	return &AddressImpl{}
}

func (ad *AddressImpl) QueryUserAddressList(userID string) []interface{} {
	return make([]interface{}, 4)
}
