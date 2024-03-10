package service

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"online-judge/define"
	"online-judge/helper"
	"online-judge/models"
	"strconv"
)

// GetCategoryList
// @Tags 管理员私有方法
// @Summary 分类列表
// @Param authorization header string true "authorization"
// @Param page query int false "page"
// @Param size query int false "size"
// @Param keyword query string false "keyword"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /category-list [get]
func GetCategoryList(ctx *gin.Context) {
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", define.DefaultSize))
	page, err := strconv.Atoi(ctx.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		log.Println("Get Problem List Page strconv Error:", err)
		return
	}
	var count int64
	page = (page - 1) * size

	keyword := ctx.Query("keyword")
	list := make([]*models.CategoryBasic, 0)
	err = models.DB.Model(new(models.CategoryBasic)).Where("name like ?", "%"+keyword+"%").
		Count(&count).Limit(size).Offset(page).Find(&list).Error
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "获取分类列表失败",
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

// CategoryCreate
// @Tags 管理员私有方法
// @Summary 分类创建
// @Param authorization header string true "authorization"
// @Param name formData string true "name"
// @Param parentId formData string false "parentId"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /category-create [post]
func CategoryCreate(ctx *gin.Context) {
	name := ctx.PostForm("name")
	parentId, _ := strconv.Atoi(ctx.PostForm("parentId"))
	category := &models.CategoryBasic{
		Identity: helper.GetUUID(),
		Name:     name,
		ParentId: parentId,
	}

	err := models.DB.Create(category).Error
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "创建分类失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "分类创建成功",
	})
}

// CategoryModify
// @Tags 管理员私有方法
// @Summary 分类修改
// @Param authorization header string true "authorization"
// @Param identity formData string true "identity"
// @Param name formData string true "name"
// @Param parentId formData string false "parentId"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /category-modify [put]
func CategoryModify(ctx *gin.Context) {
	identity := ctx.PostForm("identity")
	parentId, _ := strconv.Atoi(ctx.PostForm("parentId"))
	name := ctx.PostForm("name")
	if name == "" || identity == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数错误",
		})
	}
	category := &models.CategoryBasic{
		Identity: identity,
		Name:     name,
		ParentId: parentId,
	}

	err := models.DB.Model(new(models.CategoryBasic)).Where("identity = ?", identity).Updates(category).Error
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "分类修改失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "分类修改成功",
	})
}

// CategoryDelete
// @Tags 管理员私有方法
// @Summary 分类删除
// @Param authorization header string true "authorization"
// @Param identity query string true "identity"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /category-delete [delete]
func CategoryDelete(ctx *gin.Context) {
	identity := ctx.Query("identity")
	if identity == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不正确",
		})
		return
	}
	var cnt int64
	err := models.DB.Model(new(models.ProblemCategory)).Where("category_id = (SELECT id FROM category_basic WHERE identity = ? LIMIT 1)", identity).Count(&cnt).Error
	if err != nil {
		log.Println("Get ProblemCategory Error:", err)
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "获取分类关联的问题失败",
		})
		return
	}
	if cnt > 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "该分类下面已存在问题，不可删除",
		})
		return
	}
	err = models.DB.Where("identity = ?", identity).Delete(new(models.CategoryBasic)).Error
	if err != nil {
		log.Println("Delete CategoryBasic Error:", err)
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "删除失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}
