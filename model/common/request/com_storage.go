package request

// storageInfo
type StorageInfo struct {
	StorageName     string `json:"storageName"`    
	StorageRealPath string `json:"storageRealPath"`
	StorageURLName  string `json:"storageUrlName"` 
	DeadLine		string `json:"deadLine"`
}

// update model
type UpdateGeneric struct {
	NewName    string  `json:"newName"`  
	NewUrlName string  `json:"newUrlName"`
	NewPath	   string  `json:"newPath"`
	NewStatus  bool    `json:"newStatus"`
	DeadLine   string  `json:"deadLine,omitempty"`
	StorageKey string  `json:"storageKey"`
}

// request storage file list
type ReqStorageList struct {
	Path       string `json:"path"`      
	StorageKey string `json:"storageKey"`
}