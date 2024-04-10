package dao

import (
	"fmt"
	"time"

	"github.com/assimon/luuu/config"
	"github.com/assimon/luuu/util/log"
	"github.com/gookit/color"
	"github.com/spf13/viper"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// PostgreSQLInit 数据库初始化
func PostgreSQLInit() {
	var err error
	user := viper.GetString("postgres_user")
	pass := viper.GetString("postgres_passwd")
	host := viper.GetString("postgres_host")
	port := viper.GetString("postgres_port")
	dbname := viper.GetString("postgres_database")
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		user, pass, host, port, dbname,
	)
	Mdb, err = gorm.Open(postgres.Open(dsn),
		&gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   viper.GetString("postgres_table_prefix"),
				SingularTable: true,
			},
			Logger: logger.Default.LogMode(logger.Error),
		},
	)
	if err != nil {
		panic(err)
	}

	if config.AppDebug {
		Mdb = Mdb.Debug()
	}
	sqlDB, err := Mdb.DB()
	if err != nil {
		color.Red.Printf("[store_db] postgres get DB,err=%s\n", err)
		panic(err)
	}
	sqlDB.SetMaxIdleConns(viper.GetInt("postgres_max_idle_conns"))
	sqlDB.SetMaxOpenConns(viper.GetInt("postgres_max_open_conns"))
	sqlDB.SetConnMaxLifetime(time.Hour * time.Duration(viper.GetInt("postgres_max_life_time")))
	err = sqlDB.Ping()
	if err != nil {
		color.Red.Printf("[store_db] postgres connDB err:%s", err.Error())
		panic(err)
	}
	log.Sugar.Debug("[store_db] postgres connDB success")
}
