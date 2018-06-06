package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vinayr/go-garden/middleware"
	"github.com/vinayr/go-garden/services"
	"golang.org/x/crypto/bcrypt"
)

// UserHandler represents an HTTP handler for managing users
type UserHandler struct {
	// Services
	UserService *services.UserService
}

// NewUserHandler returns a new instance of userHandler
func NewUserHandler() *UserHandler {
	h := &UserHandler{}
	return h
}

// Signup JSON binding
type Signup struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Signup new user
func (h *UserHandler) Signup(c *gin.Context) {
	var params Signup
	if err := c.BindJSON(&params); err != nil {
		return
	}

	// Check if user already exists
	username := strings.ToLower(params.Username)
	if h.UserService.UserExists(username) {
		c.JSON(http.StatusConflict, gin.H{"message": "User already exists"})
		return
	}

	// Create password hash
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Print("Password hash error: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Create user
	user := &services.User{
		Username:     username,
		PasswordHash: string(passwordHash),
	}
	err = h.UserService.Create(user)
	if err != nil {
		log.Print("Create user error: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(200, gin.H{"id": user.ID})
}

// List all users
func (h *UserHandler) List(c *gin.Context) {
	users, err := h.UserService.List()
	if err != nil {
		log.Print("List users error: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, users)
}

// Show an user
func (h *UserHandler) Show(c *gin.Context) {
	username := middleware.JWTGetCurrentUser(c)
	user, err := h.UserService.FindUserByUsername(username)
	if err != nil {
		log.Print("Show user error: ", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, user)
}
