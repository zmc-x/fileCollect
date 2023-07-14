package request

// file info
type FileInfo struct {
	FolderID        uint   `form:"folderId"`       
	Path            string `form:"path"`           
	StorageID       uint   `form:"storageId"`      
}

// files array
type FileArray struct {
	Files []File		`json:"files"`
	Path  string		`json:"path"`
	StorageId uint	 	`json:"storageId"`
}

// singal file information
type File struct {
	FileID    uint   `json:"fileId"`  
	FileName  string `json:"fileName"`
}

// update filename structure
type UpdateNameReq struct {
	File     
	NewFileName string  `json:"newFileName"`
	Path 		string  `json:"path"`
	StorageId 	uint	`json:"storageId"`
}