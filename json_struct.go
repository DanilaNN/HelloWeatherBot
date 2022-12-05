package main

import "encoding/json"

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
