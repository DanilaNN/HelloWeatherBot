package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

func createAnswer(db *sql.DB) func(update tgbotapi.Update) tgbotapi.MessageConfig {
	var prevMsgType MsgType = Default

	return func(update tgbotapi.Update) tgbotapi.MessageConfig {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		var b strings.Builder
		switch update.Message.Text {
		case "/start":
			b.Reset()
			b.WriteString("Привет! Я HelloWeatherBot,\nВыбери свой город и пришли мне его цифру:\n")
			for cityId, cityInfo := range CityInfoMap {
				b.WriteString(strconv.Itoa(int(cityId)))
				b.WriteString(" - ")
				b.WriteString(cityInfo.name)
				b.WriteString("\n")
			}
			b.WriteString("Чтобы получить прогноз погоды пришлите любой текст\n")
			b.WriteString("Город можно сменить командой /switch\n")
			prevMsgType = Start
		case "/switch":
			b.Reset()
			b.WriteString("Выбери новый город и пришли мне его цифру:\n")
			for cityId, cityInfo := range CityInfoMap {
				b.WriteString(strconv.Itoa(int(cityId)))
				b.WriteString(" - ")
				b.WriteString(cityInfo.name)
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
					b.WriteString(strconv.Itoa(len(CityInfoMap)))
					prevMsgType = Start
				}
				if id < 1 || id > len(CityInfoMap) {
					b.WriteString("Пришлите целое число от 1 до ")
					b.WriteString(strconv.Itoa(len(CityInfoMap)))
					prevMsgType = Start
				} else {
					b.WriteString("Ваш город: ")
					b.WriteString(CityInfoMap[CityId(id)].name)
					b.WriteString("\n")
					b.WriteString("Чтобы получить прогноз погоды пришлите любой текст\n")
					addNewUser(db, update.Message.Chat.ID, id)
					prevMsgType = Default
				}
			case Switch:
				id, err := strconv.Atoi(update.Message.Text)
				if err != nil {
					b.WriteString("Пришлите целое число от 1 до ")
					b.WriteString(strconv.Itoa(len(CityInfoMap)))
					prevMsgType = Switch
				}

				if id < 1 || id > len(CityInfoMap) {
					b.WriteString("Пришлите целое число от 1 до ")
					b.WriteString(strconv.Itoa(len(CityInfoMap)))
					prevMsgType = Switch
				} else {
					b.WriteString("Ваш город: ")
					b.WriteString(CityInfoMap[CityId(id)].name)
					b.WriteString("\n")
					b.WriteString("Чтобы получить прогноз погоды пришлите любой текст\n")
					fmt.Printf("Update user %v with city %v!\n", update.Message.Chat.ID, id)
					err = switchUserCity(db, update.Message.Chat.ID, id)
					if err != nil {
						fmt.Printf("Can't update user with new city!\n")
					} else {
						fmt.Printf("User was updated with new city!\n")
					}
					prevMsgType = Default
				}
			case Default:
				weatherInfo := getCachedWeather(UserId(update.Message.Chat.ID))

				b.Reset()
				val, ok := CityInfoMap[UserCityMap[UserId(update.Message.Chat.ID)]]
				if !ok {
					b.WriteString("Не получилось получить прогноз:(\n")
				} else {
					b.WriteString("По данным Яндекс.Погоды в городе: ")
					b.WriteString(val.name)
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
		return msg
	}
}
