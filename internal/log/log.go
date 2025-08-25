package log

import "go.uber.org/zap"

// New creates a new logger instance
// In Go, functions that create things are often called "New"
func New() *zap.Logger {
	// *zap.Logger means "pointer to a zap.Logger"
	// Pointers let us share objects efficiently
	l, _ := zap.NewDevelopment() // readable for dev; swap to NewProduction in prod
	// The _ means "ignore this return value" (it would be an error)
	return l
}
