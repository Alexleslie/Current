package main

import (
	"Current/tools/logc"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func main() {
	db, _ := sql.Open("sqlite3", "gee.db")
	defer func() { _ = db.Close() }()

	_, _ = db.Exec("DROP TABLE IF EXISTS USER;")
	_, _ = db.Exec("CREATE TABLE User(Name text);")
	result, err := db.Exec("INSERT INTO User('Name') values (?),(?)", "Tom", "Sam") // ?为占位符，一般用来防SQL注入
	if err == nil {
		affected, _ := result.RowsAffected()
		log.Println(affected)
	}
	row := db.QueryRow("SELECT Name FROM User LIMIT 1") // 只返回一条查询记录，类型是 *sql.Row，Query()返回多条查询记录
	var name string
	if err := row.Scan(&name); err == nil {
		logc.Info("%+v", name)
	}

}
