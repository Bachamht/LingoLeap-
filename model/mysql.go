package db

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"os"
)

var DB *sql.DB

func ConnectMysql(dbURL string) *sql.DB {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("env load error:", err)
	}

	MySQLUser := os.Getenv("MySQLUser")
	MySQLPort := os.Getenv("MySQLPort")
	MySQLPassword := os.Getenv("MySQLPassword")
	MySQLHost := os.Getenv("MySQLHost")
	MySQLDbname := os.Getenv("MySQLDbname")
	dbURL := MySQLUser + ":" + MySQLPassword + "@tcp(" + MySQLHost + ":" + MySQLPort + ")/" + MySQLDbname

	DB, err := sql.Open("mysql", dbURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	fmt.Println("Connected to database successfully.")

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