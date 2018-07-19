# Introduction

This package provide helper functions for generate go struct from MySQL table.
Generated structs that can be use with https://github.com/jmoiron/sqlx/ package.

# How to use

## Installation

    go get github.com/jmoiron/sqlx

    go get -u github.com/ndlinh/t2t

# Usage

    t2t -dsn <dsn> -p <package name> [-o output directory] [-t tables]

    Options:

    -dsn string
    	DB connect DSN. Ex: root:root@localhost/flarum
    -o string
            Optional. Output directory. Use current directory if empty
    -p string
            Package name
    -t string
            Optional. Tables to export. Use commas to separate table's name


    Ex: t2t -p "test" -o "./test" -t "users_tags" -dsn "root:root@localhost/flarum_test"

# Todos

- [ ] Writing tests
- [ ] PostgreSQL dialect
