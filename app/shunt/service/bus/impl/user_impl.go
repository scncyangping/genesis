// @Author: YangPing
// @Create: 2023/10/23
// @Description: User业务逻辑

package impl

import (
	"errors"
	"genesis/app/common/base"
	"genesis/app/shunt/adapter/cqe/cmd"
	"genesis/app/shunt/adapter/cqe/query"
	"genesis/app/shunt/adapter/vo"
	"genesis/app/shunt/repository"
	"genesis/pkg/types"
	"genesis/pkg/util"
	"genesis/pkg/util/snowflake"
	"github.com/samber/lo"
)

type UserServiceImpl struct {
	repo repository.UserRepositoryI
}

// NewUserServiceImpl
// NEED WIRE
func NewUserServiceImpl(repo repository.UserRepositoryI) *UserServiceImpl {
	return &UserServiceImpl{repo: repo}
}

func (u *UserServiceImpl) Add(cmd *cmd.UserSaveCmd) (string, error) {
	var (
		en base.UserGorm
	)
	if err := util.Copy(&en, cmd); err != nil {
		return "", err
	}
	en.Id = snowflake.NextId()
	if err := u.repo.Save(&en); err != nil {
		return "", err
	}
	return en.Id, nil
}

func (u *UserServiceImpl) Update(cmd *cmd.UserUpdateCmd) error {
	var (
		en base.UserGorm
	)
	// step1. get en by uniq id
	if _, err := u.repo.GetById(cmd.Id); err != nil {
		if err != types.ErrRecordNotFound {
			return err
		} else {
			return errors.New("对应数据不存在")
		}
	}
	// step2. copy to en
	if err := util.Copy(&en, cmd); err != nil {
		return err
	}
	// step3. update
	if err := u.repo.Update(&en); err != nil {
		return err
	}
	return nil
}

func (u *UserServiceImpl) Query(query *query.UserListQuery) (*vo.PageResult, error) {
	var (
		qb = query.Q()
		qm = lo.MapKeys(util.StructToMap(query), func(value any, key string) string {
			return util.UnCamelize(key)
		})
		vos = make([]vo.UserVO, 0)
	)
	if ens, count, err := u.repo.QueryByPage(qm, qb.Skip, qb.Limit, qb.Sort, qb.SortBy, true); err != nil {
		return nil, err
	} else {
		if err := util.Copy(&vos, ens); err != nil {
			return nil, err
		}
		return vo.NewPageResult(vos, query.Page, query.Size, count), nil
	}
}

func (u *UserServiceImpl) GetById(id string) (*vo.UserVO, error) {
	var (
		vo vo.UserVO
	)
	en, err := u.repo.GetById(id)
	if err != nil {
		if err != types.ErrRecordNotFound {
			return nil, err
		} else {
			return nil, errors.New("对应数据不存在")
		}
	}
	if err := util.Copy(&vo, en); err != nil {
		return nil, err
	}
	return &vo, nil
}

func (u *UserServiceImpl) DeleteById(id string) error {
	if err, su := u.repo.DeleteById(id); err != nil {
		return err
	} else if su == false {
		return errors.New("删除失败,请检查对应数据是否存在")
	}
	return nil
}

func (u *UserServiceImpl) DeleteByMap(m map[string]any) error {
	if err, _ := u.repo.DeleteBy(m, true); err != nil {
		return err
	}
	return nil
}
