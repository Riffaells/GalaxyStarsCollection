package main

import (
	"TinyVerse/api"
	"TinyVerse/bot"
	"TinyVerse/config"
	"log"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	apiHandler, err := api.NewAPIHandler(
		"https://api.tonverse.app",
		cfg.SessionIDs,
		map[string]string{
			"Accept":          "*/*",
			"Accept-Encoding": "gzip, deflate, br",
			"Accept-Language": "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7",
			"Connection":      "keep-alive",
			"Content-Type":    "application/x-www-form-urlencoded;charset=UTF-8",
			"Origin":          "https://app.tonverse.app",
			"Referer":         "https://app.tonverse.app/",
			"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
		},
	)
	if err != nil {
		log.Fatalf("Failed to initialize APIHandler: %v", err)
	}

	botConfig := map[string]interface{}{
		"ToleranceFrom":   cfg.ToleranceFrom,
		"ToleranceTo":     cfg.ToleranceTo,
		"StatsPerRequest": cfg.StatsPerRequest,
		"GalaxyID":        cfg.GalaxyIDs,
		"StarsAutoBuy":    cfg.StarsAutoBuy,
		"StarsCount":      cfg.StarsAutoBuyCount,
	}

	botInstance, err := bot.NewBot(apiHandler, cfg.BotToken, cfg.TelegramID, botConfig)
	if err != nil {
		log.Fatalf("Failed to initialize bot: %v", err)
	}

	botInstance.Run()
}
