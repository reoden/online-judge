package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"online-judge/define"
	"online-judge/models"
	"strconv"
)

// GetProblemList
// @Tags 公共方法
// @Summary 问题列表
// @Param page query int false "page"
// @Param size query int false "size"
// @Param keyword query string false "keyword"
// @Param category_identity query string false "category_identity"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /problem-list [get]
func GetProblemList(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		log.Println("Get Problem List Page strconv Error:", err)
		return
	}
	// page == 1 ===> offset 0base
	size, err := strconv.Atoi(ctx.DefaultQuery("size", define.DefaultSize))
	if err != nil {
		log.Println("Get Problem List size strconv Error:", err)
		return
	}

	page = (page - 1) * size
	keyword := ctx.Query("keyword")
	categoryIdentity := ctx.Query("category_identity")
	tx := models.GetProblemList(keyword, categoryIdentity)

	list := make([]*models.ProblemBasic, 0)
	var count int64
	err = tx.Count(&count).Omit("content").Offset(page).Limit(size).Find(&list).Error
	if err != nil {
		log.Println("Get Problem List Error:", err)
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

// GetProblemDetail
// @Tags 公共方法
// @Summary 问题详情
// @Param identity query string false "problem identity"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /problem-detail [get]
func GetProblemDetail(ctx *gin.Context) {
	identity := ctx.Query("identity")
	if identity == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "问题唯一标识不能为空",
		})
		return
	}

	data := new(models.ProblemBasic)
	err := models.DB.Where("identity = ?", identity).Preload("ProblemCategories").
		Preload("ProblemCategories.CategoryBasic").First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "当前问题不存在",
			})
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "Get Problem Detail Error:" + err.Error(),
			})
		}
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": data,
	})
}
