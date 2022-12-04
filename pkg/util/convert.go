package util

import (
	"github.com/fatih/structs"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"reflect"
	"time"
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

func StructToMap(s any) any {
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
