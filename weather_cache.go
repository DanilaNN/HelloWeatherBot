package main

import (
	"database/sql"
	"fmt"
	"sync"
)

const MaxWeatherRequests = 50
const MaxCityCnt = MaxWeatherRequests

type WCache struct {
	cities map[CityId]WValues
}

var (
	weatherCache WCache
	once         sync.Once
)

func initWCache(db *sql.DB) {
	once.Do(func() {
		weatherCache.cities = make(map[CityId]WValues)

		cityNames, err := getCityIds(db)
		if err != nil {
			panic("Can not read List of cities from Database")
		}
		for _, val := range cityNames {
			weatherCache.cities[val] = WValues{}
		}
		updateWeatherCache()
	})
}

func (wc *WCache) String() string {
	var out string = ""
	for cityId, wInfo := range weatherCache.cities {
		out += fmt.Sprintf("City: %v, temp: %v, fl: %v, wspeed: %v\n",
			CityInfoMap[cityId].name, wInfo.temp, wInfo.feelsLike, wInfo.windSpeed)
	}
	return out
}

func updateWeatherCache() {
	for cityId := range weatherCache.cities {
		weatherCache.cities[cityId] = getWeather(CityInfoMap[cityId].coord)
	}
}

func getCachedWeather(userId UserId) WValues {
	return weatherCache.cities[UserCityMap[userId]]
}
