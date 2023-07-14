package request

type GenericFolder struct {
	FolderId	uint 	`json:"folderId"`
	FolderName  string	`json:"folderName"`
}

type GenericStorage struct {
	StorageId	uint	`json:"stroageID"`
	Path		string	`json:"path"`
}

// create folder request info
type CreateFolderInfo struct {
	GenericStorage
	FolderName     string `json:"folderName"`    
	ParentFolderID uint   `json:"parentFolderId"`
}

// delete folder request info
type DeleteFolderInfo struct {
	Folders   []GenericFolder `json:"folders"`  
	Path      string   		  `json:"path"`     
	StorageID uint    		  `json:"storageId"`
}

// update folder request info
type UpdateFolderInfo struct {
	GenericFolder
	GenericStorage
	FolderNewName string 	`json:"folderNewName"`
}
