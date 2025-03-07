package logger

import (
	"context"
	"log"
	"time"

	"gorm.io/gorm/logger"
)

// CustomLogger struct
type CustomLogger struct {
	logger.Interface
}

// LogMode sets the logger level
func (c *CustomLogger) LogMode(level logger.LogLevel) logger.Interface {
	return c
}

// Info logs general messages
func (c *CustomLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	log.Printf("[INFO] "+msg, data...)
}

// Warn logs warnings
func (c *CustomLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	log.Printf("[WARN] "+msg, data...)
}

// Error logs errors
func (c *CustomLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	log.Printf("[ERROR] "+msg, data...)
}

// Trace logs all SQL queries with time taken
func (c *CustomLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, rows := fc()
	elapsed := time.Since(begin)

	if err != nil {
		log.Printf("[ERROR] %s | Duration: %s | Rows: %d | Error: %v\n", sql, elapsed, rows, err)
	} else {
		log.Printf("[SQL] %s | Duration: %s | Rows: %d\n", sql, elapsed, rows)
	}
}
