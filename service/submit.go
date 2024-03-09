package service

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"online-judge/define"
	"online-judge/models"
	"strconv"
)

// GetSubmitList
// @Tags 公共方法
// @Summary 提交列表
// @Param page query int false "page"
// @Param size query int false "size"
// @Param problem_identity query string false "problem_identity"
// @Param user_identity query string false "user_identity"
// @Param status query int false "status"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /submit-list [get]
func GetSubmitList(ctx *gin.Context) {
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", define.DefaultSize))
	page, err := strconv.Atoi(ctx.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		log.Println("Get Problem List Page strconv Error:", err)
		return
	}
	var count int64
	list := make([]*models.SubmitBasic, 0)

	page = (page - 1) * size
	problemIdentity := ctx.Query("problem_identity")
	userIdentity := ctx.Query("user_identity")
	status, _ := strconv.Atoi(ctx.Query("status"))
	tx := models.GetSubmitList(problemIdentity, userIdentity, status)

	err = tx.Count(&count).Offset(page).Limit(size).Find(&list).Error
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Get Submit List Error:" + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"list":  list,
			"count": count,
		},
	})
}
