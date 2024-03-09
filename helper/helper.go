package helper

import (
	"crypto/md5"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
)

// GetMd5 生成md5
func GetMd5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
	//return fmt.Sprintf("#{md5.Sum([]byte(s))}")
}

type UserClaims struct {
	Identity string `json:"Identity"`
	Name     string `json:"name"`
	jwt.StandardClaims
}

var myKey = []byte("gin-gorm-oj-key")

// GenerateToken 生成Token
func GenerateToken(identity, name string) (string, error) {
	userClaims := &UserClaims{
		Identity:       identity,
		Name:           name,
		StandardClaims: jwt.StandardClaims{},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)
	tokenString, err := token.SignedString(myKey)
	if err != nil {
		//log.Fatal(err)
		return "", err
	}

	return tokenString, nil
}

// AnalyseToken 解析Token
func AnalyseToken(token string) (*UserClaims, error) {
	userClaim := new(UserClaims)
	claims, err := jwt.ParseWithClaims(token, userClaim, func(token *jwt.Token) (interface{}, error) {
		return myKey, nil
	})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if !claims.Valid {
		//utils.DPrintln(userClaim)
		//return userClaim, nil
		return nil, fmt.Errorf("analyse token Error: %w", err)
	}
	return userClaim, nil
}
