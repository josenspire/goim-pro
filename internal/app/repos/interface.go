package repos

type IGormRepo interface {
	InsertOne(target interface{}) (value interface{}, err error)
}
