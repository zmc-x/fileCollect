package request

// storageInfo
type StorageInfo struct {
	StorageName     string `json:"storageName"`    
	StorageRealPath string `json:"storageRealPath"`
	StorageURLName  string `json:"storageUrlName"` 
}

// update model
type UpdateGeneric struct {
	NewName    string  `json:"newName"`  
	NewUrlName string  `json:"newUrlName"`
	NewPath	   string  `json:"newPath"`
	NewStatus  bool    `json:"newStatus"`
	StorageID  uint    `json:"storageId"`
}