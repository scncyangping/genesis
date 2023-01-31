package mixture

import (
	"bytes"
	"errors"
	"genesis/pkg/config/app/shunt"
	"genesis/pkg/core/dandelion"
	"genesis/pkg/core/dandelion/padding"
	"genesis/pkg/util"
	"path/filepath"
	"strings"
	tmp "text/template"
)

type TemplateSrv interface {
	// 添加模版文件过滤器
	AddOption(dandelion.MatchOptions)
	// Template key 文件路径
	// value 模版内容
	Template(string) ([]*util.FileInfo, error)
}

type PaddingSrv interface {
	// AddPadding Padding key 填充key
	// value 填充值
	AddPadding(dandelion.PaddingOptions)
	Padding() map[string]any
}

type TargetSrv interface {
	Zip() error
}

type Manager struct {
	// 模版数据
	templateSrv TemplateSrv
	// 填充数据
	paddingSrv PaddingSrv
	// 全局参数
	managerStartOp *ManagerStartOp
}

type ManagerStartOp struct {
	ModuleName   string
	DBName       string
	TemplatePath string
	SavePath     string
	MatchSuffix  []string
	Tables       []string
}

func NewManager(m *ManagerStartOp) *Manager {
	// 构建模版读取数据
	tmp := dandelion.NewFileTemp(m.TemplatePath)
	// 设置模块文件名称过滤规则
	// 仅获取以xxx后缀结尾的数据
	sf := util.NewMatchSuffix()

	if m.MatchSuffix != nil && len(m.MatchSuffix) > 0 {
		sf.BuildMatchSuffix(m.MatchSuffix)
	}
	tmp.AddOption(sf.SuffixMatch)

	// 设置填充内容
	pad := padding.NewPadding()
	// 读取数据库中对应表信息
	tablePadding := padding.NewTablePadding(
		m.ModuleName,
		m.Tables,
		padding.NewTableOp(m.DBName).BuildGormDB(shunt.GormDB()))

	pad.AddPadding(tablePadding.Add)

	return &Manager{
		managerStartOp: m,
		templateSrv:    tmp,
		paddingSrv:     pad,
	}
}

func (m *Manager) BuildTmplate(t TemplateSrv) *Manager {
	m.templateSrv = t
	return m
}

func (m *Manager) BuildPadding(t PaddingSrv) *Manager {
	m.paddingSrv = t
	return m
}

func (m *Manager) BuildStartOp(t *ManagerStartOp) *Manager {
	m.managerStartOp = t
	return m
}

func (m *Manager) Start() error {
	// 解析模版文件
	temFiles, err := m.templateSrv.Template(m.managerStartOp.TemplatePath)
	if err != nil {
		return err
	}
	// 生成填充数据
	sourceKeyValueMap := m.paddingSrv.Padding()

	sm := make(map[string]string, 0)
	for _, file := range temFiles {
		sm[file.Path] = string(file.Data)
	}

	if z, err := m.mixture(sm, sourceKeyValueMap); err != nil {
		return err
	} else {
		return z.Zip()
	}
}

func (m *Manager) mixture(tm map[string]string, data map[string]any) (*util.FileZip, error) {
	c := util.NewFileZip(m.managerStartOp.SavePath)

	for k, v := range tm {
		// k 文件路径
		// v 匹配前数据
		f := util.FileInfo{
			Path: k,
		}
		if !strings.HasSuffix(k, "tmpl") {
			if tem, err := tmp.New(k).Parse(v); err != nil {
				return nil, err
			} else {
				var doc bytes.Buffer
				if err = tem.Execute(&doc, &data); err != nil {
					return nil, err
				}
				f.Data = doc.Bytes()
				c.Files = append(c.Files, &f)
			}
		} else {
			// tmpl 结尾 需要根据表信息生成文件
			modelList, ok := data["Models"]
			if !ok {
				return nil, errors.New("models not found")
			}
			models := modelList.([]interface{})
			for _, model := range models {
				if tem, err := tmp.New(k).Parse(v); err != nil {
					return nil, err
				} else {
					var doc bytes.Buffer

					if err = tem.Execute(&doc, &model); err != nil {
						return nil, err
					}

					split, fileName := filepath.Split(k)
					// replace file suffix
					fileName = strings.ReplaceAll(fileName, "tmpl", "go")
					// replace moudle name
					OldName := model.(map[string]any)["OldName"].(string)
					fileName = strings.ReplaceAll(fileName, "xxx", OldName)

					ff := util.FileInfo{
						Path: filepath.Join(split, fileName),
						Data: doc.Bytes(),
					}

					c.Files = append(c.Files, &ff)
				}

			}
		}
	}
	return c, nil
}
