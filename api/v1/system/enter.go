package system

import "fileCollect/service"

type ApiGroup struct {
	SystemFileApi
	SystemFolderApi
	SystemStorageApi
}

// import the service
var (
	storageService = service.ServiceGroupApp.SystemServiceGroup.StorageService
	fileService = service.ServiceGroupApp.SystemServiceGroup.FileService
	folderService = service.ServiceGroupApp.SystemServiceGroup.FolderService
)