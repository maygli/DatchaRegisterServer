package deviceservice

import (
	"datcha/servercommon"
	"errors"
	"log"

	"github.com/golang-jwt/jwt/v5"
)

type JwtDeviceClaims struct {
	jwt.RegisteredClaims
	DeviceId uint `json:"device_id"`
}

func (server *DeviceService) GenerateToken(deviceId uint) (string, error) {
	claims := JwtDeviceClaims{
		jwt.RegisteredClaims{
			Subject: server.StateSubject,
			Issuer:  server.Issuer,
		},
		deviceId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(server.StateSecretKey))
	if err != nil {
		log.Println("Error. Generate token failed. Error: " + err.Error())
		return "", errors.New(servercommon.ERROR_INTERNAL)
	}
	return tokenStr, nil
}

func (server *DeviceService) ParseToken(tokenStr string) (uint, error) {
	claims := JwtDeviceClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(server.StateSecretKey), nil
	})
	if err != nil {
		return 0, errors.New(servercommon.ERROR_PARSE_DEVICE_TOKEN)
	}
	if !token.Valid {
		return 0, errors.New(servercommon.ERROR_PARSE_DEVICE_TOKEN)
	}
	subject, err := claims.GetSubject()
	if err != nil {
		return 0, errors.New(servercommon.ERROR_PARSE_DEVICE_TOKEN)
	}
	if subject != server.StateSubject {
		return 0, errors.New(servercommon.ERROR_PARSE_DEVICE_TOKEN)
	}
	return claims.DeviceId, nil
}
