package util

import (
	"errors"
	"strings"
)

func NewError(err string) error {
	return errors.New(err)
}
func StrError(err error) string {
	//e := err.Error();
	//t := reflect.TypeOf(err)
	//if t.Kind() == reflect.Ptr {
	//	if t.Elem().String() == "url.Error" {
	//		urlErr := (url.Error{})err
	//	};
	//}
	s := err.Error()
	if s == "record not found" {
		s = "未找到数据"
	} else if strings.Contains(s, "UNIQUE constraint failed") {
		s = "数据已存在"
	} else if strings.Contains(s, "value too long for type") || strings.Contains(s, "out of range") {
		s = "数据过长"
	} else if strings.Contains(s, "Duplicate entry") {
		s = "数据已存在"
	} else if strings.Contains(s, "Data too long for") {
		s = "数据长度超出"
	} else if strings.Contains(s, "net/http: request canceled (Client.Timeout exceeded while awaiting headers)") {
		s = "net/http 连接超时"
	}
	return s
}
