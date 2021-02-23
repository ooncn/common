package util

import (
	"strconv"
	"testing"

	"github.com/pikanezi/mapslice"
	"github.com/tealeg/xlsx"
)

func TestExcel(t *testing.T) {
	var role RoleBack
	var roleLists []RoleBack
	roleList, _ := GetRoleList()
	for i := 0; i < len(roleList); i++ {
		role.Id = strconv.Itoa(roleList[i].Id)
		role.Name = roleList[i].Name
		role.Remark = roleList[i].Remark
		role.Status = strconv.Itoa(roleList[i].Status)
		roleLists = append(roleLists, role)
	}
	id, _ := mapslice.ToStrings(roleLists, "Id")
	name, _ := mapslice.ToStrings(roleLists, "Name")
	remark, _ := mapslice.ToStrings(roleLists, "Remark")
	status, _ := mapslice.ToStrings(roleLists, "Status")

	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	file = xlsx.NewFile()
	sheet, _ = file.AddSheet("Sheet1")
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.Value = "编号"
	cell = row.AddCell()
	cell.Value = "名称"
	cell = row.AddCell()
	cell.Value = "状态"
	cell = row.AddCell()
	cell.Value = "备注"
	for i := 0; i < len(id); i++ {
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.Value = id[i]
		cell = row.AddCell()
		cell.Value = name[i]
		cell = row.AddCell()
		cell.Value = status[i]
		cell = row.AddCell()
		cell.Value = remark[i]
		file.Save("File.xlsx")
	}
}
