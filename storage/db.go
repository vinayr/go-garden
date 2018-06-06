package storage

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/vinayr/go-garden/services"
	"golang.org/x/crypto/bcrypt"
	// postgres
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Database represents handle to database
type Database struct {
	DB   *gorm.DB
	Path string
}

// NewDB returns a new instance of Database
func NewDB() *Database {
	return &Database{}
}

// Open and initialize the database
func (db *Database) Open() error {
	d, err := gorm.Open("postgres", db.Path)
	if err != nil {
		return err
	}

	if err := d.DB().Ping(); err != nil {
		return err
	}

	d.LogMode(true)
	db.DB = d

	return nil
}

// Close the database
func (db *Database) Close() error {
	if db.DB != nil {
		db.DB.Close()
	}
	return nil
}

// Migrate all schemas
func (db *Database) Migrate() {
	db.DB.AutoMigrate(
		services.User{},
	)

	// Initial admin setup
	username := "admin@test.com"
	password := "admin"
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := &services.User{
		Username:     username,
		PasswordHash: string(passwordHash),
		IsAdmin:      true,
	}
	err := db.DB.FirstOrCreate(user, &services.User{Username: username}).Error
	if err != nil {
		log.Print("Create admin error: ", err)
	}
}
