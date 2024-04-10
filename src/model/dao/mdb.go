package dao

import (
	"strings"

	"github.com/spf13/viper"

	"gorm.io/gorm"
)

var Mdb *gorm.DB

func DBInit() {
	dbType := viper.GetString("db_type")
	if strings.EqualFold(dbType, "postgres") {
		PostgreSQLInit()
	} else {
		MysqlInit()
	}
	MdbTableInit()
}
