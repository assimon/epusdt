package mdb

import (
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint64         `gorm:"column:id;primary_key" json:"id"`
	CreatedAt carbon.Time    `gorm:"column:created_at" json:"created_at"`
	UpdatedAt carbon.Time    `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
