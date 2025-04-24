package configs

import (
	"os"
	"strconv"
	"strings"
)

func intEnv(key string, fallback int) int {
	value := os.Getenv(key)
	intVal, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return intVal
}

func stringEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func boolEnv(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	if value == "yes" {
		return true
	}
	return false
}

func listEnv(key string, fallback []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return strings.Split(value, ",")
}

func mapEnv(key string, fallback string) map[string]bool {
	value := os.Getenv(key)
	if value == "" {
		return map[string]bool{fallback: true}
	}
	mapper := make(map[string]bool)
	keys := strings.Split(value, ",")
	for _, k := range keys {
		mapper[k] = true
	}
	return mapper
}
