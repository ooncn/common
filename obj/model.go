package obj

import (
	"errors"
	"gorm.io/gorm"
)

var ErrorMap = map[int]ErrorObj{
	2001: {2001, "UI_NOT_DATA", "防伪信息错误"},   //找不到数据
	2002: {2002, "UI_STATUS_ERR", "防伪信息错误"}, //状态错误
	2003: {2003, "UI_END_TIME", "防伪码信息过期"},  //到期
	2004: {2004, "UI_START_TIME", "防伪信息错误"}, //未到开始时间
	2005: {2005, "UI_NOT_ID", "防伪信息错误"},     //找不到数据
	2006: {2006, "UI_PASS_CHECK", "防伪信息错误"}, //找不到数据
	2007: {2007, "UI_NOT_PASS", "请输入防伪密码"},  //找不到数据
	2008: {2008, "UI_PASS_ERR", "防伪密码错误"},   //找不到数据
}

type Model struct {
	SqlSer *gorm.DB `json:"-" gorm:"-"`
}

type ErrorObj struct {
	Code int
	Msg  string
	Data string
}

func (w *Model) Verification() interface{} {
	return nil
}

// 查询
func (m *Model) SetDB(i *gorm.DB) {
	m.SqlSer = i
}
func (m *Model) GetDB() {
}

// 查询
func (m *Model) Create(i interface{}) error {
	return m.SqlSer.Create(i).Error
}

// 更新
func (m *Model) Updates(i interface{}) error {
	return m.SqlSer.Updates(i).Error
}

// 删除
func (m *Model) Delete(id string) error {
	return m.SqlSer.Delete(id).Error
}

// 校验
func (m *Model) Check() error {
	return errors.New("")
}

// 校验
/*func (m *Model) Save(i interface{}) error {
	return m.SqlSer.Save(i).Error
}
*/

type Option struct {
	Id   string `json:"id"`   //ID
	Name string `json:"name"` //名称
}

type OptionPid struct {
	Id   string `json:"id"`   //ID
	Pid  string `json:"pid"`  //PiD
	Name string `json:"name"` //名称
}

type OptionMenu struct {
	Id   string `json:"id"`   //ID
	Pid  string `json:"pid"`  //PiD
	Icon string `json:"icon"` //PiD
	Name string `json:"name"` //名称
}

type PidSelect struct {
	Id    string `json:"id"`
	Pid   string `json:"pid"`
	Name  string `json:"name"`
	Value string `json:"value"`
}
