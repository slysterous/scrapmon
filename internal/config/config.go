package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	scrapmon "github.com/slysterous/scrapmon/internal/scrapmon"
)

// FromEnv returns the apps configuration based on environmental variables including sane defaults.
func FromEnv() scrapmon.Config {
	return scrapmon.Config{
		DatabaseHost:       getString("SCRAPMON_DB_HOST", "127.0.0.1"),
		DatabaseName:       getString("SCRAPMON_DB_NAME", "scrapmon"),
		DatabasePort:       getString("SCRAPMON_DB_PORT", "5432"),
		DatabaseUser:       getString("SCRAPMON_DB_USER", "postgres"),
		DatabasePassword:   getString("SCRAPMON_DB_PASSWORD", "password"),
		MaxDBConnections:   getInt("MAX_DB_CONNECTIONS", 100),
		Env:                getString("SCRAPMON_ENV", "dev"),
		TorHost:            getString("TOR_HOST", "127.0.0.1"),
		TorPort:            getString("TOR_PORT", "9050"),
		ScrapStorageFolder: getString("SCRAPMON_IMAGE_FOLDER", "./"),
	}
}

// getString returns the string value of an env variable.
func getString(key, fallback string) string {
	env := os.Getenv(key)
	if env == "" {
		log.Printf("debug: missing env variable for key: %s, using default: %s", key, fallback)
		return fallback
	}
	return env
}

// getString returns the converted int value from a string env variable.
func getInt(key string, fallback int) int {
	strValue := getString(key, fmt.Sprint(fallback))

	intValue, err := strconv.Atoi(strValue)
	if err != nil {
		log.Printf("debug: converting env %s to int", key)
		return fallback
	}
	return intValue
}
