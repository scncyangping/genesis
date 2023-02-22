package main

import (
	"fmt"
	"genesis/pkg/util"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var Source = ""
var Target = ""
var ZNUM = 0
var SizeNum = 0

func main() {
	// 获取文件目录
	args := os.Args

	Source = args[1] // /abc/csdf/
	Target = args[2] // /abc/csdf/
	s := args[3]     // 1 copy 2 save

	if !strings.HasSuffix(Source, string(filepath.Separator)) {
		Source += string(filepath.Separator)
	}
	if !strings.HasSuffix(Target, string(filepath.Separator)) {
		Target += string(filepath.Separator)
	}

	num, err := strconv.Atoi(s)
	if err != nil {
		fmt.Println(err)
		return
	}
	if num == 1 {
		ReadFile(Source)
	} else if num == 2 {
		ToFile(Target)
	}
}

func ReadFile(dirPath string) {
	fs, err := ioutil.ReadDir(dirPath)
	if err != nil {
		fmt.Println(fmt.Sprintf("error : %v", err))
	}
	for _, file := range fs {
		if file.IsDir() {
			ReadFile(filepath.Join(dirPath, file.Name()))
		} else {
			// file.Size()
			size := file.Size()
			if size < 81920 {
				continue
			}
			SizeNum++

			if SizeNum == 500 {
				ZNUM++
				SizeNum = 0
			}
			// 读取文件,从新命名为 0_2023_02_03_x12312313_0092039203.jpeg
			endStr := strings.Replace(filepath.Join(dirPath, file.Name()), Source, "", -1)
			newImageName := strings.Replace(endStr, string(filepath.Separator), "_", -1)

			path := filepath.Join(Target, strconv.Itoa(ZNUM))

			util.CreateDirectory(path)

			err := os.Rename(
				filepath.Join(dirPath, file.Name()),
				filepath.Join(path, newImageName),
			)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func ToFile(tar string) {
	fs, err := ioutil.ReadDir(tar)
	if err != nil {
		fmt.Println(fmt.Sprintf("error : %v", err))
	}
	for _, file := range fs {
		if file.IsDir() {
			ToFile(filepath.Join(tar, file.Name()))
		} else {
			fileExt := filepath.Ext(file.Name())

			//endStr := strings.Replace(filepath.Join(tar, file.Name()), Target, "", -1)
			endStr := file.Name()

			// 第一个目录为分类目录,去要去掉
			es := strings.ReplaceAll(endStr, "_", string(filepath.Separator))
			es = strings.ReplaceAll(es, fileExt, ".jpg")

			err := os.Rename(
				filepath.Join(tar, file.Name()),
				filepath.Join(Source, es),
			)
			if err != nil {
				fmt.Println(err)
			}
		}

	}
}
