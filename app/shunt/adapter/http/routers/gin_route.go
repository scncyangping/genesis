// @Author: YangPing
// @Create: 2023/10/23
// @Description: 全局控制类

package routers

import "genesis/app/shunt/adapter/http/handlers/business"

// 使用Wire自动构造
// 将构造完成的对象注入到Gin Router进行路由配置

type Handlers struct {
	AuthHandler *business.AuthHandler
}
