package system

import "gorm.io/gorm"

type Storage struct {
	gorm.Model
	Files           []File
	Folders         []Folder
	StorageName     string `gorm:"unique;not null"`
	StorageRealPath string 
	StorageUrlName  string `gorm:"unique;not null"`
	Status          bool
}
