package dao

import (
	"time"

	"github.com/assimon/luuu/config"
	"github.com/assimon/luuu/util/log"
	"github.com/gookit/color"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// MysqlInit 数据库初始化
func MysqlInit() error {
	var err error
	Mdb, err = gorm.Open(mysql.Open(config.MysqlDns), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   viper.GetString("mysql_table_prefix"),
			SingularTable: true,
		},
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		color.Red.Printf("[store_db] mysql open DB,err=%s\n", err)
		// panic(err)
		return err
	}
	if config.AppDebug {
		Mdb = Mdb.Debug()
	}
	sqlDB, err := Mdb.DB()
	if err != nil {
		color.Red.Printf("[store_db] mysql get DB,err=%s\n", err)
		// panic(err)
		time.Sleep(10 * time.Second)
		MysqlInit()
		return err
	}
	sqlDB.SetMaxIdleConns(viper.GetInt("mysql_max_idle_conns"))
	sqlDB.SetMaxOpenConns(viper.GetInt("mysql_max_open_conns"))
	sqlDB.SetConnMaxLifetime(time.Hour * time.Duration(viper.GetInt("mysql_max_life_time")))
	err = sqlDB.Ping()
	if err != nil {
		color.Red.Printf("[store_db] mysql connDB err:%s", err.Error())
		// panic(err)
		time.Sleep(10 * time.Second)
		MysqlInit()
		return err
	}
	log.Sugar.Debug("[store_db] mysql connDB success")
	return nil
}
