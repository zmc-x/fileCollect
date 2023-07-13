package v1

import "fileCollect/api/v1/system"

type apiGroup struct {
	SystemApiGroup system.ApiGroup
}

// return the pointer that point to apiGroup pointer
var ApiGroupApp = new(apiGroup)