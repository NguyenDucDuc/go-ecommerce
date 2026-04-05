package util

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Println("Load env không thành công")
	}
	log.Println("Load env thành công")
}

func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func GetIntEnv(key string, defaultValue int) int {
	valueStr, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	// Chuyển từ string sang int
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		// Log lỗi nếu cần: log.Printf("Invalid int value for %s: %v", key, err)
		return defaultValue
	}

	return value
}