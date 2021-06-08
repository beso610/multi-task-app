package repository

import (
	"fmt"
	"os"
	"strconv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var tables = []interface{}{
	&User{},
	&Task{},
}

type Repository interface {
	UserRepository
	TaskRepository
}

// DBRepository implements Repository
type DBRepository struct {
	DB *gorm.DB
}

func SetupDB() (*gorm.DB, error) {
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "root"
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "password"
	}

	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}

	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		port = 3306
	}

	dbname := os.Getenv("DB_DATABASE")
	if dbname == "" {
		dbname = "multitask"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, dbname)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return db, err
	}

	tx := db.Begin()

	if err := tx.AutoMigrate(tables...); err != nil {
		tx.Rollback()
		return db, fmt.Errorf("failed to sync table schema: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		return db, fmt.Errorf("failed to commit setup db: %v", err)
	}

	return db, nil
}
