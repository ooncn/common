package util

import (
	"gorm.io/gorm"
)

type Role struct {
	Id     int    `gorm:"primary_key"` //主键id
	Name   string `gorm:"not null"`    //角色名
	Remark string
	Status int
}

type RoleBack struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Remark string `json:"remark"`
	Status string `json:"status"`
}

//规定表名
func (Role) TableName() string {
	return "role"
}

var DB *gorm.DB

func GetRoleList() (roleList []Role, err error) {
	err = DB.Find(&roleList).Error
	return roleList, err
}
