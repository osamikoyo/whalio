package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"whalio/config"
	"whalio/core"
	"whalio/handlers"
	"whalio/models"
	"whalio/repository"
	"whalio/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v2"
	"github.com/rs/zerolog"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// Setup logger
	logger := setupLogger(cfg)

	db, err := gorm.Open(sqlite.Open(cfg.DatabasePath))
	if err != nil {
		logger.Error().Msgf("Failed open database: %v", err)

		return
	}

	if err = db.AutoMigrate(&models.Song{}, &models.Artist{}, &models.Album{}); err != nil {
		logger.Error().Msgf("Failed make migration: %v", err)

		return
	}

	logger.Info().Msgf("Successfully connected to db: %s", cfg.DatabasePath)

	// Ensure required directories exist
	if err := os.MkdirAll(cfg.UploadDir, 0o755); err != nil {
		logger.Fatal().Err(err).Msgf("Failed to create upload dir: %s", cfg.UploadDir)
	}
	if err := os.MkdirAll(cfg.ImageDir, 0o755); err != nil {
		logger.Fatal().Err(err).Msgf("Failed to create image dir: %s", cfg.ImageDir)
	}
	if err := os.MkdirAll(cfg.StaticDir, 0o755); err != nil {
		logger.Fatal().Err(err).Msgf("Failed to create static dir: %s", cfg.StaticDir)
	}

	core := core.NewCore(repository.NewRepository(&logger, db), storage.NewStorage(&logger), cfg, 30*time.Second)
	// Create router
	r := chi.NewRouter()

	// Setup middleware
	setupMiddleware(r, cfg)

	// Initialize handlers
	h := handlers.New(core)

	// Register routes
	h.RegisterRoutes(r)

	// Create server
	srv := &http.Server{
		Addr:         cfg.Address(),
		Handler:      r,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	// Start server in goroutine
	go func() {
		logger.Info().Msgf("🐋 Whalio server starting on %s", cfg.Address())
		logger.Info().Msgf("🌍 Environment: %s", cfg.Environment)
		logger.Info().Msgf("🎯 Open: http://%s", cfg.Address())

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info().Msg("🛑 Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	logger.Info().Msg("✅ Server exited")
}

// setupLogger configures the application logger
func setupLogger(cfg *config.Config) zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339

	var logger zerolog.Logger

	if cfg.LogFormat == "json" {
		logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	} else {
		logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
	}

	// Set log level
	switch cfg.LogLevel {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	return logger
}

// setupMiddleware configures all middleware
func setupMiddleware(r *chi.Mux, cfg *config.Config) {
	// Request ID middleware
	r.Use(middleware.RequestID)

	// Real IP middleware
	r.Use(middleware.RealIP)

	// Logging middleware
	r.Use(httplog.RequestLogger(httplog.NewLogger("whalio", httplog.Options{
		JSON:             cfg.LogFormat == "json",
		Concise:          true,
		RequestHeaders:   cfg.Debug,
		MessageFieldName: "message",
		Tags: map[string]string{
			"version": "1.0.0",
			"env":     cfg.Environment,
		},
	})))

	// Recoverer middleware
	r.Use(middleware.Recoverer)

	// Timeout middleware
	r.Use(middleware.Timeout(30 * time.Second))

	// Compress middleware
	r.Use(middleware.Compress(5))

	// CORS middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.AllowedOrigins,
		AllowedMethods:   cfg.AllowedMethods,
		AllowedHeaders:   cfg.AllowedHeaders,
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Security headers
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("X-XSS-Protection", "1; mode=block")

			if cfg.IsProduction() {
				w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
			}

			next.ServeHTTP(w, r)
		})
	})

	// Rate limiting (if enabled)
	if cfg.RateLimitEnabled {
		r.Use(middleware.Throttle(cfg.RateLimit))
	}

	// Development middleware
	if cfg.IsDevelopment() {
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
				next.ServeHTTP(w, r)
			})
		})
	}
}
