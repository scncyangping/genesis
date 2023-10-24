package dandelion

import "genesis/pkg/util"

type FileTemp struct {
	path string
	file *util.FileUnZip
}

func NewFileTemp(path string) *FileTemp {
	return &FileTemp{
		path: path,
		file: util.NewFileUnzip(),
	}
}

// // 添加模版文件过滤器
// Template key 文件路径
// value 模版内容

func (f *FileTemp) AddOption(op MatchOptions) {
	f.file.AddOption(util.FileMatch{
		MType: util.FileMatchPath,
		DoFn:  op,
	})
}

func (f *FileTemp) Template(string) ([]*util.FileInfo, error) {
	return f.file.UnZip(f.path)
}
