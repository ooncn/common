package util

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

func DBQuerySql(sqlSer *gorm.DB, sql string, values ...interface{}) (list []map[string]string, err error) {
	rows, err := sqlSer.Raw(sql, values...).Rows() // (*sql.Rows, error)
	//defer rows.Close()
	if err != nil {
		return
	}
	columns, _ := rows.Columns()
	columnLength := len(columns)
	cache := make([]interface{}, columnLength) //临时存储每行数据
	for index, _ := range cache {              //为每一列初始化一个指针
		var a interface{}
		cache[index] = &a
	}
	var o string
	for rows.Next() {
		err = rows.Scan(cache...)
		if err != nil {
			return
		}
		item := make(map[string]string)
		for i, data := range cache {
			bytes := *data.(*interface{})
			switch bytes.(type) {
			case int64:
				o = fmt.Sprintf("%d", bytes.(int64))
			case float64:
				o = fmt.Sprintf("%f", bytes.(float64))
			case []uint8:
				o = UiToS(bytes.([]uint8))
			case time.Time:
				o = TimeUtil.DateToyMdHmsSepTo(bytes.(time.Time))
			default:
				o = ""
			}
			item[columns[i]] = o //取实际类型
		}
		list = append(list, item)
	}
	return
}
func UiToS(bs []uint8) string {
	var ba = make([]byte, 0)
	for _, b := range bs {
		ba = append(ba, b)
	}
	return string(ba)
}
