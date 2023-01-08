package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var (
	DBCon *sql.DB

	dbInfo = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s", host, port, user, dbname, sslmode)
)

func createTables(db *sql.DB) {
	createTableCity(db)
	createTableUsers(db)
}

func createTableCity(db *sql.DB) error {
	if _, err := db.Exec(`CREATE TABLE cities(ID SERIAL PRIMARY KEY, NAME TEXT, LONGITUDE REAL, LATITUDE REAL);`); err != nil {
		return err
	} else {
		fmt.Printf("Table cities was created")
	}

	insertInitCities(db)

	return nil
}

func createTableUsers(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE users(ID SERIAL PRIMARY KEY, USER_ID BIGINT, CITY_ID INT);`)
	return err
}

func insertInitCities(db *sql.DB) error {
	req := `INSERT INTO cities (name, longitude, latitude) VALUES($1, $2, $3);`
	_, err := db.Exec(req, "Нижний Новгород", 44.002, 56.3287)
	if err != nil {
		return err
	}
	if _, err = db.Exec(req, "Москва", 37.6173, 55.7558); err != nil {
		return err
	}
	if _, err = db.Exec(req, "Санкт-Петербург", 30.3351, 59.9343); err != nil {
		return err
	}
	if _, err = db.Exec(req, "Казань", 49.1233, 55.7879); err != nil {
		return err
	}

	return nil
}

func addNewUser(db *sql.DB, user_id int64, city_id int) error {
	rows, err := db.Query(fmt.Sprintf("SELECT exists (SELECT * FROM users WHERE user_id = %v  LIMIT 1);", user_id))
	if err != nil {
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		var user_exist bool
		if err := rows.Scan(&user_exist); err != nil {
			return err
		}
		if user_exist {
			fmt.Printf("User exist\n")
			return nil
		} else {
			fmt.Printf("User not exist - add new user\n")
			break
		}
	}

	data := `INSERT INTO users (user_id, city_id) VALUES($1, $2);`
	_, err = db.Exec(data, user_id, city_id)
	if err != nil {
		fmt.Printf("Can not insert user: user_id=%v, city_id=%v", user_id, city_id)
		return err
	}

	UserCityMap[UserId(user_id)] = CityId(city_id)
	return nil
}

func switchUserCity(db *sql.DB, user_id int64, new_city_id int) error {
	rows, err := db.Query(fmt.Sprintf("SELECT exists (SELECT * FROM users WHERE user_id = %v  LIMIT 1);", user_id))
	if err != nil {
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		var user_exist bool
		if err := rows.Scan(&user_exist); err != nil {
			return err
		}
		if user_exist == true {
			fmt.Printf("User not exist\n")
			break
		} else {
			fmt.Printf("User exist\n")
			return nil
		}
	}

	data := `UPDATE users SET city_id=$1 WHERE user_id=$2;`
	_, err = db.Exec(data, new_city_id, user_id)
	if err != nil {
		fmt.Printf("Can't update user: user_id=%v, city_id=%v\n", new_city_id, user_id)
		return err
	}

	return nil
}

func getCityIds(db *sql.DB) ([]CityId, error) {
	out := make([]CityId, 0, MaxCityCnt)

	rows, err := db.Query("SELECT id FROM cities")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var city_id CityId
		if err := rows.Scan(&city_id); err != nil {
			return nil, err
		}
		out = append(out, city_id)
	}

	return out, nil
}

func fillMapsFromDb(db *sql.DB) {
	rows, err := db.Query("SELECT id, name, longitude, latitude FROM cities")
	if err != nil {
		panic("fillCityIdNameMap: Error reading cities table")
	}
	defer rows.Close()

	for rows.Next() {
		var id CityId
		var name string
		var lon float64
		var lat float64
		if err := rows.Scan(&id, &name, &lon, &lat); err != nil {
			panic("fillMapsFromDb: Error reading cities table")
		}
		CityInfoMap[id] = CityInfo{name, Coordinates{lon, lat}}
	}

	rows, err = db.Query("SELECT user_id, city_id FROM users")
	if err != nil {
		panic("fillCityIdNameMap: Error reading users table")
	}
	for rows.Next() {
		var user_id UserId
		var city_id CityId
		if err := rows.Scan(&user_id, &city_id); err != nil {
			panic("fillMapsFromDb: Error reading users table")
		}
		UserCityMap[user_id] = city_id
	}
}
