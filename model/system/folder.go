package system

import "gorm.io/gorm"

type Folder struct {
	gorm.Model
	Files	[]File
	FolderName string `gorm:"uniqueIndex:folderInfo"`
	StorageId  uint `gorm:"uniqueIndex:folderInfo"`
	ParentFolderId uint `gorm:"uniqueIndex:folderInfo"`
}