package response

import "time"

const (
	// file = 0, folder = 1
	File = false 
	Folder = true
)

// file and folder information in the storage
type StorageFileList struct {
	UpdateAt time.Time
	FName	 string 
	FSize	 uint
	FType	 bool
}

// storage info
type StorageInfo struct {
	StorageName 	string 
	StorageKey	 	string 
	Status	 		bool
}