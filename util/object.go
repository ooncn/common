package util

import (
	"fmt"
	"reflect"
	"strconv"
)

/*
获取对象参数,将参数变量指针地址拼接成变量数组返回
*/
func ObjectParameterArray(obj interface{}) ([]interface{}, error) {
	value := reflect.ValueOf(obj)
	// value.Elem() 遍历对象的值
	// value.Elem().NumField() 获取对象中的属性数量
	var err error
	fmt.Println(value.Elem().NumField())
	s := value.Elem() //遍历对象的值

	length := s.NumField() //获取值的数量
	result := make([]interface{}, 0)
	row := make([]interface{}, length)
	for i := 0; i < length; i++ {
		f := s.Field(i)        // 根据索引获取索引所对应的的属性变量
		a := f.Addr()          //获得其地址
		row[i] = a.Interface() //以空接口类型获得值
		result = append(result, s.Interface())
	}
	fmt.Println(row)
	fmt.Println(result)
	return row, err
}

func ObjToStr(i interface{}) (s string) {
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)
	switch t.Kind() {
	case reflect.String:
		s = v.String()
	case reflect.Bool:
		if v.Bool() {
			s = "true"
		} else {
			s = "false"
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		s = strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		s = strconv.FormatUint(v.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		s = strconv.FormatFloat(v.Float(), 'f', 5, 32)
	case reflect.Ptr:
		s = ObjToStr(v.Elem())
	}
	return
}
