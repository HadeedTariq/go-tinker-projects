package utils

import (
	"fmt"
	"math/rand"
	"time"
)

// Function to generate a random string of a given length
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// Function to simulate logging
func logMessage(isError bool) string {
	var logType string
	if isError {
		logType = "ERROR"
	} else {
		logType = "SUCCESS"
	}
	randomString := generateRandomString(10) // Generate a random string of length 10
	log := fmt.Sprintf("[%s] %s: %s", time.Now().Format(time.RFC3339), logType, randomString)
	return log
}

func GenerateLog() string {

	isError := rand.Intn(2) == 0 // Randomly decide if it's an error log or success log
	return logMessage(isError)
}
func IsErrorLog(log string) bool {
	return len(log) > 0 && log[1:6] == "ERROR" // Check if the log type is ERROR
}
