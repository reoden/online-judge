package models

import (
	"gorm.io/gorm"
)

type ProblemBasic struct {
	gorm.Model
	// 问题表的唯一标识
	Identity string `gorm:"column:identity;type:varchar(36);" json:"identity"`
	//关联问题分类表
	ProblemCategories []*ProblemCategory `gorm:"foreignKey:problem_id;references:id"`
	// 文章的标题
	Title string `gorm:"column:title;type:varchar(255);" json:"title"`
	// 文章正文
	Content string `gorm:"column:content;type:text;" json:"content"`
	//最大运行时间
	MaxRunTime   int         `gorm:"column:max_runtime;type:int(11);" json:"max_runtime"`
	MaxMem       int         `gorm:"column:max_mem;type:int(11);" json:"max_mem"`
	TestCases    []*TestCase `gorm:"foreignKey:problem_identity;references:identity"`
	AcNumber     int64       `gorm:"column:ac_number;type:int(11);" json:"ac_number"`
	SubmitNumber int64       `gorm:"column:submit_number;type:int(11);" json:"submit_number"`
}

func (table *ProblemBasic) TableName() string {
	return "problem_basic"
}

func GetProblemList(keyword, categoryIdentity string) *gorm.DB {
	tx := DB.Model(new(ProblemBasic)).Preload("ProblemCategories").Preload("ProblemCategories.CategoryBasic").
		Where("title like ? OR content like ?", "%"+keyword+"%", "%"+keyword+"%")

	if categoryIdentity != "" {
		tx.Joins("RIGHT JOIN problem_category pc on pc.problem_id = problem_basic.id").
			Where("pc.category_id = (SELECT cb.id FROM category_basic cb WHERE cb.identity = ? )", categoryIdentity)
	}
	return tx
}
