package main

import (
	"time"

	"github.com/dgrijalva/jwt-go"

	pb "github.com/bobcats/shipper/user-service/proto/user"
)

var (
	key = []byte("mySuperSecretKeyLol")
)

type CustomClaims struct {
	User *pb.User
	jwt.StandardClaims
}

type Authable interface {
	Decode(tokenString string) (*CustomClaims, error)
	Encode(user *pb.User) (string, error)
}

type TokenService struct {
	repo Repository
}

func (srv *TokenService) Decode(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

func (srv *TokenService) Encode(user *pb.User) (string, error) {
	expireToken := time.Now().Add(time.Hour * 72).Unix()

	claims := CustomClaims{
		user,
		jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    "go.micro.srv.user",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(key)
}
