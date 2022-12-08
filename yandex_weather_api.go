package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type SearchResults struct {
	Now      int             `json:"now"`
	Now_dt   string          `json:"now_dt"`
	Info     json.RawMessage `json:"info"`
	Fact     json.RawMessage `json:"fact"`
	Forecast json.RawMessage `json:"forecast"`
}

type WeatherInfo struct {
	Obs_time    int     `json:"obs_time"`
	Temp        int     `json:"temp"`
	Feels_like  int     `json:"feels_like"`
	Icon        string  `json:"icon"`
	Condition   string  `json:"condition"`
	Wind_speed  float32 `json:"wind_speed"`
	Wind_dir    string  `json:"wind_dir"`
	Pressure_mm int     `json:"pressure_mm"`
	Pressure_pa int     `json:"pressure_pa"`
	Humidity    int     `json:"humidity"`
	Daytime     string  `json:"daytime"`
	Polar       bool    `json:"polar"`
	Season      string  `json:"season"`
	Wind_gust   float32 `json:"wind_gust"`
}

func getWeather(coord Coordinates) WValues {
	client := &http.Client{}
	req_str := fmt.Sprintf("https://api.weather.yandex.ru/v2/informers?lat=%v&lon=%v&[lang=ru_RU]", coord.lat, coord.lon)
	fmt.Printf("request = %v", req_str)
	req, err := http.NewRequest("GET", req_str, nil)
	if err != nil {
		fmt.Printf("Bad Link\n")
	}
	req.Header.Add("X-Yandex-API-Key", "44616025-9b90-4b28-8b90-90afef470b2f")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Bad request\n")
	}
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var sr SearchResults
	err = json.Unmarshal(contents, &sr)
	if err != nil {
		panic(err)
	}
	var w_info WeatherInfo
	err = json.Unmarshal(sr.Fact, &w_info)
	if err != nil {
		panic(err)
	}

	out := WValues{}
	out.temp = w_info.Temp
	out.feelsLike = w_info.Feels_like
	out.windSpeed = w_info.Wind_speed
	return out
}
