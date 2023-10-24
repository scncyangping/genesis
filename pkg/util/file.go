// @Author: YangPing
// @Create: 2023/10/23
// @Description: 文件操作工具类

package util

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"strings"

	"github.com/samber/lo"
)

// FileMatchEnum MatchOptions 传入文件名称 返回是否需要match
// type FileMatchOptions func(string) bool
type FileMatchEnum string

const (
	FileMatchPath FileMatchEnum = "PATH"
	FileMatchBody FileMatchEnum = "BODY"
)

type FileMatch struct {
	MType FileMatchEnum
	DoFn  func(string) bool
}

type FileInfo struct {
	Path string // file path
	Data []byte // file data
}

// FileZip zip file struct
type FileZip struct {
	SavePath string
	Files    []*FileInfo
}

func NewZip(path string, files []*FileInfo) *FileZip {
	return &FileZip{
		SavePath: path,
		Files:    files,
	}
}

func NewFileZip(path string) *FileZip {
	return &FileZip{
		SavePath: path,
		Files:    make([]*FileInfo, 0),
	}
}

func (c *FileZip) Append(path string, data []byte) *FileZip {
	c.Files = append(c.Files, &FileInfo{
		Path: path,
		Data: data,
	})
	return c
}

func (c *FileZip) Zip() error {
	// 创建一个打包文件
	outFile, err := os.Create(c.SavePath)
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()
	// 创建zip writer
	zipWriter := zip.NewWriter(outFile)
	// 往打包文件中写文件
	for _, file := range c.Files {
		if fileWriter, err := zipWriter.Create(file.Path); err != nil {
			return err
		} else {
			if _, err := fileWriter.Write(file.Data); err != nil {
				return err
			}

		}
	}
	// 清理
	return zipWriter.Close()
}

// FileUnZip un zip file struct
type FileUnZip struct {
	files        []*FileInfo
	matchOptions []FileMatch
}

func NewFileUnzip() *FileUnZip {
	return &FileUnZip{
		files:        make([]*FileInfo, 0),
		matchOptions: make([]FileMatch, 0),
	}
}

func (t *FileUnZip) AddOption(m FileMatch) {
	if t.matchOptions == nil {
		t.matchOptions = make([]FileMatch, 0)
	}
	t.matchOptions = append(t.matchOptions, m)
}

func (t *FileUnZip) UnZip(path string) ([]*FileInfo, error) {
	if zipReader, err := zip.OpenReader(path); err != nil {
		return nil, err
	} else {
		defer zipReader.Close()

		// m := make(map[string]string)
		for _, file := range zipReader.Reader.File {
			t.readFile(file)
		}
	}

	return t.files, nil
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

		em := lo.GroupBy(t.matchOptions, func(item FileMatch) FileMatchEnum {
			return item.MType
		})

		if fn, ok := em[FileMatchPath]; ok {
			for _, v := range fn {
				rs := v.DoFn(file.Name)
				if rs {
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
					Path: file.Name,
					Data: content,
				})
			}
		}
	}

	return nil
}

// MatchSuffix file match options: suffix matches
type MatchSuffix struct {
	suffix []string
}

var defaultMatchSuffix = []string{"go", "mod", "tmpl", "xml"}

func NewMatchSuffix() *MatchSuffix {
	return &MatchSuffix{suffix: defaultMatchSuffix}
}

func (m *MatchSuffix) BuildMatchSuffix(s []string) *MatchSuffix {
	m.suffix = s
	return m
}

func (m *MatchSuffix) SuffixMatch(url string) bool {
	for _, n := range m.suffix {
		if strings.HasSuffix(url, n) {
			return true
		}
	}
	return false
}
