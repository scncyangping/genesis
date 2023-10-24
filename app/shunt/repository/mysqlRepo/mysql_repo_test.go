// @Author: YangPing
// @Create: 2023/10/23
// @Description: Gorm基础操作测试

package mysqlRepo

import (
	"fmt"
	"genesis/app/common/base"
	"genesis/app/common/config"
	"genesis/app/shunt/adapter/cqe/query"
	"genesis/pkg/plugin/mysql"
	"gorm.io/gorm"
	"strconv"
	"testing"
)

func TestUserMysqlRepo(t *testing.T) {
	repo := NewConn(t)

	t.Run("Save", func(t *testing.T) {
		var user base.UserGorm
		user.Id = strconv.Itoa(1)
		user.Name = "张三"
		user.Age = 12
		user.Status = 1
		if err := repo.Save(&user); err != nil {
			t.Fatalf("保存用户失败: %v", err)
		}
	})

	t.Run("SaveBatch", func(t *testing.T) {
		var users []*base.UserGorm
		for i := 2; i < 6; i++ {
			var user base.UserGorm
			user.Id = strconv.Itoa(i)
			user.Name = fmt.Sprintf("张三%d", i)
			user.Age = 12 + i
			if i%2 == 0 {
				user.Status = 1
			} else {
				user.Status = 2
			}
			users = append(users, &user)
		}
		if err := repo.SaveBatch(users, 1); err != nil {
			t.Fatalf("保存用户失败: %v", err)
		}
	})

	t.Run("Update", func(t *testing.T) {
		up1 := base.UserGorm{}
		up1.Id = strconv.Itoa(1)
		up1.Name = "张三_up1"
		err := repo.Update(&up1)
		if err != nil {
			t.Fatalf("实体更新用户失败: %v", err)
		}

		if qu, err := repo.GetById(strconv.Itoa(1)); err != nil {
			t.Fatalf("实体更新->查询用户失败: %v", err)
		} else if qu.Name != "张三_up1" {
			t.Fatalf("实体更新->更新用户失败,查询值不相等")
		}
	})

	t.Run("UpdateByMap", func(t *testing.T) {
		um := map[string]any{
			"name": "张三_up2",
		}
		err, count := repo.UpdateByMap(um, []string{"1"})
		t.Logf("受影响行数:%d", count)
		if err != nil {
			t.Fatalf("Map更新用户失败: %v", err)
		}
		if qu, err := repo.GetById(strconv.Itoa(1)); err != nil {
			t.Fatalf("实体更新->查询用户失败: %v", err)
		} else if qu.Name != "张三_up2" {
			t.Fatalf("实体更新->更新用户失败,查询值不相等")
		}
	})

	t.Run("DeleteById", func(t *testing.T) {
		err, su := repo.DeleteById("1")
		t.Logf("删除结果:%v", su)
		if err != nil {
			t.Fatalf("根据ID删除用户失败: %v", err)
		}
		if qu, err := repo.GetById(strconv.Itoa(1)); err != nil && err != gorm.ErrRecordNotFound {
			t.Fatalf("实体删除->查询用户失败: %v", err)
		} else if qu != nil && qu.Id != "" {
			t.Fatalf("根据ID删除用户失败,数据仍然存在")
		}
	})

	t.Run("QueryByPage", func(t *testing.T) {
		qm := map[string]any{
			"name like": "张三",
		}
		ens, i, err := repo.QueryByPage(qm, &query.PageQBase{
			Skip:   1,
			Limit:  2,
			Sort:   "created_at",
			SortBy: "Desc",
		}, true)
		if err != nil || len(ens) != 2 {
			t.Fatalf("page query error")
		}
		for _, en := range ens {
			fmt.Println(fmt.Sprintf("count:%d, en: %v", i, en))
		}
	})

	t.Run("DeleteBy", func(t *testing.T) {
		um := map[string]any{
			"name like": "张三",
			"status":    0,
		}
		err, count := repo.DeleteBy(um, true)
		t.Logf("受影响行数:%d", count)
		if err != nil {
			t.Fatalf("根据ID删除用户失败: %v", err)
		}
		if qu, err := repo.FindBy(um, true); err != nil {
			t.Fatalf("查询用户失败: %v", err)
		} else if len(qu) > 0 {
			t.Fatalf("数据仍然存在")
		}
	})
	// 真删除
	t.Cleanup(func() {
		repo.gm.Exec(fmt.Sprintf("DELETE FROM %s WHERE ID < 6", base.UserGorm{}.TableName()))
	})
}

func NewConn(t *testing.T) *UserMysqlRepo {
	mysqlConfig := config.MysqlConfig{
		User:     "xxx",
		Password: "xxx",
		Host:     "xxx",
		DbName:   "xxx",
	}
	mysqlConfig.Sanitize()
	mysqlConfig.Validate()
	conn, err := mysql.NewMysqlConn(&mysqlConfig)
	if err != nil {
		t.Fatalf("初始化数据库连接失败:%v", err)
	}

	repo := NewUserMysqlRepo(conn)
	return repo
}
