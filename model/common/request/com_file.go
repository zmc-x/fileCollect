package request

// file info
type FileInfo struct {
	StorageKey string `form:"storageKey"`
	Path	   string `form:"path"`
}

// files array
type FileArray struct {
	Files          []File `json:"files"`
	Path           string `json:"path"`
	StorageKey     string `json:"storageKey"`
}

// singal file information
type File struct {
	FileName string `json:"fileName"`
}

// update filename structure
type UpdateNameReq struct {
	File
	NewFileName    string `json:"newFileName"`
	Path           string `json:"path"`
	StorageKey     string `json:"storageKey"`
}
