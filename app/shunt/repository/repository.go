// @Author: YangPing
// @Create: 2023/10/23
// @Description: 全局Repository定义

package repository

import (
	"genesis/app/common/base"
	"genesis/app/shunt/adapter/cqe/query"
)

type UniversalRepositoryI[T base.EnI] interface {
	// GetById 根据主键查询数据
	GetById(string) (*T, error)
	// FindBy 根据条件查询数据
	// bool: map中是否去除空元素
	FindBy(map[string]any, bool) ([]*T, error)
	// QueryByPage
	// any: struct or map
	// bool: map中是否去除空元素
	QueryByPage(any, *query.PageQBase, bool) ([]*T, int64, error)
	// SaveBatch 批量保存
	// int: 每一批次数量
	SaveBatch([]*T, int) error
	// Save 保存/批量保存
	Save(*T) error
	// Update 实体更新
	Update(*T) error
	// UpdateByMap 根据主键ID批量更新
	UpdateByMap(map[string]any, []string) (error, int64)
	// DeleteById 根据主键删除数据
	/**/
	DeleteById(string) (error, bool)
	// DeleteBy 根据条件删除数据
	// bool: map中是否去除空元素
	DeleteBy(map[string]any, bool) (error, int64)
}

type UserRepositoryI interface {
	UniversalRepositoryI[base.UserGorm]
}

type User2RepositoryI interface {
	UniversalRepositoryI[base.UserGorm]
}
