package response

const (
	// file = 0, folder = 1
	File = false 
	Folder = true
)

// data
type Storages struct {
	StorageList []StorageItem `json:"storageList,omitempty"`
}

// storageList
type StorageItem struct {
	Path           string `json:"path"`          
	Status         bool   `json:"status"`        
	StorageKey     string `json:"storageKey"`     
	StorageName    string `json:"storageName"`
}

// fileInfo
type Files struct {
	FileList []FileItem `json:"fileList,omitempty"`
}

type FileItem struct {
	FName    string 	`json:"fName"`   
	FSize    uint	  	`json:"fSize"`   
	FType    bool   	`json:"fType"`   
	UpdateAt string  	`json:"updateAt"`
}