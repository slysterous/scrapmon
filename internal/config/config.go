package config

import (
	"log"
	"os"
	"strconv"
	printscrape "github.com/slysterous/print-scrape/internal/domain"
)


// FromEnv returns the apps configuration based on environmental variables including sane defaults.
func FromEnv() printscrape.Config {
	return printscrape.Config{
		DatabaseHost:      getString("WAVE_ACCOUNTS_DB_HOST", "127.0.0.1"),
		DatabaseName:      getString("WAVE_ACCOUNTS_DB_NAME", "wave-accounts"),
		DatabasePort:      getString("WAVE_ACCOUNTS_DB_PORT", "5432"),
		DatabaseUser:      getString("WAVE_ACCOUNTS_DB_USER", "postgres"),
		DatabasePassword:  getString("WAVE_ACCOUNTS_DB_PASSWORD", "password"),
		MaxDBConnections:  getInt("MAX_DB_CONNECTIONS", 100),
		Env:               getString("WAVE_ACCOUNTS_ENV", "dev"),
		HTTPClientTimeout: getInt("HTTP_CLIENT_TIMEOUT_SECONDS", 15),
	}
}

// getString returns the string value of an env variable.
func getString(key, fallback string) string {
	env := os.Getenv(key)
	if env == "" {
		log.Printf("debug: missing env variable for key: %s",key)
		return fallback
	}
	return env
}

// getString returns the converted int value from a string env variable.
func getInt(key string,fallback int) int  {
	strValue := getString(key,string(fallback))

	intValue, err := strconv.Atoi(strValue)
	if err !=nil {
		log.Printf("debug: converting env %s to int",key)
		return fallback
	}
	return intValue
}