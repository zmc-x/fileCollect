package response

// data
type StorageInfo struct {
	StorageList []StorageList `json:"storageList"`
}

// storageList
type StorageList struct {
	Path           string `json:"path"`          
	Status         bool   `json:"status"`        
	StorageKey     string `json:"storageKey"`     
	StorageName    string `json:"storageName"`
}

// fileInfo
type FilesInfo struct {
	FileList []FileList `json:"fileList"`
}

type FileList struct {
	FName    string 	`json:"fName"`   
	FSize    uint	  	`json:"fSize"`   
	FType    bool   	`json:"fType"`   
	UpdateAt string  	`json:"updateAt"`
}
