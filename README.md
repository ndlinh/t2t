# Introduction

This package provide helper functions for generate go struct from MySQL table.
Generated structs can be use with `https://github.com/jmoiron/sqlx/issues/135` package.

# How to use

## Installation

    go get github.com/jmoiron/sqlx

    go get -u github.com/ndlinh/t2t

# Usage

    t2t -dsn <dsn> -p <package name> [-o output directory] [-t tables]

    Ex: t2t -p "test" -o "./test" -t "users_tags" -dsn "root:root@localhost/flarum_test"

# Todos

- [ ] Writing tests
- [ ] PostgreSQL dialect
