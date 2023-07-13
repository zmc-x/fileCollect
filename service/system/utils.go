package system

import (
	"errors"
	"fileCollect/global"
)

// check whether the record exists
func checkRecordById(id uint, dataModel interface{}) error {
	db := global.MysqlDB
	tmp := db.Find(&dataModel, id)
	if tmp.RowsAffected == 0 {
		return errors.New("the database don't have this record")
	}
	return nil
}