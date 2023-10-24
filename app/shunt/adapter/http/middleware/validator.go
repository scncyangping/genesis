// @Author: YangPing
// @Create: 2023/10/23
// @Description: 参数校验构建

package middleware

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// required: 指定字段为必填字段，不能为空。
// len: 指定字符串、切片或数组的长度必须满足特定要求。
// min: 指定数字的最小值要求。
// max: 指定数字的最大值要求。
// eq 和 ne: 指定相等和不相等的值要求。
// oneof: 要求字段的值必须是给定的一组值之一。
// alphanum: 要求字段只包含字母和数字字符。
// alpha: 要求字段只包含字母字符。
// numeric: 要求字段只包含数字字符。
// email: 要求字段是有效的电子邮件地址。
// url: 要求字段是有效的 URL。
// ip: 要求字段是有效的 IP 地址。
// iscolor: 要求字段是有效的颜色值。
// mac: 要求字段是有效的 MAC 地址。
// file: 要求字段是有效的文件路径。
// http_url: 要求是有效的http地址
func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("notEmpty", notEmpty)
	}
}

// 定义一个自定义验证函数
func notEmpty(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return value != ""
}
