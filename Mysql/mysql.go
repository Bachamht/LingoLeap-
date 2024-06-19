package db

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB(dbURL string) *sql.DB {

	DB, err := sql.Open("mysql", dbURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	fmt.Println("Connected to database successfully.")
	fmt.Println("connect success")

	// 读取 SQL 脚本内容
	sqlScript, err := ioutil.ReadFile("./db/db.sql")
	if err != nil {
		log.Fatal(err)
	}

	// 执行 SQL 脚本

	statements := strings.Split(string(sqlScript), ";")

	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt != "" {
			_, err := DB.Exec(stmt)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	fmt.Println("SQL script executed successfully.")
	return DB
}