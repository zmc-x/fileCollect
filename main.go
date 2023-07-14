package main

import (
	"fileCollect/global"
	"fileCollect/initialize"
	"fmt"
)

func main() {
	serverConfig := initialize.InitConfig()
	initialize.InitMysql(serverConfig)
	// defer close the database connect
	defer global.SqlDb.Close()
	initialize.InitTable()
	r := initialize.Router()
	r.Run(fmt.Sprintf(":%d", serverConfig.GinPort))
}

