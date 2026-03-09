package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB
func ConnectDB(){
	username := "root" // Sesuaikan dengan username di Komputer Anda
	password := "" // Sesuaikan dengan password di Komputer Anda
	server := "127.0.0.1:3306" // Sesuaikan dengan servername di Komputer Anda
	dbName := "tokolepkom_npm"
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true",username, password, server, dbName)
	db, err := sql.Open("mysql", dsn)
	
	if err != nil {
		log.Fatal("Failed open DB:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Database not connected:", err)
	}

	log.Println("\nDatabase connected!")
	DB = db
}