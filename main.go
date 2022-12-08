package main

import (
	"fmt"
	"strconv"
	"strings"

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

	cities, err := getCities()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", cities)

	// bot, err := tgbotapi.NewBotAPI(os.Getenv("5844951877:AAE6bcT2BxNZK5D_NjI5wM6AsLxrSWT9AzA"))
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

	var prevMsgType MsgType = Default
	for update := range updates {
		// Telegram can send many types of updates depending on what your Bot
		// is up to. We only want to look at messages for now, so we can
		// discard any other updates.
		if update.Message == nil {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		var b strings.Builder
		switch update.Message.Text {
		case "/start":
			b.Reset()
			b.WriteString("Привет! Я HelloWeatherBot,\nВыбери свой город и пришли мне его цифру:\n")
			for id, city := range cities {
				b.WriteString(strconv.Itoa(id))
				b.WriteString(" - ")
				b.WriteString(city)
				b.WriteString("\n")
			}
			b.WriteString("Чтобы получить прогноз погоды пришлите любой текст\n")
			b.WriteString("Город можно сменить командой /switch\n")
			prevMsgType = Start
		case "/switch":
			b.Reset()
			b.WriteString("Выбери новый город и пришли мне его цифру:\n")
			for id, city := range cities {
				b.WriteString(strconv.Itoa(id))
				b.WriteString(" - ")
				b.WriteString(city)
				b.WriteString("\n")
			}
			prevMsgType = Switch
		default:
			b.Reset()
			switch prevMsgType {
			case Start:
				id, err := strconv.Atoi(update.Message.Text)
				if err != nil {
					b.WriteString("Пришлите целое число от 1 до ")
					b.WriteString(strconv.Itoa(len(cities)))
					prevMsgType = Start
				}
				if id < 1 || id > len(cities) {
					b.WriteString("Пришлите целое число от 1 до ")
					b.WriteString(strconv.Itoa(len(cities)))
					prevMsgType = Start
				} else {
					b.WriteString("Ваш город: ")
					b.WriteString(cities[id])
					b.WriteString("\n")
					b.WriteString("Чтобы получить прогноз погоды пришлите любой текст\n")
					addNewUser(update.Message.Chat.ID, id)
					prevMsgType = Default
				}
			case Switch:
				id, err := strconv.Atoi(update.Message.Text)
				if err != nil {
					b.WriteString("Пришлите целое число от 1 до ")
					b.WriteString(strconv.Itoa(len(cities)))
					prevMsgType = Switch
				}

				if id < 1 || id > len(cities) {
					b.WriteString("Пришлите целое число от 1 до ")
					b.WriteString(strconv.Itoa(len(cities)))
					prevMsgType = Switch
				} else {
					b.WriteString("Ваш город: ")
					b.WriteString(cities[id])
					b.WriteString("\n")
					b.WriteString("Чтобы получить прогноз погоды пришлите любой текст\n")
					fmt.Printf("Update user %v with city %v!\n", update.Message.Chat.ID, id)
					err = switchUserCity(update.Message.Chat.ID, id)
					if err != nil {
						fmt.Printf("Can't update user with new city!\n")
					} else {
						fmt.Printf("User was updated with new city!\n")
					}
					prevMsgType = Default
				}
			case Default:
				coord, err := getCityCoordinates(update.Message.Chat.ID)
				if err != nil {
					prevMsgType = Default
					fmt.Printf("Can't read coordinates from data base! - %v\n", err.Error())
					continue
				}
				// fmt.Printf("Coordinates: lon=%v, lat=%v\n", coord.lon, coord.lat)

				weatherInfo := getWeather(coord)

				b.Reset()
				city_name, err := getCityName(update.Message.Chat.ID)
				if err != nil {
					b.WriteString("Не получилось получить прогноз:(\n")
				} else {
					b.WriteString("По данным Яндекс.Погоды в городе: ")
					b.WriteString(city_name)
					b.WriteString("\nTемпература: ")
					b.WriteString(strconv.Itoa(weatherInfo.temp))
					b.WriteString(" C\n")
					b.WriteString("Ощущается как: ")
					b.WriteString(strconv.Itoa(weatherInfo.feelsLike))
					b.WriteString(" C\n")
					b.WriteString("Скорость ветра: ")
					b.WriteString(fmt.Sprintf("%.1f", weatherInfo.windSpeed))
					b.WriteString(" м/с\n")
					b.WriteString("\nХорошего дня!")
					prevMsgType = Default
				}
			}
		}

		msg.Text = b.String()
		if _, err := bot.Send(msg); err != nil {
			panic(err)
		}

	}
}
