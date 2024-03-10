package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"online-judge/helper"
	"online-judge/models"
	"time"
)

// GetUserDetail
// @Tags 公共方法
// @Summary 用户详情
// @Param identity query string false "user identity"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /user-detail [get]
func GetUserDetail(ctx *gin.Context) {
	identity := ctx.Query("identity")
	if identity == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"coed": -1,
			"msg":  "用户唯一标识不能为空",
		})
		return
	}

	data := new(models.UserBasic)
	err := models.DB.Where("identity = ?", identity).Omit("password").First(&data).Error
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Get User Detail Error:" + err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": data,
	})
}

// Login
// @Tags 公共方法
// @Summary 用户登录
// @Param username formData string false "username"
// @Param password formData string false "password"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /login [post]
func Login(ctx *gin.Context) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	if username == "" || password == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "必填信息为空",
		})
	}
	password = helper.GetMd5(password)
	//utils.DPrintf("username = %s, password = %s\n", username, password)
	//print(username, password)
	data := new(models.UserBasic)
	err := models.DB.Where("name = ? AND password = ? ", username, password).First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "用户名或密码错误",
			})
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "Get UserBasic Error:" + err.Error(),
			})
		}
		return
	}

	token, err := helper.GenerateToken(data.Identity, data.Name)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Generate Token Error:" + err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"token": token,
		},
	})
}

// SendCode
// @Tags 公共方法
// @Summary 发送验证骂
// @Param email formData string true "email"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /send-code [post]
func SendCode(ctx *gin.Context) {
	email := ctx.PostForm("email")
	if email == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不正确",
		})
		return
	}

	code := "123456"
	models.RDB.Set(ctx, email, code, time.Second*300)
	err := helper.SendCode(email, code)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Send Code Error:" + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Success",
	})
}

// Register
// @Tags 公共方法
// @Summary 用户注册
// @Param email formData string true "email"
// @Param code formData string true "code"
// @Param name formData string true "name"
// @Param password formData string true "password"
// @Param phone formData string false "phone"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /register [post]
func Register(c *gin.Context) {
	email := c.PostForm("email")
	userCode := c.PostForm("code")
	name := c.PostForm("name")
	password := c.PostForm("password")
	phone := c.PostForm("phone")
	if email == "" || userCode == "" || name == "" || password == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不正确",
		})
		return
	}
	// 验证验证码是否正确
	sysCode, err := models.RDB.Get(c, email).Result()
	if err != nil {
		log.Printf("Get Code Error:%v \n", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "验证码不正确，请重新获取验证码",
		})
		return
	}
	if sysCode != userCode {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "验证码不正确",
		})
		return
	}
	// 判断邮箱是否已存在
	var cnt int64
	err = models.DB.Where("email = ?", email).Model(new(models.UserBasic)).Count(&cnt).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Get User Error:" + err.Error(),
		})
		return
	}
	if cnt > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "该邮箱已被注册",
		})
		return
	}

	// 数据的插入
	userIdentity := helper.GetUUID()
	data := &models.UserBasic{
		Identity: userIdentity,
		Name:     name,
		Password: helper.GetMd5(password),
		Phone:    phone,
		Email:    email,
	}
	err = models.DB.Create(data).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Crete User Error:" + err.Error(),
		})
		return
	}

	// 生成 token
	token, err := helper.GenerateToken(userIdentity, name)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Generate Token Error:" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"token": token,
		},
	})
}
