package response

import "time"

// data
type StorageInfo struct {
	StorageList []StorageList `json:"storageList"`
}

// storageList
type StorageList struct {
	Path           string `json:"path"`          
	Status         bool   `json:"status"`        
	StorageID      uint   `json:"storageId"`     
	StorageName    string `json:"storageName"`   
	StorageURLName string `json:"storageUrlName"`
}

// fileInfo
type FilesInfo struct {
	FileList []FileList `json:"fileList"`
}

type FileList struct {
	FileID   uint   	`json:"fileId"`  
	FName    string 	`json:"fName"`   
	FSize    uint	  	`json:"fSize"`   
	FType    bool   	`json:"fType"`   
	UpdateAt time.Time 	`json:"updateAt"`
}