package main

import (
	"fmt"
	"time"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

// One connection to database

type MsgType int

const (
	Start   MsgType = 0
	Switch          = 1
	Default         = 2
)

func main() {
	createTables()
	fillCityInfoMap()
	initWCache()
	fmt.Println(&weatherCache)

	bot, err := tgbotapi.NewBotAPI(TELEGRAM_TOKEN)
	if err != nil {
		panic(err)
	}

	bot.Debug = false

	// Create a new UpdateConfig struct with an offset of 0. Offsets are used
	// to make sure Telegram knows we've handled previous values and we don't
	// need them repeated.
	updateConfig := tgbotapi.NewUpdate(0)

	// Tell Telegram we should wait up to 30 seconds on each request for an
	// update. This way we can get information just as quickly as making many
	// frequent requests without having to send nearly as many.
	updateConfig.Timeout = 30

	// Start polling Telegram for updates.
	updates, _ := bot.GetUpdatesChan(updateConfig)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			updateWeatherCache()
		case update := <-updates:
			if update.Message == nil {

				msg := createAnswer(update)
				if _, err := bot.Send(msg); err != nil {
					panic(err)
				}
			}
		}
	}
}
