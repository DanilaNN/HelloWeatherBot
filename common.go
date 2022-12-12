package main

// Weather data
type WValues struct {
	temp      int
	feelsLike int
	windSpeed float32
}

// Coordinates of cities
type Coordinates struct {
	lon float64
	lat float64
}

type CityInfo struct {
	name  string
	coord Coordinates
}

type CityId int
type UserId int64

var CityInfoMap = make(map[CityId]CityInfo)
var UserCityMap = make(map[UserId]CityId)
