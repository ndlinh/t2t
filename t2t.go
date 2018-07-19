package main

import (
	"flag"

	"github.com/jmoiron/sqlx"
	"github.com/ndlinh/t2t/dialect"
)

func main() {
	pk := flag.String("p", "", "Package name")
	table := flag.String("t", "", "Optional. Tables to export. Use commas to separate table's name")
	output := flag.String("o", "", "Optional. Output directory. Use current directory if empty")
	dsn := flag.String("dsn", "", "DB connect DSN. Ex: root:root@localhost/flarum?parseTime=true")

	flag.Parse()

	if *dsn == "" {
		flag.PrintDefaults()
		return
	}

	if *pk == "" {
		flag.PrintDefaults()
		return
	}

	if *output == "" {
		*output = "./"
	}

	db, err := sqlx.Connect("mysql", *dsn)
	if err != nil {
		panic(err)
	}

	d := dialect.NewMySQLDialect(*pk, *output, db)
	if *table == "" {
		d.BuildStruct()
	} else {
		d.BuildStructForTable(*table)
	}
}
