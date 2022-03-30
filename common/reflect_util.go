package common

import (
	"fmt"
	"reflect"
	"strings"
)

// GetFieldName 获取结构体中字段的名称
func GetFieldName(columnName string, info interface{}) interface{} {
	var val interface{}
	t := reflect.TypeOf(info)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		fmt.Println("Check type error not Struct")
		return nil
	}
	fieldNum := t.NumField()
	for i := 0; i < fieldNum; i++ {
		name := t.Field(i).Name
		if strings.ToUpper(name) == strings.ToUpper(columnName) {
			v := reflect.ValueOf(info).Elem()
			val := v.FieldByName(name)
			return val.String()
		}
	}
	return val
}
