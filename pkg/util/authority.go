package util

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type UserClaims struct {
	UserID   int    `json:"userid"`
	UserName string `json:"username"`
	Email    string `json:"email"`
	jwt.StandardClaims
}

var myKey = []byte("Soft-ware-Engineering-Team-Project-Children-Math")

func GenerateToken(userID int, userName, email string) (string, error) {
	// 令牌的有限时间为 24 小时
	expireTime := time.Now().Add(time.Hour * 24)
	//expireTime := time.Now().Add(time.Second * 30)

	uerClaim := &UserClaims{
		UserID:   userID,
		UserName: userName,
		Email:    email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, uerClaim)
	tokenString, err := token.SignedString(myKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func AnalyseToken(tokenString string) (*UserClaims, error) {
	claims := new(UserClaims)
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return myKey, nil
	})
	if token != nil {
		if claims.ExpiresAt < time.Now().Unix() {
			return nil, errors.New("timeout")
		}
		if !token.Valid {
			return nil, errors.New("analyse token failed")
		}
	}
	return claims, err
}
