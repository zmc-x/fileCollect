package system

import (
	"errors"
	"fileCollect/global"
	model "fileCollect/model/system"
	"strings"
)

// get StorageId by storageKey
func getStorageId(storageKey string) (uint, error) {
	db := global.MysqlDB
	var res model.Storage
	if temp := db.Where("storage_url_name = ?", storageKey).Find(&res); temp.RowsAffected == 0 {
		return res.ID, errors.New("the corresponding record cannot be queried")
	}
	return res.ID, nil
}

// get FolderId by path and storageId
func getFolderId(path string, storageId uint) (uint, error) {
	db :=global.MysqlDB
	var res model.Folder
	folders := strings.Split(path, "/")
	lenFolders := len(folders)
	for i := 1; i < lenFolders && folders[i] != ""; i++ {
		if temp := db.Where("storage_id = ? and folder_name = ? and parent_folder_id = ?", storageId, folders[i], res.ID).Find(&res); temp.RowsAffected == 0 {
			return res.ID, errors.New("the corresponding record cannot be queried")
		}
	}
	return res.ID, nil
}