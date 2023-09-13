package request

type GenericInfo struct {
	Path       string `json:"path"`
	StorageKey string `json:"storageKey"`
}

type Folder struct {
	FolderName string `json:"folderName"`
}

// create folder info
type CreateFolderInfo struct {
	GenericInfo
	FolderName string `json:"folderName"`
}

// delete folders info
type DeleteFolderInfo struct {
	GenericInfo
	Folders    []Folder `json:"folders"`   
}

// update folder info
type UpdateFolderInfo struct {
	GenericInfo
	Folder
	FolderNewName string `json:"folderNewName"`
}