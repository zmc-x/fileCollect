package system

import (
	"os"
	"path/filepath"
)

type FolderService struct{}

// create the folder in the storage
func (fs *FolderService) CreateFolder(storagePath, folderSrc string) error {
	return os.Mkdir(filepath.Join(storagePath, folderSrc), 0644)
}

// update the folder Name
func (fs *FolderService) UpdateFolderName(storagePath, oldFolderName, newFolderName string) error {
	return os.Rename(filepath.Join(storagePath, oldFolderName), filepath.Join(storagePath, newFolderName))
}


// delete the folder information 
func (fs *FolderService) DeleteFolder(storagePath, folderSrc string) error {
	return os.RemoveAll(filepath.Join(storagePath, folderSrc))
}
