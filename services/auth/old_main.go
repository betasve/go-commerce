// The main package for user authentication
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/betasve/go-commerce/services/auth/router"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func migrateDb() {
	m, err := migrate.New(
		fmt.Sprintf("file:///%s/db/migrations/", os.Getenv("SERVICE_ROOT_PATH")),
		os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed while instantiating migration object: %v", err)
	}

	defer m.Close()

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed while migrating the database: %v", err)
	}
}

func connectDb() *gorm.DB {
	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the db: %v", err)
	}

	return db
}

func main() {
	migrateDb()
	// db := connectDb()

	router.Run()
}
