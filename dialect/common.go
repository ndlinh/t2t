package dialect

import (
	"database/sql"
	"strings"
)

// TField - table field information
type TField struct {
	Field   string
	Type    TType
	Null    string
	Key     string
	Default sql.NullString
	Extra   string
}

// TType - type information
type TType struct {
	Name  string
	Size  int
	Extra string
}

// Dialect - common dialect
type Dialect interface {
	BuildStruct() string
}

// StructTemplate is template for genrate file output
const StructTemplate = `
package %s

%s// %s - 
type %s struct {
%s}
`

// CamelCase convert table name to cammel case
func CamelCase(s string) string {
	list := strings.Split(s, "_")
	for i := 0; i < len(list); i++ {
		val := list[i]
		check := strings.ToLower(val)
		if check == "id" || check == "api" || check == "ip" {
			list[i] = strings.ToUpper(val)
		} else {
			list[i] = strings.Title(val)
		}
	}

	return strings.Join(list, "")
}
