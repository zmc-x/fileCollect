package initialize

import (
	"fileCollect/global"
	model "fileCollect/model/system"
	"fileCollect/utils/zaplog"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)


func InitMysql(sc *global.ServerConfig) {
	url, port, user, passwd, databasename := sc.DbUrl, sc.DbPort, sc.DbUser, sc.DbPasswd, sc.DbName
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", user, passwd, url, port, databasename), // DSN data source name
		DefaultStringSize: 256, // string 类型字段的默认长度
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: false,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		zaplog.GetLogLevel(zaplog.FATAL, err.Error())
	}
	global.MysqlDB = db
	global.SqlDb, err = db.DB()
	if err != nil {
		zaplog.GetLogLevel(zaplog.FATAL, err.Error())
	}
	// set the connection pool
	global.SqlDb.SetMaxIdleConns(sc.DbPoolMaxIdleConns)
	global.SqlDb.SetMaxOpenConns(sc.DbPoolMaxOpenConns)
	global.SqlDb.SetConnMaxLifetime(time.Hour * time.Duration(sc.DbPoolConnMaxLifetime))
	zaplog.GetLogLevel(zaplog.INFO, "database connection successful")
}

// create table
func InitTable() {
	err := global.MysqlDB.AutoMigrate(&model.Storage{}, &model.Folder{}, &model.File{})
	if err != nil {
		zaplog.GetLogLevel(zaplog.FATAL, err.Error())
	}
	zaplog.GetLogLevel(zaplog.INFO, "database table initialization completed")
}
