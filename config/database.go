package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDatabase() {
	var err error
	dsn := "host=localhost user=file_admin password=dhairya dbname=file_sharing sslmode=disable"

	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("❌ Could not connect to the database:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("❌ Database connection failed:", err)
	}

	fmt.Println("🚀 Database connected successfully!")
}
