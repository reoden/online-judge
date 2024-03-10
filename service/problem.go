package service

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"online-judge/define"
	"online-judge/helper"
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
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", define.DefaultSize))
	page, err := strconv.Atoi(ctx.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		log.Println("Get Problem List Page strconv Error:", err)
		return
	}
	var count int64
	list := make([]*models.ProblemBasic, 0)
	page = (page - 1) * size

	keyword := ctx.Query("keyword")
	categoryIdentity := ctx.Query("category_identity")
	tx := models.GetProblemList(keyword, categoryIdentity)

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

// ProblemCreate
// @Tags 管理员私有方法
// @Summary 问题创建
// @Param authorization header string true "authorization"
// @Param title formData string true "title"
// @Param content formData string true "content"
// @Param max_runtime formData int false "max_runtime"
// @Param max_mem formData int false "max_mem"
// @Param category_ids formData []string false "category_ids" collectionFormat(multi)
// @Param test_cases formData []string true "test_cases" collectionFormat(multi)
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/problem-create [post]
func ProblemCreate(ctx *gin.Context) {
	title := ctx.PostForm("title")
	content := ctx.PostForm("content")
	maxRunTime, _ := strconv.Atoi(ctx.PostForm("max_runtime"))
	maxMem, _ := strconv.Atoi(ctx.PostForm("max_mem"))
	categoryIds := ctx.PostFormArray("category_ids")
	testCases := ctx.PostFormArray("test_cases")
	if title == "" || content == "" || len(testCases) == 0 || len(categoryIds) == 0 || maxMem == 0 || maxRunTime == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不正确",
		})
		return
	}

	userIdentity := helper.GetUUID()
	data := models.ProblemBasic{
		Identity:   userIdentity,
		Title:      title,
		Content:    content,
		MaxRunTime: maxRunTime,
		MaxMem:     maxMem,
	}

	// 处理分类
	problemCategoryBasics := make([]*models.ProblemCategory, 0)
	for _, id := range categoryIds {
		categoryId, _ := strconv.Atoi(id)
		problemCategoryBasics = append(problemCategoryBasics, &models.ProblemCategory{
			ProblemId:  data.ID,
			CategoryId: uint(categoryId),
		})
	}

	data.ProblemCategories = problemCategoryBasics
	// 处理测试样例

	testCaseBasics := make([]*models.TestCase, 0)
	for _, testCase := range testCases {
		// {"input":"1 2\n", "output":"3\n"}
		caseMap := make(map[string]string)
		err := json.Unmarshal([]byte(testCase), &caseMap)
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "测试用例错误",
			})
			return
		}
		if _, ok := caseMap["input"]; !ok {
			ctx.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "测试用例错误",
			})
			return
		}

		if _, ok := caseMap["output"]; !ok {
			ctx.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "测试用例错误",
			})
			return
		}
		testCaseBasic := &models.TestCase{
			Identity:        helper.GetUUID(),
			ProblemIdentity: userIdentity,
			Input:           caseMap["input"],
			Output:          caseMap["output"],
		}
		testCaseBasics = append(testCaseBasics, testCaseBasic)
	}

	data.TestCases = testCaseBasics
	err := models.DB.Create(&data).Error
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "创建失败" + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"identity": userIdentity,
		},
	})
}

// ProblemModify
// @Tags 管理员私有方法
// @Summary 问题修改
// @Param authorization header string true "authorization"
// @Param identity formData string true "identity"
// @Param title formData string true "title"
// @Param content formData string true "content"
// @Param max_runtime formData int false "max_runtime"
// @Param max_mem formData int false "max_mem"
// @Param category_ids formData []string false "category_ids" collectionFormat(multi)
// @Param test_cases formData []string true "test_cases" collectionFormat(multi)
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/problem-modify [put]
func ProblemModify(ctx *gin.Context) {
	identity := ctx.PostForm("identity")
	title := ctx.PostForm("title")
	content := ctx.PostForm("content")
	maxRunTime, _ := strconv.Atoi(ctx.PostForm("max_runtime"))
	maxMem, _ := strconv.Atoi(ctx.PostForm("max_mem"))
	categoryIds := ctx.PostFormArray("category_ids")
	testCases := ctx.PostFormArray("test_cases")
	if identity == "" || title == "" || content == "" || len(testCases) == 0 || len(categoryIds) == 0 || maxMem == 0 || maxRunTime == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不正确",
		})
		return
	}

	if err := models.DB.Transaction(func(tx *gorm.DB) error {
		// 问题基础信息的保存
		problemBasic := &models.ProblemBasic{
			Identity:   identity,
			Title:      title,
			Content:    content,
			MaxRunTime: maxRunTime,
			MaxMem:     maxMem,
		}

		err := tx.Where("identity = ?", identity).Updates(problemBasic).Error
		if err != nil {
			return err
		}
		//查询问题详情
		err = tx.Where("identity = ?", identity).Find(problemBasic).Error
		if err != nil {
			return err
		}
		// 关联问题分类的更新
		err = tx.Where("problem_id = ?", problemBasic.ID).Delete(new(models.ProblemCategory)).Error
		if err != nil {
			return err
		}

		pcs := make([]*models.ProblemCategory, 0)
		for _, id := range categoryIds {
			categoryId, _ := strconv.Atoi(id)
			pcs = append(pcs, &models.ProblemCategory{
				ProblemId:  problemBasic.ID,
				CategoryId: uint(categoryId),
			})
		}

		err = tx.Create(&pcs).Error
		if err != nil {
			return err
		}
		// 关联测试案例的更新
		err = tx.Where("problem_identity = ?", identity).Delete(new(models.TestCase)).Error
		if err != nil {
			return err
		}

		tcs := make([]*models.TestCase, 0)
		for _, testCase := range testCases {
			// {"input":"1 2\n", "output":"3\n"}
			caseMap := make(map[string]string)
			err := json.Unmarshal([]byte(testCase), &caseMap)
			if err != nil {
				return err
			}
			if _, ok := caseMap["input"]; !ok {
				return errors.New("测试样例input错误")
			}
			if _, ok := caseMap["output"]; !ok {
				return errors.New("测试样例output错误")
			}
			testCaseBasic := &models.TestCase{
				Identity:        helper.GetUUID(),
				ProblemIdentity: identity,
				Input:           caseMap["input"],
				Output:          caseMap["output"],
			}
			tcs = append(tcs, testCaseBasic)
		}
		err = tx.Create(&tcs).Error
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "problem modify error" + err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "问题修改成功",
	})
}
