package common

import (
	"log"
)

// Used for logging general information
func LogInfo(message string) {
	log.Println("INFO:", message)
}

// Used for logging errors
func LogError(message string) {
	log.Println("ERROR:", message)
}
