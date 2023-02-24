module goim-pro

go 1.13

require (
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/davecgh/go-spew v1.1.1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/edwingeng/wuid v0.0.0-20191129013112-1bf3d98c028d
	github.com/go-redis/redis/v7 v7.0.0-beta.5
	github.com/go-sql-driver/mysql v1.4.1
	github.com/golang/protobuf v1.3.5
	github.com/jinzhu/gorm v1.9.11
	github.com/mozillazg/go-pinyin v0.18.0
	github.com/oklog/ulid v1.3.1
	github.com/pmezard/go-difflib v1.0.0
	github.com/sirupsen/logrus v1.4.2
	github.com/smartystreets/goconvey v1.6.4
	github.com/spf13/viper v1.6.1
	github.com/stretchr/objx v0.1.1
	github.com/stretchr/testify v1.4.0
	golang.org/x/net v0.7.0
	google.golang.org/genproto v0.0.0-20191206224255-0243a4be9c8f // indirect
	google.golang.org/grpc v1.25.1
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.28.1
