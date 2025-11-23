package pkg

import (
	"fmt"
	"log"
	"os"
	"taskflow/internal/common"
)

func GetEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	// If fallback is empty, treat it as a required variable.
	if fallback == "" {
		// Fatal → exit with a clear error message
		log.Fatal(
			common.ErrorResponse{
				Message: fmt.Sprintf("missing required environment variable: %s", key),
			}.Error(),
		)
	}

	return fallback
}
