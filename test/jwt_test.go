package test

import (
	"github.com/dgrijalva/jwt-go"
	"online-judge/utils"
	"testing"
)

type UserClaims struct {
	Identity string `json:"Identity"`
	Name     string `json:"name"`
	jwt.StandardClaims
}

var myKey = []byte("gin-gorm-oj-key")

// 生成Token
func TestGenerateToken(t *testing.T) {
	userClaims := &UserClaims{
		Identity:       "user_1",
		Name:           "Get",
		StandardClaims: jwt.StandardClaims{},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)
	tokenString, err := token.SignedString(myKey)
	if err != nil {
		t.Fatal(err)
	}

	// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZGVudGl0eSI6InVzZXJfMSIsIm5hbWUiOiJHZXQifQ.I6cd7qGZ31QtZ3FUZGZfwcnXiCPA_wBrDcWhfHk6Cuc
	utils.DPrintln(tokenString)
}

// 解析Token
func TestAnalyseToken(t *testing.T) {
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZGVudGl0eSI6InVzZXJfMSIsIm5hbWUiOiJHZXQifQ.I6cd7qGZ31QtZ3FUZGZfwcnXiCPA_wBrDcWhfHk6Cuc"
	userClaim := new(UserClaims)
	claims, err := jwt.ParseWithClaims(tokenString, userClaim, func(token *jwt.Token) (interface{}, error) {
		return myKey, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if claims.Valid {
		utils.DPrintln(userClaim)
	}
}
