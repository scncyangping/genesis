// @Author: YangPing
// @Create: 2023/10/23
// @Description: 转换工具类

package util

import (
	"reflect"
	"strings"
	"time"
	"unicode"

	"github.com/fatih/structs"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
)

func Copy(target, source any) error {
	err := copier.CopyWithOption(target, source, copier.Option{
		IgnoreEmpty: false,
		DeepCopy:    false,
		Converters: []copier.TypeConverter{
			{
				SrcType: time.Time{},
				DstType: copier.String,
				Fn: func(src interface{}) (interface{}, error) {
					s, ok := src.(time.Time)
					if !ok {
						return nil, errors.New("src type not matching")
					}
					return s.Format("2006-01-02 15:04:05"), nil
				},
			},
		},
	})

	if err != nil {
		return errors.Wrap(err, "Copy Struct Error")
	}
	return nil
}

func DeepCopy(target, source any) error {
	return copier.CopyWithOption(target, source, copier.Option{DeepCopy: true})
}

func CopyIgnoreEmpty(target, source any) error {
	return copier.CopyWithOption(target, source, copier.Option{IgnoreEmpty: true})
}

func DeepCopyIgnoreEmpty(target, source any) error {
	return copier.CopyWithOption(target, source, copier.Option{DeepCopy: true, IgnoreEmpty: true})
}

//func init() {
//	structs.DefaultTagName = "json"
//}

func StructToMap(s any) map[string]any {
	return structs.Map(s)
}

func StructCopy(DstStructPtr any, SrcStructPtr any) error {
	a := reflect.ValueOf(SrcStructPtr)
	b := reflect.ValueOf(DstStructPtr)
	c := reflect.TypeOf(SrcStructPtr)
	d := reflect.TypeOf(DstStructPtr)
	if c.Kind() != reflect.Ptr || d.Kind() != reflect.Ptr ||
		c.Elem().Kind() == reflect.Ptr || d.Elem().Kind() == reflect.Ptr {

		return errors.New("StructCopy error:type of parameters must be Ptr of value")
	}
	if a.IsNil() || b.IsNil() {
		return errors.New("StructCopy error:value of parameters should not be nil")
	}
	srcV := a.Elem()
	dstV := b.Elem()
	fields := deepFields(reflect.ValueOf(SrcStructPtr).Elem().Type())
	for _, v := range fields {
		if v.Anonymous {
			continue
		}
		dst := dstV.FieldByName(v.Name)
		src := srcV.FieldByName(v.Name)
		if !dst.IsValid() {
			continue
		}
		if src.Type() == dst.Type() && dst.CanSet() {
			dst.Set(src)
			continue
		}
		if src.Kind() == reflect.Ptr && !src.IsNil() && src.Type().Elem() == dst.Type() {
			dst.Set(src.Elem())
			continue
		}
		if dst.Kind() == reflect.Ptr && dst.Type().Elem() == src.Type() {
			dst.Set(reflect.New(src.Type()))
			dst.Elem().Set(src)
			continue
		}
	}
	return nil
}

func deepFields(baseType reflect.Type) []reflect.StructField {
	var fields []reflect.StructField

	for i := 0; i < baseType.NumField(); i++ {
		v := baseType.Field(i)
		if v.Anonymous && v.Type.Kind() == reflect.Struct {
			fields = append(fields, deepFields(v.Type)...)
		} else {
			fields = append(fields, v)
		}
	}
	return fields
}

// Camelize
// name: convert string
// big: whether the first letter is capitalized
func Camelize(name string, big bool) string {
	temp := strings.Split(name, "_")
	var s string
	for i, v := range temp {
		vv := []rune(v)
		if len(vv) > 0 {
			if !(!big && i == 0) {
				if vv[0] >= 'a' && vv[0] <= 'z' { //首字母大写
					vv[0] -= 32
				}
			}
			s += string(vv)
		}
	}
	return s
}

func UnCamelize(name string) string {
	buffer := strings.Builder{}

	for i, r := range name {
		if unicode.IsUpper(r) {
			if i != 0 {
				buffer.WriteRune('_')
			}
			buffer.WriteRune(unicode.ToLower(r))
		} else {
			buffer.WriteRune(r)
		}
	}
	return buffer.String()
}

func ArrayInGroupsOf[T any](arr []T, num int) [][]T {
	max := int(len(arr))
	//判断数组大小是否小于等于指定分割大小的值，是则把原数组放入二维数组返回
	if max <= num {
		return [][]T{arr}
	}
	//获取应该数组分割为多少份
	var quantity int
	if max%num == 0 {
		quantity = max / num
	} else {
		quantity = (max / num) + 1
	}
	//声明分割好的二维数组
	var segments = make([][]T, 0)
	//声明分割数组的截止下标
	var start, end, i int
	for i = 1; i <= quantity; i++ {
		end = i * num
		if i != quantity {
			segments = append(segments, arr[start:end])
		} else {
			segments = append(segments, arr[start:])
		}
		start = i * num
	}
	return segments
}

func String(s string) *string {
	return &s
}
