package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/vinayr/go-garden/env"
	"github.com/vinayr/go-garden/http"
	"github.com/vinayr/go-garden/services"
	"github.com/vinayr/go-garden/storage"
)

func main() {
	m := NewMain()

	// Load configuration
	if err := m.Config.Load(); err != nil {
		log.Print("LoadConfig error: ", err)
		os.Exit(1)
	}

	// Execute program.
	if err := m.Run(); err != nil {
		log.Print(err)
		os.Exit(1)
	}

	// Shutdown on SIGINT (CTRL-C).
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Print("Received interrupt, shutting down...")
	m.Close()
}

// Main represents the main program execution
type Main struct {
	Config  *env.Config
	closeFn func() error
}

// NewMain returns a new instance of Main
func NewMain() *Main {
	return &Main{
		Config:  env.NewConfig(),
		closeFn: func() error { return nil },
	}
}

// Close cleans up the program
func (m *Main) Close() error { return m.closeFn() }

// Run executes the program
func (m *Main) Run() error {
	// Open database
	db := storage.NewDB()
	db.Path = fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=disable",
		m.Config.DbHost,
		m.Config.DbName,
		m.Config.DbUser,
		m.Config.DbPassword,
	)
	if err := db.Open(); err != nil {
		return err
	}
	log.Print("Database initialized: ", db.Path)

	// Database migration
	db.Migrate()

	// Instantiate db services
	userService := services.NewUserService(db.DB)

	// Initialize HTTP server
	httpServer := http.NewServer()
	httpServer.Addr = m.Config.HttpAddr
	httpServer.JwtSecret = m.Config.JwtSecret
	httpServer.UserService = userService

	// Open HTTP server
	if err := httpServer.Open(); err != nil {
		return err
	}
	log.Print("HTTP listening: ", httpServer.Addr)

	// Assign close function
	m.closeFn = func() error {
		log.Print("Closing...")
		httpServer.Close()
		db.Close()
		return nil
	}

	return nil
}
