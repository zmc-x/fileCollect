package main

import (
	"fileCollect/global"
	"fileCollect/initialize"
	"fmt"
)

func main() {
	// initialize zap
	initialize.InitLogger()
	defer global.Logger.Sync()
	serverConfig := initialize.InitConfig()
	initialize.InitMysql(serverConfig)
	initialize.InitReids(serverConfig)
	// defer close the database connect
	defer global.SqlDb.Close()
	initialize.InitTable()
	r := initialize.Router()
	r.Run(fmt.Sprintf(":%d", serverConfig.GinPort))
}

