package util

import (
	"encoding/json"
	"fmt"
	"reflect"
)

/**
json 转指定对象
*/
func JsonToType(jsonStr string, i interface{}) (err error) {
	err = json.Unmarshal([]byte(jsonStr), i)
	return err
}
func JsonTypeToType(json interface{}, i interface{}) error {
	return JsonToType(JsonToStr(json), i)
	//fmt.Println(i)
}
func jsonDataTypeConv(m map[string]interface{}) {
	for k, v := range m {
		switch vv := v.(type) {
		case string:
			fmt.Println(k, "is string", vv)
		case float64:
			fmt.Println(k, "is float", int64(vv))
		case int:
			fmt.Println(k, "is int", vv)
		case []interface{}:
			fmt.Println(k, "is an array:")
			for i, u := range vv {
				fmt.Println(i, u)
			}
		case nil:
			fmt.Println(k, "is nil", "null")
		case map[string]interface{}:
			fmt.Println(k, "is an map:")
			jsonDataTypeConv(vv)
		default:
			fmt.Println(k, "is of a type I don't know how to handle ", fmt.Sprintf("%T", v))
		}
	}
}

/**
json 转指定字符串
*/
func JsonToStr(i interface{}) string {

	t := reflect.TypeOf(i)
	switch t.Name() {
	case "string":
		return i.(string)
	}
	data, _ := json.Marshal(i)
	//fmt.Println(data)
	return string(data)
}
func JsonToMap(i string) map[string]interface{} {
	m := make(map[string]interface{})
	JsonToType(i, &m)
	return m
}
