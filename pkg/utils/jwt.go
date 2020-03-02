package utils

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

const (
	secretKey   = "SaltyIM"
	expiresTime = time.Hour * time.Duration(24*3) // 3 days
)

type MyClaims struct {
	Foo []byte `json:"foo"`
	jwt.StandardClaims
}

func NewToken(foo []byte) string {
	claims := MyClaims{
		Foo: foo,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expiresTime).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "salty_im",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(secretKey))
	if err != nil {
		logger.Errorf("[jwt] signed string error: %s", err.Error())
		return ""
	}
	return tokenStr
}

func TokenVerify(tokenStr string) (isValid bool, payload []byte, err error) {
	logger.Infof("token string: %s", tokenStr)
	claims := &MyClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (i interface{}, e error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		logger.Errorf("authorized error: %s", err.Error())
		return false, nil, err
	}
	if !token.Valid {
		logger.Warnf("unauthorized access to this resource")
		return false, nil, nil
	}
	return true, claims.Foo, nil
}
