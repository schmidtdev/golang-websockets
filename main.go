package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"schmidtdev/golang-websockets/handlers"
	"schmidtdev/golang-websockets/types"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to database
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"), os.Getenv("SSL_MODE"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	// Migrate the schema
	db.AutoMigrate(&types.Channel{}, &types.Subscription{})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handlers.WsHandler(w, r, db)
	})

	fmt.Println("Listening on :3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		fmt.Println("Server error:", err)
	}
}
