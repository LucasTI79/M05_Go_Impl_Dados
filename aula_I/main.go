package main

import (
	_ "github.com/go-sql-driver/mysql"
	db "github.com/lucasti79/go-impl-dados/database"
)

func main() {
	db.Init()

}
