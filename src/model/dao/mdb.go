package dao

import (
	"strings"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var Mdb *gorm.DB

func DBInit() error {
	dbType := viper.GetString("db_type")
	if strings.EqualFold(dbType, "postgres") {
		if err := PostgreSQLInit(); err != nil {
			return err
		}
	} else if strings.EqualFold(dbType, "sqlite") {
		if err := SqliteInit(); err != nil {
			return err
		}
	} else {
		if err := MysqlInit(); err != nil {
			return err
		}
	}

	MdbTableInit()
	return nil
}
