package service

import "fileCollect/service/system"

type serviceGroup struct {
	SystemServiceGroup system.ServiceGroup
}

// return point to serviceGroup's pointer
var ServiceGroupApp = new(serviceGroup)