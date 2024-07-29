package dao

import (
	"github.com/assimon/luuu/config"
	"github.com/assimon/luuu/util/log"
	"github.com/gookit/color"
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// SqliteInit 数据库初始化
func SqliteInit() {
	var err error
	dbname := viper.GetString("sqlite_database_filename")
	if len(dbname) == 0 {
		dbname = ".db"
	}
	Mdb, err = gorm.Open(sqlite.Open(dbname), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   viper.GetString("sqlite_table_prefix"),
			SingularTable: true,
		},
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		panic(err)
	}
	if config.AppDebug {
		Mdb = Mdb.Debug()
	}
	sqlDB, err := Mdb.DB()
	if err != nil {
		color.Red.Printf("[store_db] sqlite get DB,err=%s\n", err)
		panic(err)
	}
	// sqlDB.SetMaxIdleConns(viper.GetInt("sqlite_max_idle_conns"))
	// sqlDB.SetMaxOpenConns(viper.GetInt("sqlite_max_open_conns"))
	// sqlDB.SetConnMaxLifetime(time.Hour * time.Duration(viper.GetInt("sqlite_max_life_time")))
	err = sqlDB.Ping()
	if err != nil {
		color.Red.Printf("[store_db] sqlite connDB err:%s", err.Error())
		panic(err)
	}
	log.Sugar.Debug("[store_db] sqlite connDB success")
}
