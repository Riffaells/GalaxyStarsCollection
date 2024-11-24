package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken          string
	SessionIDs        []string
	TelegramID        int64
	SendToTelegram    bool
	ToleranceFrom     int
	ToleranceTo       int
	StatsPerRequest   int
	StarsAutoBuy      bool
	StarsAutoBuyCount int
	GalaxyIDs         []string
}

func GetEnvAsSlice(key string, delimiter string) []string {
	value := os.Getenv(key)
	if value == "" {
		return []string{}
	}

	return strings.Split(value, delimiter)
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables.")
	}

	sessionIDs := GetEnvAsSlice("SESSION_ID", ",")
	galaxyIDs := GetEnvAsSlice("GALAXY_ID", ",")

	if len(sessionIDs) == 0 {
		log.Fatalf("SESSION_ID is required and must contain at least one session ID")
	}
	if len(galaxyIDs) == 0 {
		log.Fatalf("GALAXY_ID is required and must contain at least one galaxy ID")
	}

	return &Config{
		BotToken:          os.Getenv("BOT_TOKEN"),
		SessionIDs:        sessionIDs,
		TelegramID:        getEnvAsInt64("TELEGRAM_ID"),
		SendToTelegram:    getEnvAsBool("SEND_TO_TELEGRAM"),
		ToleranceFrom:     getEnvAsInt("TOLERANCE_FROM"),
		ToleranceTo:       getEnvAsInt("TOLERANCE_TO"),
		StatsPerRequest:   getEnvAsInt("STATS_PER_REQUEST"),
		StarsAutoBuy:      getEnvAsBool("STARS_AUTO_BUY"),
		StarsAutoBuyCount: getEnvAsInt("STARS_AUTO_BUY_COUNT"),
		GalaxyIDs:         galaxyIDs,
	}, nil
}

func getEnvAsInt(key string) int {
	value := os.Getenv(key)
	if value == "" {
		return 0
	}
	v, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("Invalid value for %s: must be an integer. Error: %v", key, err)
	}
	return v
}

func getEnvAsInt64(key string) int64 {
	value := os.Getenv(key)
	if value == "" {
		return 0
	}
	v, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		log.Fatalf("Invalid value for %s: must be a 64-bit integer. Error: %v", key, err)
	}
	return v
}

func getEnvAsBool(key string) bool {
	value := os.Getenv(key)
	if value == "" {
		return false
	}
	v, err := strconv.ParseBool(value)
	if err != nil {
		log.Fatalf("Invalid value for %s: must be a boolean. Error: %v", key, err)
	}
	return v
}
