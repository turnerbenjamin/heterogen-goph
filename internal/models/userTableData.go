package models

import (
	"fmt"
	"slices"
)

type UserTableData struct {
	Headers []string
	Rows    [][]string
}

type userColumn string

const NameUC = userColumn("Name")
const BusinessUS = userColumn("Business")
const AdminUC = userColumn("Admin")
const EmailUC = userColumn("Email")

var columnDataMapping = map[userColumn]func(User) string{
	NameUC: func(u User) string {
		return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
	},
	BusinessUS: func(u User) string {
		return "-"
	},
	AdminUC: func(u User) string {
		if slices.Contains(u.Permissions, "admin") {
			return "✓"
		}
		return "✕"
	},
	EmailUC: func(u User) string {
		return string(u.EmailAddress)
	},
}

func GetUserTableData(users []User, columns ...userColumn) UserTableData {
	tableData := UserTableData{
		Headers: []string{},
		Rows:    [][]string{},
	}

	for _, col := range columns {
		tableData.Headers = append(tableData.Headers, string(col))
	}

	for _, user := range users {
		rowData := []string{}
		for _, col := range columns {
			getColData := columnDataMapping[col]
			rowData = append(rowData, getColData(user))
		}
		tableData.Rows = append(tableData.Rows, rowData)
	}

	return tableData

}

// var nameUserColumn column = column{
// 	name: "Name",
// 	getData: func(u User) string{
// 		return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
// 	},
// }
