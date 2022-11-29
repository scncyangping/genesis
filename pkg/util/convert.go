package util

import (
	"github.com/pkg/errors"
	"reflect"
)

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
