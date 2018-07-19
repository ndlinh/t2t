package dialect

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var fieldMaps = map[string]string{
	"int":       "int64",
	"tinyint":   "int",
	"smallint":  "int",
	"mediumint": "int",
	"bigint":    "int64",

	"int.unsigned":       "uint64",
	"tinyint.unsigned":   "uint",
	"smallint.unsigned":  "uint",
	"mediumint.unsigned": "uint",
	"bigint.unsigned":    "int64",

	"float":      "float64",
	"decimal":    "float64",
	"double":     "float64",
	"timestamp":  "int64",
	"varchar":    "string",
	"text":       "string",
	"longtext":   "string",
	"mediumtext": "string",
	"datetime":   "time.Time",
	"date":       "time.Time",
	"bit":        "bool",
	"enum":       "string",
	"blob":       "string",
}

var nullableFieldMaps = map[string]string{
	"int":       "sql.NullInt64",
	"tinyint":   "sql.NullInt64",
	"smallint":  "sql.NullInt64",
	"mediumint": "sql.NullInt64",
	"bigint":    "sql.NullInt64",

	"int.unsigned":       "sql.NullInt64",
	"tinyint.unsigned":   "sql.NullInt64",
	"smallint.unsigned":  "sql.NullInt64",
	"mediumint.unsigned": "sql.NullInt64",
	"bigint.unsigned":    "sql.NullInt64",

	"timestamp":  "sql.NullInt64",
	"varchar":    "sql.NullString",
	"text":       "sql.NullString",
	"longtext":   "sql.NullString",
	"mediumtext": "sql.NullString",
	"float":      "sql.NullFloat64",
	"datetime":   "mysql.NullTime",
	"date":       "mysql.NullTime",
	"bit":        "bool",
	"enum":       "sql.NullString",
	"blob":       "sql.NullString",
}

// MySQLDialect -
type MySQLDialect struct {
	pk     string
	output string
	db     *sqlx.DB
}

// NewMySQLDialect -
func NewMySQLDialect(pk string, output string, db *sqlx.DB) *MySQLDialect {
	//ensure output folder
	if _, err := os.Stat(output); os.IsNotExist(err) {
		os.Mkdir(output, 0777)
	}
	return &MySQLDialect{pk: pk, output: output, db: db}
}

// BuildStruct -
func (d *MySQLDialect) BuildStruct() string {
	tables := d.getTables()

	for i := range tables {
		d.BuildStructForTable(tables[i])
	}

	return ""
}

// BuildStructForTable -
func (d *MySQLDialect) BuildStructForTable(table string) {
	tables := strings.Split(table, ",")

	for i := 0; i < len(tables); i++ {
		t := tables[i]
		tableName := CamelCase(t)
		str, importStr := d.buildTableStruct(t)
		fileStr := fmt.Sprintf(StructTemplate, d.pk, importStr, tableName, tableName, str)

		filePath := fmt.Sprintf("%s/%s.go", d.output, t)
		f, err := os.Create(filePath)
		if err != nil {
			panic("Can not write file")
		}
		defer f.Close()

		f.WriteString(fileStr)
	}
}

func (d *MySQLDialect) buildTableStruct(table string) (string, string) {
	rows := d.getTableFields(table)
	rowTemplate := "    %s %s `db:\"%s\"`\n"

	var str = ""
	var imports = make(map[string]string)
	for r := range rows {
		f := rows[r]
		str += fmt.Sprintf(rowTemplate, CamelCase(f.Field), d.getMapFieldType(f), f.Field)
		if f.Type.Name == "datetime" || f.Type.Name == "date" {
			if f.Null != "NO" {
				imports["mysql"] = `    "github.com/go-sql-driver/mysql"`
			} else {
				imports["time"] = `    "time"`
			}
		} else if f.Null != "NO" {
			imports["sql"] = `    "database/sql"`
		}
	}
	var importList = []string{}
	for _, v := range imports {
		importList = append(importList, v)
	}
	if len(importList) > 0 {
		importStr := "import (\n" + strings.Join(importList, "\n") + "\n)\n\n"
		return str, importStr
	}

	return str, ""
}

// GetTables return all table in database
func (d *MySQLDialect) getTables() []string {
	query := "SHOW TABLES"
	rows, err := d.db.Query(query)
	defer rows.Close()

	if err != nil {
		log.Fatal(err)
	}

	var result []string
	for rows.Next() {
		var name string
		e := rows.Scan(&name)
		if e != nil {
			log.Fatal(e)
			break
		}

		result = append(result, name)
	}

	return result
}

// GetTableFields return all fields of tables
func (d *MySQLDialect) getTableFields(name string) []TField {
	query := "DESCRIBE " + name
	rows, err := d.db.Query(query)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}

	var ret = []TField{}
	for rows.Next() {
		var (
			field, ftype, key, isNull, extra string
			isDefault                        sql.NullString
		)
		err := rows.Scan(&field, &ftype, &isNull, &key, &isDefault, &extra)
		if err != nil {
			log.Fatal(err)
		}

		tt := d.extractType(ftype)
		if tt == nil {
			log.Printf("Can not determine type: %s of field %s", ftype, field)
			continue
		}

		var r = TField{
			Field:   field,
			Type:    *tt,
			Key:     key,
			Null:    isNull,
			Extra:   extra,
			Default: isDefault,
		}

		ret = append(ret, r)
	}

	return ret
}

// ExtractType - extract type to TType
func (d *MySQLDialect) extractType(ttype string) *TType {
	reg := regexp.MustCompile(`([a-z]+)|(\d+)`)
	parts := reg.FindAllString(ttype, -1)
	l := len(parts)
	if l == 0 {
		return nil
	}

	if parts[0] == "enum" {
		return &TType{
			Name:  parts[0],
			Size:  0,
			Extra: strings.Join(parts[1:], ","),
		}
	}
	switch l {
	case 3:
		size, _ := strconv.Atoi(parts[1])
		return &TType{
			Name:  parts[0],
			Size:  size,
			Extra: parts[2],
		}
	case 2:
		size, _ := strconv.Atoi(parts[1])
		return &TType{
			Name:  parts[0],
			Size:  size,
			Extra: "",
		}
	default:
		return &TType{
			Name:  parts[0],
			Size:  0,
			Extra: "",
		}
	}
}

func (d *MySQLDialect) getMapFieldType(from TField) string {
	var key = from.Type.Name
	if from.Type.Extra == "unsigned" {
		key = key + ".unsigned"
	}

	var (
		val string
		ok  bool
	)
	if from.Null == "NO" {
		val, ok = fieldMaps[key]
	} else {
		val, ok = nullableFieldMaps[key]
	}

	if ok == true {
		return val
	}

	return "/* <error> */" + from.Type.Name
}
