// The main package for user authentication
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func connectDb() *gorm.DB {
	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	return db
}

func main() {
	db := connectDb()
	result := db.Raw("SELECT 1;")
	if result.Error != nil {
		log.Fatalf("Failed to execute query: %v", result.Error)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain")

		fmt.Fprint(w, "User Auth Service is running.")
	})

	port := fmt.Sprintf(":%s", os.Getenv("APP_PORT"))

	fmt.Printf("User Auth Service started on port %s\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
