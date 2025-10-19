package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"todo-service/internal/handler"
	"todo-service/internal/infrastructure/migration"
	"todo-service/internal/infrastructure/repository"
	"todo-service/internal/infrastructure/storage"
	"todo-service/internal/infrastructure/stream"
	"todo-service/internal/usecase"

	"github.com/rs/zerolog/log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func cmd() {

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open database")

	}
	defer db.Close()

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Retry connection with exponential backoff
	if err := waitForDatabase(db, 30*time.Second); err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	log.Info().Msg("Database connected successfully")

	migrationDSN := fmt.Sprintf(
		"mysql://%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	migrationsPath := "./migrations"
	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		migrationsPath = "/app/migrations"
	}

	log.Info().Msgf("migrationsPath: %+v", migrationsPath)

	if err := migration.RunMigrations(migrationDSN, migrationsPath); err != nil {
		log.Fatal().Err(err).Msg("Failed to run migrations")
	}
	log.Info().Msg("Migrations completed successfully")

	// Initialize repositories
	todoRepo := repository.NewTodoRepository(db)
	fileRepo := repository.NewFileRepository(db)

	// Get S3 configuration
	s3Bucket := os.Getenv("S3_BUCKET")
	if s3Bucket == "" {
		s3Bucket = "todos-bucket" // default
	}

	s3Endpoint := os.Getenv("S3_ENDPOINT")
	if s3Bucket == "" || s3Endpoint == "" {
		log.Fatal().Msg("S3_BUCKET or S3_ENDPOINT not set")
	}

	log.Info().Msgf("S3_BUCKET=%s, S3_ENDPOINT=%s", s3Bucket, s3Endpoint)

	s3Repo, err := storage.NewS3Repository(context.Background(), s3Bucket, s3Endpoint)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize S3")
	}

	redisRepo := stream.NewRedisStreamRepository(os.Getenv("REDIS_ADDR"), os.Getenv("REDIS_STREAM"))

	// Initialize use cases
	todoUC := usecase.NewTodoUseCase(todoRepo, s3Repo, redisRepo, os.Getenv("S3_BUCKET"))
	fileUC := usecase.NewFileUseCase(fileRepo, s3Repo, s3Bucket)

	// Initialize handlers
	h := handler.NewHandler(todoUC, fileUC)

	// Setup routes
	r := mux.NewRouter()
	handler.RegisterRoutes(r, h)

	APP_PORT := os.Getenv("APP_PORT")

	if APP_PORT == "" {
		log.Fatal().Msg("APP_PORT not set")
	}

	srv := &http.Server{Addr: fmt.Sprintf(":%s", APP_PORT), Handler: r}

	go func() {
		log.Info().Msgf("Server starting on :%s", APP_PORT)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Server error")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server shutdown error")
	}
	log.Info().Msg("Server shutdown successfully")
}

// waitForDatabase retries database connection with exponential backoff
func waitForDatabase(db *sql.DB, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	var lastErr error

	for time.Now().Before(deadline) {
		if err := db.Ping(); err == nil {
			return nil
		} else {
			lastErr = err
			time.Sleep(2 * time.Second)
		}
	}

	return fmt.Errorf("timeout waiting for database: %w", lastErr)
}
