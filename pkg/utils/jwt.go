package utils

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

const (
	SecretKey = "SaltyIM"
)

type MyClaims struct {
	Foo []byte `json:"foo"`
	jwt.StandardClaims
}

func NewToken(foo []byte) string {
	claims := MyClaims{
		Foo: foo,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(5)).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "salty_im",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		logger.Errorf("[jwt] signed string error: %s", err.Error())
		return ""
	}
	return tokenStr
}

func TokenVerify(tokenStr string) (bool, error) {
	logger.Infof("token string: %s", tokenStr)
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		logger.Errorf("authorized error: %s", err.Error())
		return false, err
	}
	if !token.Valid {
		logger.Warnf("unauthorized access to this resource")
		return false, nil
	}
	return true, nil
}
