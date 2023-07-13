package request

// file info
type FileInfo struct {
	FolderID        uint   `json:"folderId"`       
	Path            string `json:"path"`           
	StorageID       uint   `json:"storageId"`      
	StorageRealPath string `json:"storageRealPath"`
}

// files array
type FileId struct {
	Files []uint			`json:"files"`
}

// update filename structure
type UpdateNameReq struct {
	FileID      uint   `json:"fileId"`     
	NewFileName string `json:"newFileName"`
}