package initialize

import (
	"fileCollect/global"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// return server config information
func InitConfig() *global.ServerConfig {
	type Config struct {
		Mysql struct {
			DatabaseName string `yaml:"databaseName"`
			Url string `yaml:"url"` 
			Port uint16 `yaml:"port"`
			Username string `yaml:"userName"`
			Password string `yaml:"passWord"`
			MaxIdleConns int `yaml:"maxIdleConns"`
			MaxOpenConns int `yaml:"maxOpenConns"`
			ConnMaxLifetime int `yaml:"connMaxLifetime"`
		}
		Service struct {
			Port	uint16 `yaml:"port"`
		}
		Redis struct {
			RedisAddr	string `yaml:"addr"`
			RedisPort	uint16 `yaml:"port"`
			RedisPasswd	string `yaml:"passWord"`
			RedisDb		uint8	`yaml:"db"`
		}
	}
	config := Config{}
	// read the config.yml
	data, err := os.ReadFile("config.yml")
	if err != nil {
		log.Fatal("Global config modul InitConfig function read the file error, this error is " + err.Error())
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatal("this yml file unmarshal error, this error is" + err.Error())
	}
	return &global.ServerConfig{
		DbUser: config.Mysql.Username,
		DbPasswd: config.Mysql.Password,
		DbUrl: config.Mysql.Url,
		DbPort: config.Mysql.Port,
		DbName: config.Mysql.DatabaseName,
		DbPoolMaxIdleConns: config.Mysql.MaxIdleConns,
		DbPoolConnMaxLifetime: config.Mysql.ConnMaxLifetime,
		DbPoolMaxOpenConns: config.Mysql.MaxOpenConns,
		GinPort: config.Service.Port,
		RedisAddr: config.Redis.RedisAddr,
		RedisPort: config.Redis.RedisPort,
		RedisPasswd: config.Redis.RedisPasswd,
		RedisDb: config.Redis.RedisDb,
	}
}