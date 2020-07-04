package config

import (
	printscrape "github.com/slysterous/print-scrape/internal/domain"
	"log"
	"os"
	"strconv"
)

// FromEnv returns the apps configuration based on environmental variables including sane defaults.
func FromEnv() printscrape.Config {
	return printscrape.Config{
		DatabaseHost:     getString("PRINT_SCRAPE_DB_HOST", "127.0.0.1"),
		DatabaseName:     getString("PRINT_SCRAPE_DB_NAME", "print-scrape"),
		DatabasePort:     getString("PRINT_SCRAPE_DB_PORT", "5432"),
		DatabaseUser:     getString("PRINT_SCRAPE_DB_USER", "postgres"),
		DatabasePassword: getString("PRINT_SCRAPE_DB_PASSWORD", "password"),
		MaxDBConnections: getInt("MAX_DB_CONNECTIONS", 100),
		Env:              getString("PRINT_SCRAPE_ENV", "dev"),
		TorHost:          getString("TOR_HOST", "tor"),
		TorPort:          getString("TOR_PORT", "9050"),
	}
}

// getString returns the string value of an env variable.
func getString(key, fallback string) string {
	env := os.Getenv(key)
	if env == "" {
		log.Printf("debug: missing env variable for key: %s", key)
		return fallback
	}
	return env
}

// getString returns the converted int value from a string env variable.
func getInt(key string, fallback int) int {
	strValue := getString(key, string(fallback))

	intValue, err := strconv.Atoi(strValue)
	if err != nil {
		log.Printf("debug: converting env %s to int", key)
		return fallback
	}
	return intValue
}
