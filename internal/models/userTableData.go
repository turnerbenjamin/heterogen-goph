package models

import (
	"github.com/turnerbenjamin/heterogen-go/internal/httpErrors"
)

var nameColumn = NameUC()
var businessColumn = BusinessUC()
var emailColumn = EmailUC()
var adminColumn = AdminUC()

var columnLabelMapping = map[string]tableColumn{
	nameColumn.Label:     nameColumn,
	businessColumn.Label: businessColumn,
	adminColumn.Label:    adminColumn,
	emailColumn.Label:    emailColumn,
}

var defaultColumnConfig = ColumnConfig{nameColumn, businessColumn, emailColumn, adminColumn}

func GetColumnConfig(columnLabels []string) (ColumnConfig, error) {
	if len(columnLabels) == 0 {
		return defaultColumnConfig, nil
	}

	columnConfig := ColumnConfig{}
	for _, columnLabel := range columnLabels {
		if col, ok := columnLabelMapping[columnLabel]; ok {
			columnConfig = append(columnConfig, col)
		} else {
			return nil, httpErrors.InvalidColumnConfig(columnLabel)
		}
	}
	return columnConfig, nil
}

func GetUserTableData(users []User, columns ColumnConfig) TableData {
	tableData := TableData{
		Headers: []tableColumn{},
		Rows:    [][]string{},
	}

	for _, col := range columns {
		tableData.Headers = append(tableData.Headers, col)
	}

	for _, user := range users {
		rowData := []string{}
		for _, col := range columns {
			rowData = append(rowData, col.Data(user))
		}
		tableData.Rows = append(tableData.Rows, rowData)
	}

	return tableData

}
