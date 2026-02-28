// Database connection
package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// func ConnectDB() {
// dsn := "host=localhost user=postgres password=Msdhoni@7 dbname=todo_db port=5432 sslmode=disable"
// 	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		panic("Database connection failed")
// 	}

// 	database.AutoMigrate(&User{}, &Todo{})
// 	DB = database
// }
func ConnectDB() error {
    dsn := "host=localhost user=postgres password=Msdhoni@7 dbname=todo_db port=5432 sslmode=disable"
    database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return err
    }

    database.AutoMigrate(&User{}, &Todo{})
    DB = database
    return nil
}