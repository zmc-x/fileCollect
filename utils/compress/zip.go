package compress

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

type Compress struct{
	zipWri *zip.Writer
}

func NewCompress(w io.Writer) *Compress {
	return &Compress{
		zipWri: zip.NewWriter(w),
	}
}

// close zipWri
func (c *Compress) Close() {
	c.zipWri.Close()
}

// create zip archive
func (c *Compress) CreateZip(path []string, prefix string) (err error) {
	for _, v := range path {
		zipName, err := filepath.Rel(prefix, v)
		if err != nil {
			return err
		}
		dstf, err := c.zipWri.Create(zipName)
		if err != nil {
			return err
		}
		srcf, err := os.Open(v)
		if err != nil {
			return err
		}
		defer srcf.Close()
		_, err = io.Copy(dstf, srcf)
		if err != nil {
			return err
		}
	}
	return nil
}