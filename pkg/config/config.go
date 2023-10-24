// @Author: YangPing
// @Create: 2023/10/23
// @Description: 配置文件实体基础接口

package config

type Config interface {
	Sanitize()
	Validate() error
}

type InitConfig interface {
	Init() error
}
