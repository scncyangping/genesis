// @Author: YangPing
// @Create: 2023/10/21
// @Description: 基础接口定义

package base

// EnI 实体基础接口
type EnI interface {
	// TableName 返回实体对应表名
	TableName() string
}
