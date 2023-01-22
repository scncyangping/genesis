package util

import (
	"archive/zip"
	"log"
	"os"
)

type FileInfo struct {
	path string // file path
	data []byte // file data
}

// zip file struct
type FileZip struct {
	savePath string
	files    []*FileInfo
}

func NewZip(path string, files []*FileInfo) *FileZip {
	return &FileZip{
		savePath: path,
		files:    files,
	}
}

func NewFileZip(path string) *FileZip {
	return &FileZip{
		savePath: path,
		files:    make([]*FileInfo, 0),
	}
}

func (c *FileZip) Append(path string, data []byte) *FileZip {
	c.files = append(c.files, &FileInfo{
		path: path,
		data: data,
	})
	return c
}

func (c *FileZip) Zip() error {
	// 创建一个打包文件
	outFile, err := os.Create(c.savePath)
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()
	// 创建zip writer
	zipWriter := zip.NewWriter(outFile)
	// 往打包文件中写文件
	for _, file := range c.files {
		if fileWriter, err := zipWriter.Create(file.path); err != nil {
			return err
		} else {
			if _, err := fileWriter.Write(file.data); err != nil {
				return err
			}

		}
	}
	// 清理
	return zipWriter.Close()
}
