package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // Driver PostgreSQL untuk Go
)

// Variabel global agar database bisa diakses dari folder/file lain
var DB *sql.DB

func ConnectDB() {
	dbUrl := os.Getenv("DATABASE_URL")
	var connStr string

	if dbUrl != "" {
		// Jika dideploy di Cloud (Render), gunakan URL langsung dari Supabase
		connStr = dbUrl
	} else {
		// Fallback untuk local Docker Compose
		dbHost := os.Getenv("DB_HOST")
		if dbHost == "" { dbHost = "localhost" }
		dbUser := os.Getenv("DB_USER")
		if dbUser == "" { dbUser = "postgres" }
		dbPass := os.Getenv("DB_PASS")
		if dbPass == "" { dbPass = "admin123" }
		
		connStr = fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=gotracker sslmode=disable", dbHost, dbUser, dbPass)
	}
	
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Gagal membuka koneksi database: ", err)
	}

	// Tes koneksi (Ping) untuk memastikan PostgreSQL benar-benar menyala dan merespons
	err = DB.Ping()
	if err != nil {
		log.Fatal("Database tidak merespons (Ping gagal): ", err)
	}

	fmt.Println("✅ Berhasil terhubung ke PostgreSQL (Database: gotracker)!")
	
	// Secara otomatis membuat tabel jika belum ada (Sangat praktis!)
	createTables()
}

// Fungsi internal untuk auto-create tabel
func createTables() {
	query := `
	CREATE TABLE IF NOT EXISTS websites (
		id SERIAL PRIMARY KEY,
		url VARCHAR(255) NOT NULL UNIQUE
	);
	
	CREATE TABLE IF NOT EXISTS checks (
		id SERIAL PRIMARY KEY,
		website_id INT REFERENCES websites(id),
		status VARCHAR(10) NOT NULL,
		checked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	
	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal("Gagal membuat tabel: ", err)
	}
	fmt.Println("✅ Tabel database siap digunakan.")
}
