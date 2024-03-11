package helper

import (
	"crypto/md5"
	"crypto/tls"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/jordan-wright/email"
	uuid "github.com/satori/go.uuid"
	"math/rand"
	"net/smtp"
	"os"
	"strconv"
	"time"
)

// GetMd5 生成md5
func GetMd5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
	//return fmt.Sprintf("#{md5.Sum([]byte(s))}")
}

type UserClaims struct {
	Identity string `json:"Identity"`
	Name     string `json:"name"`
	IsAdmin  int    `json:"is_admin"`
	jwt.StandardClaims
}

var myKey = []byte("gin-gorm-oj-key")

// GenerateToken 生成Token
func GenerateToken(identity, name string, isAdmin int) (string, error) {
	userClaims := &UserClaims{
		Identity:       identity,
		Name:           name,
		IsAdmin:        isAdmin,
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
		return nil, err
	}

	if !claims.Valid {
		//utils.DPrintln(userClaim)
		//return userClaim, nil
		return nil, fmt.Errorf("analyse token Error: %w", err)
	}
	return userClaim, nil
}

func SendCode(toUserEmail, code string) error {
	e := email.NewEmail()
	e.From = "Get <staraino_o@163.com>"
	e.To = []string{toUserEmail}
	e.Subject = "发送验证码"
	e.HTML = []byte("<b>" + code + "</b>")
	return e.SendWithTLS("smtp.163.com:465", smtp.PlainAuth("", "staraino_o@163.com", "SIYICTUUUKWYNGKO", "smtp.163.com"),
		&tls.Config{ServerName: "smtp.163.com", InsecureSkipVerify: true})
}

func GenerateCode() (code string) {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 6; i++ {
		code += strconv.Itoa(rand.Intn(10))
	}
	return
}

// GetUUID 生成唯一码
func GetUUID() string {
	return uuid.NewV4().String()
}

func CodeSave(code []byte) (string, error) {
	dirName := "code/" + GetUUID()
	path := dirName + "/main.go"
	err := os.Mkdir(dirName, 0777)
	if err != nil {
		return "", err
	}
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	f.Write(code)
	defer f.Close()
	return path, nil
}
