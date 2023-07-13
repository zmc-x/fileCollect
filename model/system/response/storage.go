package response

import "gorm.io/gorm"

const (
	File = false 
	Folder = true
)

// file and folder information in the storage
type StorageFileList struct {
	gorm.Model
	FName	string 
	FSize	uint
	FType	bool
}

// storage info
type StorageInfo struct {
	Id				uint 
	StorageName 	string 
	StorageUrlName 	string 
	StorageRealPath string
	// storage's path
	Path			string 
	Status	 		bool
}