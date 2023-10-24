// @Author: YangPing
// @Create: 2023/10/23
// @Description: Gorm基础操作

package base

import (
	"fmt"
	"genesis/pkg/plugin/mysql"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UniversalRepositoryI[T EnI] interface {
	// GetById 根据主键查询数据
	GetById(string) (*T, error)
	// FindBy 根据条件查询数据
	// bool: map中是否去除空元素
	FindBy(map[string]any, bool) ([]*T, error)
	// QueryByPage
	// any: struct or map
	// bool: map中是否去除空元素
	QueryByPage(any, int, int, string, string, bool) ([]*T, int64, error)
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

type UniversalGormRepo[T EnI] struct {
	gm    *gorm.DB
	model T
}

func NewUniversalGormRepo[T EnI](model T, db *gorm.DB) *UniversalGormRepo[T] {
	return &UniversalGormRepo[T]{model: model, gm: db}
}

func (u *UniversalGormRepo[T]) Gm() *gorm.DB {
	return u.gm
}

func (u *UniversalGormRepo[T]) Save(data *T) error {
	if err := u.gm.Model(&u.model).Create(data).Error; err != nil {
		return errors.New(fmt.Sprintf("[%s]新增数据异常:%v", u.model.TableName(), err))
	}
	return nil
}

func (u *UniversalGormRepo[T]) Update(data *T) error {
	if err := u.gm.Updates(data).Error; err != nil {
		return errors.New(fmt.Sprintf("[%s]更新实体数据异常:%v", u.model.TableName(), err))
	}
	return nil
}

func (u *UniversalGormRepo[T]) UpdateByMap(m map[string]any, ids []string) (error, int64) {
	var (
		err error
	)
	tx := u.gm.Model(&u.model).Where("id IN ?", ids).Updates(m)
	if err := tx.Error; err != nil {
		err = errors.New(fmt.Sprintf("[%s]根据主键ID更新数据异常:%v", u.model.TableName(), err))
	}
	return err, tx.RowsAffected
}

func (u *UniversalGormRepo[T]) DeleteById(id string) (error, bool) {
	tx := u.gm.Model(&u.model).Delete("id=?", id)
	if err := tx.Error; err != nil {
		return errors.New(fmt.Sprintf("[%s]根据主键ID删除数据异常:%v", u.model.TableName(), err)), false
	}
	return nil, tx.RowsAffected == 1
}

func (u *UniversalGormRepo[T]) DeleteBy(m map[string]any, removeZero bool) (error, int64) {
	var (
		err   error
		count int64
	)
	if gormDb, err := mysql.BuildGormQuery(m, removeZero, u.gm.Model(&u.model)); err != nil {
		err = errors.New(fmt.Sprintf("[%s]根据指定条件组装删除条件异常:%v", u.model.TableName(), err))
	} else {
		tx := gormDb.Table(u.model.TableName()).Delete(nil)
		if err := tx.Error; err != nil {
			err = errors.New(fmt.Sprintf("[%s]根据指定条件删除数据异常:%v", u.model.TableName(), err))
		} else {
			count = tx.RowsAffected
		}
	}
	return err, count
}

func (u *UniversalGormRepo[T]) SaveBatch(data []*T, batchSize int) error {
	if err := u.gm.Model(&u.model).CreateInBatches(data, batchSize).Error; err != nil {
		return errors.New(fmt.Sprintf("[%s]批量新增数据异常:%v", u.model.TableName(), err))
	}
	return nil
}

func (u *UniversalGormRepo[T]) QueryByPage(data any, skip, limit int, sort, sortBy string, removeZero bool) ([]*T, int64, error) {
	var (
		ens   []*T
		count int64
	)

	if gormDb, err := mysql.BuildGormQuery(data, removeZero, u.gm.Model(&u.model)); err != nil {
		return ens, count, errors.New(fmt.Sprintf("[%s]构建查询条件异常:%v", u.model.TableName(), err))
	} else {
		if err := gormDb.Count(&count).
			Limit(limit).
			Offset(skip).
			Order(clause.OrderByColumn{
				Column: clause.Column{Name: sort},
				Desc:   sortBy != "ASC",
			}).Find(&ens).Error; err != nil {
			return nil, count, errors.New(fmt.Sprintf("[%s]根据指定条件查询异常:%v", u.model.TableName(), err))
		}
	}
	return ens, count, nil
}

func (u *UniversalGormRepo[T]) FindBy(m map[string]any, removeZero bool) ([]*T, error) {
	var (
		ens []*T
	)
	if gormQuery, err := mysql.BuildGormQuery(m, removeZero, u.gm.Model(&u.model)); err != nil {
		return ens, errors.New(fmt.Sprintf("[%s]构建查询条件异常:%v", u.model.TableName(), err))
	} else {
		if err := gormQuery.Find(&ens).Error; err != nil {
			return nil, errors.New(fmt.Sprintf("[%s]根据指定条件查询异常:%v", u.model.TableName(), err))
		}
	}
	return ens, nil
}

func (u *UniversalGormRepo[T]) GetById(id string) (*T, error) {
	var (
		en *T
	)
	if err := u.gm.Where("id", id).Take(&en).Error; err != nil {
		return nil, err
	}
	return en, nil
}
