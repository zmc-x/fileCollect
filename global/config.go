package global

import (
	"database/sql"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ServerConfig struct {
	// database config
	DbName                string
	DbUser                string
	DbPasswd              string
	DbUrl                 string
	DbPort                uint16
	DbPoolMaxIdleConns    int
	DbPoolMaxOpenConns    int
	DbPoolConnMaxLifetime int
	// gin config
	GinPort uint16
	// redis config
	RedisAddr   string
	RedisPort   uint16
	RedisPasswd string
	RedisDb     uint8
}

// define the global variable
var (
	MysqlDB  *gorm.DB
	SqlDb    *sql.DB
	Rdb      *redis.Client
	RealPath map[string]string // storage's realpath cache
	Logger   *zap.SugaredLogger
)

const (
	Format = "2006-01-02 15:04:05(MST)"
)
