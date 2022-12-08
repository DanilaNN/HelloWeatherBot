package main

import (
	"sync"
)

const MaxWeatherRequests = 50
const MaxCityCnt = MaxWeatherRequests

type WCache struct {
	cities map[string]WValues
}

var (
	weatherCache WCache
	once         sync.Once
)

func NewWCache() WCache {
	once.Do(func() {
		weatherCache.cities = make(map[string]WValues)
		cityNames, err := getCityNames()
		if err != nil {
			panic("Can not read List of cities from Database")
		}
		for _, val := range cityNames {
			weatherCache.cities[val] = WValues{}
		}
	})

	return weatherCache
}
