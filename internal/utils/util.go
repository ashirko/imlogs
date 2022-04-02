package utils

import (
	"os"
	"strconv"
)

func GetEnvStr(key string, defaultVal string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	return defaultVal
}

func GetEnvInt(key string, defaultVal int) int {
	if valStr, exists := os.LookupEnv(key); exists {
		val, err := strconv.Atoi(valStr)
		if err == nil {
			return val
		}
	}
	return defaultVal
}
