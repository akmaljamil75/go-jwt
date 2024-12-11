package database

import (
	"fmt"
	"time"

	"go-jwt/modules/role"
	"go-jwt/modules/user"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbGlobal *gorm.DB

func InitDB() *gorm.DB {

	if err := godotenv.Load(); err != nil {
		panic("failed to load.env file")
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable search_path=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SCHEMA"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true})

	if err != nil {
		panic(err)
	}
	err = migrateDatabase(db)
	if err != nil {
		log.Fatal("failed to migrate database")
		panic("failed to migrate database")
	}
	log.Printf("Success Connection DB")

	sql, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get SQL DB instance: %v", err)
		panic("failed to get SQL DB instance")
	}

	dbGlobal = db
	sql.SetMaxIdleConns(100)                 // Maksimal 100 koneksi idle
	sql.SetMaxOpenConns(1000)                // Maksimal 1000 koneksi aktif
	sql.SetConnMaxIdleTime(10 * time.Minute) // Koneksi idle maksimal 10 menit
	sql.SetConnMaxLifetime(30 * time.Minute) // Koneksi maksimal bertahan 30 menit

	return db
}

func GetDB() *gorm.DB {
	if dbGlobal == nil {
		log.Fatal("Database is not initialized yet")
	}
	return dbGlobal
}

func migrateDatabase(db *gorm.DB) error {
	err := db.AutoMigrate(&user.User{}, &role.Role{})
	if err != nil {
		return err
	}
	return nil
}
