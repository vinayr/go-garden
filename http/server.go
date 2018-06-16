package http

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/vinayr/go-garden/handlers"
	"github.com/vinayr/go-garden/middleware"
	"github.com/vinayr/go-garden/services"
)

// Server represents a HTTP server
type Server struct {
	// Server options
	Addr      string
	JwtSecret string

	// Services
	UserService *services.UserService
}

// NewServer returns a new instance of Server
func NewServer() *Server {
	return &Server{}
}

// Open the server
func (s *Server) Open() error {
	// Start HTTP server
	go http.ListenAndServe(s.Addr, s.router())

	return nil
}

// Close ...
func (s *Server) Close() error {
	return nil
}

func (s *Server) router() http.Handler {
	r := gin.Default()

	// TODO disable in production
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Authorization", "Content-Type"}
	r.Use(cors.New(config))
	// r.Use(cors.Default())

	user := s.userHandler()
	authMiddleware := middleware.Auth(user.UserService, s.JwtSecret)

	r.POST("/signup", user.Signup)
	r.POST("/signin", authMiddleware.LoginHandler)

	admin := r.Group("/admin", authMiddleware.MiddlewareFunc(), middleware.Admin())
	admin.GET("/users", user.List)
	admin.GET("/users/:id", user.Show)

	// auth := r.Group("/")
	// auth.Use(authMiddleware.MiddlewareFunc())
	// {
	// 	auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	// 	auth.GET("/users/me", user.Show)
	// }
	auth := r.Group("/", authMiddleware.MiddlewareFunc())
	auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	auth.GET("/profile", user.Profile)

	return r
}

func (s *Server) userHandler() *handlers.UserHandler {
	h := handlers.NewUserHandler()
	h.UserService = s.UserService
	return h
}
