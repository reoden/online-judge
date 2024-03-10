package router

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "online-judge/docs"
	"online-judge/middlewares"
	"online-judge/service"
)

func Router() *gin.Engine {
	r := gin.Default()

	//Swagger 配置
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	// 路由规则

	// 公用方法
	// problem
	r.GET("/problem-list", service.GetProblemList)
	r.GET("/problem-detail", service.GetProblemDetail)

	// user
	r.GET("/user-detail", service.GetUserDetail)
	r.POST("/login", service.Login)
	r.POST("/send-code", service.SendCode)
	r.POST("/register", service.Register)

	// User rank
	r.GET("/rank-list", service.GetRankList)
	//submission
	r.GET("/submit-list", service.GetSubmitList)

	// 私有方法
	//创建问题
	r.POST("/problem-create", middlewares.AuthAdminCheck(), service.ProblemCreate)
	//分类列表
	r.GET("/category-list", middlewares.AuthAdminCheck(), service.GetCategoryList)
	//分类创建
	r.POST("/category-create", middlewares.AuthAdminCheck(), service.CategoryCreate)
	// 分类修改
	r.PUT("/category-modify", middlewares.AuthAdminCheck(), service.CategoryModify)
	// 分类删除
	r.DELETE("/category-delete", middlewares.AuthAdminCheck(), service.CategoryDelete)

	return r
}
