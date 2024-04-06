package main

import (
	"go.uber.org/zap"
)

func main() {
	// Create a logger
	logger, _ := zap.NewProduction()

	// Log messages
	logger.Info("This is an info message", zap.String("key", "value"))
}
