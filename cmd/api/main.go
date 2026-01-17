// cmd/api/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/invoice-app-be/config"
	"github.com/invoice-app-be/internal/domain/invoice"
	"github.com/invoice-app-be/internal/domain/timeentry"
	"github.com/invoice-app-be/internal/domain/user"
	"github.com/invoice-app-be/internal/infrastructure/auth"
	"github.com/invoice-app-be/internal/infrastructure/database/postgres"
	"github.com/invoice-app-be/internal/infrastructure/integrations/jira"
	"github.com/invoice-app-be/internal/infrastructure/integrations/square"
	"github.com/invoice-app-be/internal/infrastructure/pdf"
	infraHTTP "github.com/invoice-app-be/internal/interfaces/http"
	"github.com/invoice-app-be/internal/interfaces/http/handlers"
	"github.com/invoice-app-be/internal/interfaces/http/middleware"
	"github.com/invoice-app-be/internal/pkg/logger"
)

func main() {
	// Load configuration
	cfg, err := config.Load()

	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Setup logger
	appLogger := logger.New("debug")
	logLevel := slog.LevelInfo
	if cfg.Server.Environment == "development" {
		logLevel = slog.LevelDebug
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))
	slog.SetDefault(logger)

	// Connect to database
	logger.Info("Connecting to database...")

	db, err := postgres.NewConnection(&cfg.Database)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Run migrations - UPDATED
	logger.Info("Running database migrations...")
	if err := runMigrations(&cfg.Database); err != nil {
		logger.Error("Failed to run migrations", "error", err)
		// Don't exit - migrations might already be applied
		logger.Warn("Continuing without migrations...")
	}

	// Initialize repositories
	invoiceRepo := postgres.NewInvoiceRepository(db)
	timeEntryRepo := postgres.NewTimeEntryRepository(db)
	userRepo := postgres.NewUserRepository(db)
	clientRepo := postgres.NewClientRepository(db)

	// Initialize integrations (optional, based on config)
	var jiraClient *jira.Client
	if cfg.Jira.Enabled && cfg.Jira.BaseURL != "" {
		jiraClient = jira.NewClient(cfg.Jira.BaseURL, "", "") // User-specific credentials loaded per request
		logger.Info("Jira integration enabled")
	}

	var squareClient *square.Client
	if cfg.Square.Enabled && cfg.Square.AccessToken != "" {
		squareClient = square.NewClient(cfg.Square.AccessToken, cfg.Square.Environment)
		logger.Info("Square integration enabled", "environment", cfg.Square.Environment)
	}

	// Initialize PDF generator
	pdfGenerator := pdf.NewGenerator()

	// Initialize services
	invoiceService := invoice.NewService(invoiceRepo, pdfGenerator, squareClient)
	timeEntryService := timeentry.NewService(timeEntryRepo, jiraClient)
	userService := user.NewService(userRepo, cfg.Auth.JWTSecret, appLogger)

	// Initialize auth components
	jwtManager := auth.NewJWTManager(cfg.Auth.JWTSecret, cfg.Auth.TokenDuration)
	authMiddleware := middleware.NewAuthMiddleware(jwtManager)

	// Initialize HTTP handlers
	authHandler := handlers.NewAuthHandler(userService, jwtManager)
	invoiceHandler := handlers.NewInvoiceHandler(invoiceService, clientRepo)
	timeEntryHandler := handlers.NewTimeEntryHandler(timeEntryService)

	// Setup router
	router := infraHTTP.NewRouter(
		invoiceHandler,
		timeEntryHandler,
		authHandler,
		authMiddleware,
	)
	handler := router.Setup()

	// Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      handler,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in goroutine
	go func() {
		logger.Info("Starting server",
			"port", cfg.Server.Port,
			"environment", cfg.Server.Environment)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server error", "error", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
	}

	logger.Info("Server exited")
}

// Add this helper function at the bottom of main.go
func runMigrations(cfg *config.DatabaseConfig) error {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode,
	)

	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		return fmt.Errorf("creating migration instance: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("running migrations: %w", err)
	}

	return nil
}
