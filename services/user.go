package services

import (
	"time"

	"github.com/jinzhu/gorm"
)

// UserService represents a service to manage users
type UserService struct {
	db *gorm.DB
}

// NewUserService returns a new instance of UserService
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db: db,
	}
}

// User represents a user in the system
type User struct {
	ID           uint       `json:"id"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
	DeletedAt    *time.Time `json:"deletedAt,omitempty"`
	Username     string     `json:"username" gorm:"unique"`
	PasswordHash string     `json:"-"`
	IsAdmin      bool       `json:"isAdmin" gorm:"default:false"`
}

// Create a new user
func (s *UserService) Create(user *User) error {
	err := s.db.Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

// List all users
func (s *UserService) List() ([]User, error) {
	var users []User
	err := s.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// FindUserById ...
func (s *UserService) FindUserById(id int) (*User, error) {
	user := &User{}
	err := s.db.Where("id= ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

// FindUserByUsername ...
func (s *UserService) FindUserByUsername(username string) (*User, error) {
	user := &User{}
	err := s.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UserExists ...
func (s *UserService) UserExists(username string) bool {
	user := &User{}
	if s.db.Where("username = ?", username).First(&user).RecordNotFound() {
		return false
	}
	return true
}
