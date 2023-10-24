// @Author: YangPing
// @Create: 2023/10/19
// @Description: 命令参数定义

package cmd

type UserSaveCmd struct {
	Name   string `json:"name" binding:"required" msg:"name不能为空"`
	Age    int    `json:"age" binding:"max=10,min=5" msg:"age必须大于5小于10"` // gte=5,lte=10
	Status int    `json:"status" binding:"required,oneof=1 2" msg:"status状态值错误"`
}

type UserUpdateCmd struct {
	Id string `json:"id" binding:"required"`
	UserSaveCmd
}
