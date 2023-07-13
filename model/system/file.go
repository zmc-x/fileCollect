package system

import "gorm.io/gorm"

type File struct {
	gorm.Model
	StorageId uint `gorm:"uniqueIndex:fileInfo"`
	FileSize  uint
	FileName  string `gorm:"uniqueIndex:fileInfo;not null"`
	FolderId  uint	`gorm:"uniqueIndex:fileInfo"`
}