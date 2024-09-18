package models

import (
	"fmt"
	"slices"
	"strings"
)

type TableData struct {
	Headers []tableColumn
	Rows    [][]string
}

// *SORT ENUMS
type tableSortDirection string
type TableSortConfig struct {
	Fieldname string
	Direction tableSortDirection
}

const sortAsc = tableSortDirection("ASC")
const sortDesc = tableSortDirection("DESC")

// * COLUMN FIELD ENUMS
type tableColumn struct {
	Label     string
	Sort      tableSortDirection
	Fieldname string
	Data      func(User) string
	Centered  bool
}
type ColumnConfig []tableColumn

// user table columns
func NameUC() tableColumn {
	return tableColumn{
		Label:     "Name",
		Fieldname: "full_name",
		Data: func(u User) string {
			return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
		},
	}
}

func BusinessUC() tableColumn {
	return tableColumn{
		Label:     "Business",
		Fieldname: "business",
		Data: func(u User) string {
			return "-"
		},
	}
}

func AdminUC() tableColumn {
	return tableColumn{
		Label:     "Admin",
		Fieldname: "is_admin",
		Data: func(u User) string {
			if slices.Contains(u.Permissions, "admin") {
				return "✓"
			}
			return "✕"
		},
		Centered: true,
	}
}

func EmailUC() tableColumn {
	return tableColumn{
		Label:     "Email",
		Fieldname: "email_address",
		Data: func(u User) string {
			return string(u.EmailAddress)
		},
	}
}

func (cc *ColumnConfig) ApplySortingQuery(query string) *TableSortConfig {
	if query == "" {
		return nil
	}

	colLabel := query
	direction := sortAsc
	if strings.HasPrefix(query, "-") {
		colLabel = strings.TrimPrefix(query, "-")
		direction = sortDesc
	}

	for i, col := range *cc {
		if col.Label == colLabel {
			(*cc)[i].Sort = direction
			return &TableSortConfig{
				Fieldname: col.Fieldname,
				Direction: direction,
			}
		}

	}

	return nil
}
