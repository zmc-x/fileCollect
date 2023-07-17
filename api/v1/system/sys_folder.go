package system

import (
	"fileCollect/model/common/request"
	"fileCollect/model/common/response"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type SystemFolderApi struct{}

// router:/api/folder/createFolder
// method:post
func (sf *SystemFolderApi) CreateFolder(c *gin.Context) {
	var createFolderInfo request.CreateFolderInfo
	if err := c.ShouldBindJSON(&createFolderInfo); err != nil {
		processError(c, "api/v1/system/sys_folder.go CreateFolder method:", err)
		return
	}
	// get the storageRealPath
	storagePath, err := storageService.QueryStorageRealPath(createFolderInfo.StorageId)
	if err != nil {
		processError(c, "api/v1/system/sys_folder.go CreateFolder method:", err)
		return
	}
	// create folder in system
	err = os.Mkdir(filepath.Join(storagePath, createFolderInfo.Path, createFolderInfo.FolderName), 0644)
	if err != nil {
		processError(c, "api/v1/system/sys_folder.go CreateFolder method:", err)
		return
	}
	// database
	switch createFolderInfo.Path {
	case "/":
		err = folderService.CreateFolder(createFolderInfo.FolderName, createFolderInfo.StorageId, nil)
	default:
		err = folderService.CreateFolder(createFolderInfo.FolderName, createFolderInfo.StorageId, &createFolderInfo.ParentFolderID)
	}
	if err != nil {
		defer os.Remove(filepath.Join(storagePath, createFolderInfo.Path, createFolderInfo.FolderName))
		processError(c, "api/v1/system/sys_folder.go CreateFolder method:", err)
		return
	}
	response.Ok(c)
}


// router:/api/folder/deletefolder/
// method:delete
func (sf *SystemFolderApi) DeleteFolders(c *gin.Context) {
	var info request.DeleteFolderInfo
	if err := c.ShouldBindJSON(&info); err != nil {
		processError(c, "api/v1/system/sys_folder.go DeleteFolder method:", err)
		return
	}
	// get the storageRealPath
	storagePath, err := storageService.QueryStorageRealPath(info.StorageID)
	if err != nil {
		processError(c, "api/v1/system/sys_folder.go DeleteFolder method:", err)
		return
	}
	for _, v := range info.Folders {
		if err := folderService.DeleteFolder(v.FolderId, info.StorageID); err != nil {
			log.Println("api/v1/system/sys_folder.go DeleteFolder method:" + err.Error())
			continue
		}
		// delete system folder
		if err := os.Remove(filepath.Join(storagePath, info.Path, v.FolderName)); err != nil {
			log.Println("api/v1/system/sys_folder.go DeleteFolder method:" + err.Error())
		}
	}
	response.Ok(c)
}

// router:/api/folder/updatefolder
// method:post
func (sf *SystemFolderApi) UpdateFolder(c *gin.Context) {
	var info request.UpdateFolderInfo
	if err := c.ShouldBindJSON(&info); err != nil {
		processError(c, "api/v1/system/sys_folder.go UpdateFolder method:", err)
		return
	}
	// get the storageRealPath
	storagePath, err := storageService.QueryStorageRealPath(info.StorageId)
	if err != nil {
		processError(c, "api/v1/system/sys_folder.go UpdateFolder method:", err)
		return
	}
	folderPre := filepath.Join(storagePath, info.Path)
	nName, oName := filepath.Join(folderPre, info.FolderNewName), filepath.Join(folderPre, info.FolderName)
	// update system folder
	if err := os.Rename(oName, nName); err != nil {
		processError(c, "api/v1/system/sys_folder.go UpdateFolder method:", err)
		return
	}
	if err := folderService.UpdateFolderName(info.FolderId, info.StorageId, info.FolderNewName); err != nil {
		// restore
		defer os.Rename(nName, oName)
		processError(c, "api/v1/system/sys_folder.go UpdateFolder method:", err)
		return
	}
	response.Ok(c)
}