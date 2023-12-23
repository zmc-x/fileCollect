package system

import (
	"errors"
	"fileCollect/global"
	model "fileCollect/model/system"
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