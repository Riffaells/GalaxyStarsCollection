package bot

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"TinyVerse/api"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	ToleranceFrom   int
	ToleranceTo     int
	StatsPerRequest int
	APIHandler      *api.APIHandler
	TelegramBot     *tgbotapi.BotAPI
	TelegramChatID  int64
	GalaxyIDs       []string
	StarsAutoBuy    bool
	StarsCount      int
}

func NewBot(apiHandler *api.APIHandler, telegramToken string, telegramChatID int64, config map[string]interface{}) (*Bot, error) {
	tgBot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		return nil, err
	}

	toleranceFrom, _ := config["ToleranceFrom"].(int)

	toleranceTo, _ := config["ToleranceTo"].(int)

	statsPerRequest, _ := config["StatsPerRequest"].(int)

	galaxyIDs, _ := config["GalaxyIDs"].([]string)
	starsAutoBuy, _ := config["StarsAutoBuy"].(bool)
	starsCount, _ := config["StarsCount"].(int)

	return &Bot{
		ToleranceFrom:   toleranceFrom,
		ToleranceTo:     toleranceTo,
		StatsPerRequest: statsPerRequest,
		APIHandler:      apiHandler,
		TelegramBot:     tgBot,
		TelegramChatID:  telegramChatID,
		GalaxyIDs:       galaxyIDs,
		StarsAutoBuy:    starsAutoBuy,
		StarsCount:      starsCount,
	}, nil
}

func (b *Bot) Run() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	b.notifyTelegram("🤖 Bot started!")

	for range ticker.C {
		b.collectStars()

		b.handleStatistics()

		if b.StarsAutoBuy {
			b.attemptAutoBuy()
		}

		timeToSleep := generateRandomBetweenRange(b.ToleranceFrom, b.ToleranceTo)

		nextExecutionTime := time.Now().Add(time.Duration(timeToSleep) * time.Second)

		log.Printf("Sleeping for %d seconds. Next iteration at: %s", timeToSleep, nextExecutionTime.Format("15:04:05"))

		time.Sleep(time.Duration(timeToSleep) * time.Second)
	}
}

func (b *Bot) collectStars() {
	results, err := b.APIHandler.CollectStars()
	if err != nil {
		log.Printf("Error collecting stars: %v", err)
		b.notifyTelegram(fmt.Sprintf("⚠️ Error collecting stars: %v", err))
		return
	}

	for _, result := range results {

		session, _ := result["session"].(string)

		if errMsg, exists := result["error"]; exists {
			log.Printf("Session %s: error collecting stars: %v", session, errMsg)
			continue
		}

		response, ok := result["response"].(map[string]interface{})
		if !ok {
			log.Printf("Session %s: 'response' field not found or is invalid", session)
			continue
		}

		success, ok := response["success"].(float64)
		if !ok || success != 1 {
			log.Printf("Session %s: failed to collect stars.", session)
			continue
		}

		dust, ok := response["dust"].(float64)
		if !ok {
			log.Printf("Session %s: unexpected format for 'dust'.", session)
			continue
		}

		message := fmt.Sprintf("🌌 Successfully collected %v stardust", dust)
		log.Println(message)
		b.notifyTelegram(message)
	}
}

func (b *Bot) handleStatistics() {
	b.StatsPerRequest--
	if b.StatsPerRequest > 0 {
		return
	}

	results, err := b.APIHandler.CheckStats()
	if err != nil {
		log.Printf("Error checking stats: %v", err)
		b.notifyTelegram(fmt.Sprintf("⚠️ Error checking stats: %v", err))
		return
	}

	for i, result := range results {
		session := result["session"].(string)

		if err, exists := result["error"]; exists {
			log.Printf("Session %s: error checking stats: %v", session, err)
			continue
		}

		stats := formatStats(result)
		log.Printf("SessionId %d stats: %s", i, stats)
		b.notifyTelegram(stats)
	}

	b.StatsPerRequest = 10
}

func (b *Bot) attemptAutoBuy() {
	for _, galaxyID := range b.GalaxyIDs {
		results, err := b.APIHandler.BuyStars(galaxyID, fmt.Sprintf("%d", b.StarsCount))
		if err != nil {
			log.Printf("Error processing BuyStars for galaxy %s: %v", galaxyID, err)
			b.notifyTelegram(fmt.Sprintf("⚠️ Error processing BuyStars for galaxy %s: %v", galaxyID, err))
			continue
		}

		for _, result := range results {
			session := result["session"].(string)

			if errMsg, exists := result["error"]; exists {
				log.Printf("Session %s, Galaxy %s: Error buying stars: %v", session, galaxyID, errMsg)
				b.notifyTelegram(fmt.Sprintf("⚠️ Session %s, Galaxy %s: Error buying stars: %v", session, galaxyID, errMsg))
				continue
			}

			apiResponse, ok := result["response"].(map[string]interface{})
			if !ok {
				log.Printf("Session %s, Galaxy %s: Unexpected response format", session, galaxyID)
				b.notifyTelegram(fmt.Sprintf("⚠️ Session %s, Galaxy %s: Unexpected response format", session, galaxyID))
				continue
			}

			invoice, ok := apiResponse["invoice"].(string)
			if ok {
				message := fmt.Sprintf("🌌 Invoice generated for session %s, galaxy %s: %s", session, galaxyID, invoice)
				log.Println(message)
				b.notifyTelegram(message)
			} else {
				log.Printf("Session %s, Galaxy %s: Invoice not found in response", session, galaxyID)
				b.notifyTelegram(fmt.Sprintf("⚠️ Session %s, Galaxy %s: Invoice not found in response", session, galaxyID))
			}
		}
	}
}

func (b *Bot) notifyTelegram(message string) {
	msg := tgbotapi.NewMessage(b.TelegramChatID, message)
	_, err := b.TelegramBot.Send(msg)
	if err != nil {
		log.Printf("Failed to send message to Telegram: %v", err)
	}
}

func generateRandomBetweenRange(min, max int) int {
	return rand.Intn(max-min+1) + min
}

func formatStats(stats map[string]interface{}) string {
	response := stats["response"].(map[string]interface{})
	return fmt.Sprintf(
		`📊 User Statistics
	- Stardust: %v
	- Stars: %v / %v`,
		response["dust"], response["stars"], response["stars_max"],
	)
}
