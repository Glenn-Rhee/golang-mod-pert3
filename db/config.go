package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB
func ConnectDB(){
	/* KONEKSI DI BAWAH UNTUK KAMPUS F4 */
	username := "root" 
	password := "" 
	server := "127.0.0.1:3306" 
	dbName := "tokolepkom_npm"

	/* KONEKSI DI BAWAH UNTUK KAMPUS f8 */
	username := "APCx" // Ganti nilai x dengan no PC masing - masing 
	password := "lepkom@123" 
	server := "dbms.lepkom.f4.com:3306" 
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