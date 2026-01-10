package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"vendix/internal/config"
	"vendix/internal/database"
	"vendix/internal/jobs"
	"vendix/internal/logger"

	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize logger
	logger.InitLogger()
	defer logger.Sync()

	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		logger.Fatal("Failed to connect to database", "error", err)
	}
	defer db.Close()

	logger.Info("Worker database connection established")

	// Create Redis connection options
	redisOpt := asynq.RedisClientOpt{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	}

	// Initialize task handlers
	handlers := jobs.NewHandlers(db, cfg)

	// Create asynq server
	srv := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	// Register task handlers
	mux := asynq.NewServeMux()
	mux.HandleFunc(jobs.TaskSendEmail, handlers.HandleSendEmail)
	mux.HandleFunc(jobs.TaskSignInvoice, handlers.HandleSignInvoice)
	mux.HandleFunc(jobs.TaskSendInvoiceToDGII, handlers.HandleSendInvoiceToDGII)
	mux.HandleFunc(jobs.TaskProcessRecurringBilling, handlers.HandleProcessRecurringBilling)
	mux.HandleFunc(jobs.TaskGenerateMonthlyReports, handlers.HandleGenerateMonthlyReports)
	mux.HandleFunc(jobs.TaskCheckPurchaseDueDates, handlers.HandleCheckPurchaseDueDates)
	mux.HandleFunc(jobs.TaskSchedulePurchaseDueChecks, handlers.HandleSchedulePurchaseDueChecks)

	// Create scheduler for periodic tasks
	scheduler := asynq.NewScheduler(
		redisOpt,
		&asynq.SchedulerOpts{
			LogLevel: asynq.InfoLevel,
		},
	)

	// Schedule daily purchase due date checks at 9:00 AM
	// This will create jobs for all active tenants
	_, err = scheduler.Register(
		"0 9 * * *", // Cron expression: Every day at 9:00 AM
		asynq.NewTask(
			jobs.TaskSchedulePurchaseDueChecks,
			nil,
		),
		asynq.Queue("default"),
	)
	if err != nil {
		logger.Error("Failed to register purchase due dates scheduler", "error", err)
	} else {
		logger.Info("Scheduled daily purchase due date checks at 9:00 AM")
	}

	logger.Info("Starting worker server and scheduler")

	// Start scheduler in a goroutine
	go func() {
		if err := scheduler.Run(); err != nil {
			logger.Fatal("Scheduler failed", "error", err)
		}
	}()

	// Start worker in a goroutine
	go func() {
		if err := srv.Run(mux); err != nil {
			logger.Fatal("Worker server failed", "error", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down worker and scheduler...")
	scheduler.Shutdown()
	srv.Shutdown()

	logger.Info("Worker exited gracefully")
}
