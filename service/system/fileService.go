package system

import (
	"fileCollect/global"
	model "fileCollect/model/system"
)

type FileService struct{}

// store files' information into database
func (s *FileService) StoreFile(fileSize, storageId uint, fileName string, folderId *uint) error {
	db := global.MysqlDB
	var tmp uint 
	// root directory translate
	if folderId == nil {
		tmp = 0
	} else {
		tmp = *folderId
	}
	res := db.Create(&model.File{
		FileSize: fileSize,
		FolderId: tmp,
		StorageId: storageId,
		FileName: fileName,
	})
	return res.Error
}

// update file related information
func (s *FileService) UpdateFileName(fileId uint, newName string) error {
	db := global.MysqlDB
	res := db.Model(&model.File{}).Where("id = ?", fileId).Update("FileName", newName)
	return res.Error
}

// delete file record
func (s *FileService) DeleteFile(fileId uint) error {
	db := global.MysqlDB
	res := db.Where("id = ?", fileId).Delete(&model.File{})
	return res.Error
}
