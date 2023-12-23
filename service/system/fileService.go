package system

import (
	"fileCollect/utils/compress"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type FileService struct{}

// store files locally
func (s *FileService) StoreFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	if err = os.MkdirAll(filepath.Dir(dst), 0750); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

// update file related information
func (s *FileService) UpdateFileName(storagePath, oldFileSrc, newFileSrc string) error {
	return os.Rename(filepath.Join(storagePath, oldFileSrc), filepath.Join(storagePath, newFileSrc))
}

// delete file record
func (s *FileService) DeleteFile(storagePath, fileSrc string) error {
	return os.Remove(filepath.Join(storagePath, fileSrc))
}

// Download the package
// zip
func (s *FileService) DownloadCompressFile(files, folders []string, fileSrc string) (string, *os.File, *compress.Compress, error) {
	// translate the zip file
	path := []string{}
	var findFile func(src string)
	findFile = func(src string) {
		dir, _ := os.Stat(src)
		if !dir.IsDir() {
			path = append(path, src)
			return
		}
		files, _ := os.ReadDir(src)
		for _, v := range files {
			findFile(filepath.Join(src, v.Name()))
		}
	}
	zipname := strconv.Itoa(int(time.Now().Unix())) + ".zip"
	zipPath := filepath.Join(fileSrc, zipname)
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return "", nil, nil, err
	}
	comp := compress.NewCompress(zipFile)
	for _, folder := range folders {
		findFile(filepath.Join(fileSrc, folder))
		err = comp.CreateZip(path, fileSrc)
		if err != nil {
			return "", nil, nil, err
		}
		path = nil
	}
	for _, file := range files {
		path = append(path, filepath.Join(fileSrc, file))
	}
	if err = comp.CreateZip(path, fileSrc); err != nil {
		return "", nil, nil, err
	}
	return zipPath, zipFile, comp, nil
}

// download file
func (s *FileService) Download(files []string, fileSrc string) (string, error) {
	// single file
	filename := files[0]
	src := filepath.Join(fileSrc, filename)
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return src, err
	}
	return src, nil
}
