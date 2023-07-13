package router

import "fileCollect/router/system"

type routerGroup struct {
	SystemRouter system.RouterGroup
}

// return the point to routerGroup structure's pointer
var RouterGroupApp = new(routerGroup)