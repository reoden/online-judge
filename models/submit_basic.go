package models

import "gorm.io/gorm"

type SubmitBasic struct {
	gorm.Model
	Identity        string        `gorm:"column:identity;type:varchar(36);" json:"identity"`
	ProblemIdentity string        `gorm:"column:problem_identity;type:varchar(36);" json:"problem_identity"`
	ProblemBasic    *ProblemBasic `gorm:"foreignKey:identity;references:problem_identity"` //关联问题基础表
	UserIdentity    string        `gorm:"column:user_identity;type:varchar(36);" json:"user_identity"`
	UserBasic       *UserBasic    `gorm:"foreignKey:identity;references:user_identity"` //关联用户基础表
	Path            string        `gorm:"column:path;type:varchar(255);" json:"path"`
	Status          int           `gorm:"column:status;type:tinyint(1);" json:"status"` //[-1-pending, 1-ac, 2-wa, 3-tle, 4-mle ,5-ce]
}

func (table *SubmitBasic) TableName() string {
	return "submit_basic"
}

func GetSubmitList(problemIdentity, userIdentity string, status int) *gorm.DB {
	tx := DB.Model(new(SubmitBasic)).Preload("ProblemBasic", func(db *gorm.DB) *gorm.DB {
		return db.Omit("content")
	}).Preload("UserBasic")

	if problemIdentity != "" {
		tx.Where("problem_identity = ? ", problemIdentity)
	}
	if userIdentity != "" {
		tx.Where("user_identity = ? ", userIdentity)
	}
	if status != 0 {
		tx.Where("status = ? ", status)
	}
	return tx
}
