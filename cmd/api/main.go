package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/krizad/go-gin-api/config"
	"github.com/krizad/go-gin-api/db"
	"github.com/krizad/go-gin-api/internal/adapters/http/handlers"
	"github.com/krizad/go-gin-api/internal/adapters/http/middleware"
	"github.com/krizad/go-gin-api/internal/adapters/http/routes"
	"github.com/krizad/go-gin-api/internal/adapters/repositories"
	"github.com/krizad/go-gin-api/internal/core/services"
	// "github.com/krizad/go-gin-api/docs"
)

// @title			Go Gin Example API
// @version		1.0
// @description	CRUD API built with Go and Gin framework, connected to PostgreSQL.
// @termsOfService	http://swagger.io/terms/

// @contact.name	API Support
// @contact.url	http://www.example.com/support
// @contact.email	support@example.com

// @license.name	MIT
// @license.url	https://opensource.org/licenses/MIT

// @BasePath	/api/v1

func main() {
	migrateCmd := flag.String("migrate", "", "run migration command: up, down")
	flag.Parse()

	cfg := config.Load()

	db.Connect(cfg)
	defer db.Close()

	if *migrateCmd != "" {
		handleMigrateCmd(*migrateCmd)
		return
	}

	if cfg.AutoMigrate {
		if err := db.RunMigrations(); err != nil {
			log.Fatalf("migration failed: %v", err)
		}
	} else {
		log.Println("auto-migrate disabled (set AUTO_MIGRATE=true to enable)")
	}

	r := gin.New()
	r.Use(middleware.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())

	// Example handler setup
	repo := repositories.NewExampleRepository(db.DB)
	svc := services.NewExampleService(repo)
	exampleHandler := handlers.NewExampleHandler(svc)
	transacRepo := repositories.NewTransactionRepository(db.DB)

	// Account handler setup
	accountRepo := repositories.NewAccountRepository(db.DB)
	transacSvc := services.NewTransactionService(accountRepo, transacRepo)
	// transactionRepo := repositories.NewTransactionRepository(db.DB)
	accountSvc := services.NewAccountService(accountRepo, transacRepo)
	transacHandler := handlers.NewTransactionHandler(transacSvc)
	accountHandler := handlers.NewAccountHandler(accountSvc)

	routes.Setup(r, exampleHandler, accountHandler, transacHandler)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		addr := fmt.Sprintf(":%s", cfg.AppPort)
		log.Printf("server starting on %s", addr)
		if err := r.Run(addr); err != nil {
			log.Fatalf("server failed: %v", err)
		}
	}()

	<-quit
	log.Println("shutting down server...")
}

func handleMigrateCmd(cmd string) {
	switch cmd {
	case "up":
		if err := db.RunMigrations(); err != nil {
			log.Fatalf("migration up failed: %v", err)
		}
		log.Println("migrations applied successfully")
	case "down":
		if err := db.RollbackMigrations(); err != nil {
			log.Fatalf("migration down failed: %v", err)
		}
		log.Println("migrations rolled back successfully")
	default:
		log.Fatalf("unknown migrate command: %s (use 'up' or 'down')", cmd)
	}
}
