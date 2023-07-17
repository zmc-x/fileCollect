package response

import "time"

const (
	File = false 
	Folder = true
)

// file and folder information in the storage
type StorageFileList struct {
	ID		 uint
	UpdateAt time.Time
	FName	 string 
	FSize	 uint
	FType	 bool
}

// storage info
type StorageInfo struct {
	Id				uint 
	StorageName 	string 
	StorageUrlName 	string 
	// storage's path
	Path			string 
	Status	 		bool
}