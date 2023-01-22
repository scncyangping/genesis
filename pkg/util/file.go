package util

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"strings"
)

// MatchOptions 传入文件名称 返回是否需要match
type FileMatchOptions func(string) bool

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

// un zip file struct
type FileUnZip struct {
	path         string
	files        []*FileInfo
	matchOptions []FileMatchOptions
}

func NewFileUnzip(path string) *FileUnZip {
	return &FileUnZip{
		path:         path,
		files:        make([]*FileInfo, 0),
		matchOptions: make([]FileMatchOptions, 0),
	}
}

func (t *FileUnZip) UnZip(path string) error {
	if zipReader, err := zip.OpenReader(path); err != nil {
		return err
	} else {
		defer zipReader.Close()

		// m := make(map[string]string)
		for _, file := range zipReader.Reader.File {
			t.readFile(file)
		}
	}

	return nil
}

func (t *FileUnZip) readFile(file *zip.File) error {
	// 打包文件中的文件就像普通的一个文件对象一样
	zippedFile, err := file.Open()
	if err != nil {
		return err
	}
	defer zippedFile.Close()

	if !file.FileInfo().IsDir() {
		flag := false
		for _, n := range t.matchOptions {
			if n(file.Name) {
				flag = true
				break
			}
		}
		if flag {
			content, err := io.ReadAll(zippedFile)
			if err != nil {
				log.Println(err)
			}
			t.files = append(t.files, &FileInfo{
				path: file.Name,
				data: content,
			})
		}
	}

	return nil
}

// file match options: suffix matches
type MatchSuffix struct {
	suffix []string
}

func NewMatchSuffix(s []string) *MatchSuffix {
	return &MatchSuffix{suffix: s}
}

func (m *MatchSuffix) SuffixMatch(url string) bool {
	for _, n := range m.suffix {
		if strings.HasSuffix(url, n) {
			return true
		}
	}
	return false
}
