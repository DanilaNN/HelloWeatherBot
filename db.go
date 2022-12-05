package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var dbInfo = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s", host, port, user, dbname, sslmode)

func createTables() {
	createTableCity()
	createTableUsers()
}

func createTableCity() error {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		fmt.Printf("Can not connect to db\n")
		return err
	}
	defer db.Close()

	if _, err = db.Exec(`CREATE TABLE cities(ID SERIAL PRIMARY KEY, NAME TEXT, LONGITUDE REAL, LATITUDE REAL);`); err != nil {
		return err
	} else {
		fmt.Printf("Table cities was created")
	}

	insertInitCities(db)

	return nil
}

func createTableUsers() error {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		fmt.Printf("Can not connect to db\n")
		return err
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE users(ID SERIAL PRIMARY KEY, USER_ID BIGINT, CITY_ID INT);`)
	return err
}

func insertInitCities(db *sql.DB) error {
	data := `INSERT INTO cities (name, longitude, latitude) VALUES($1, $2, $3);`
	_, err := db.Exec(data, "Нижний Новгород", 44.002, 56.3287)
	if err != nil {
		return err
	}
	if _, err = db.Exec(data, "Москва", 37.6173, 55.7558); err != nil {
		return err
	}
	if _, err = db.Exec(data, "Санкт-Петербург", 30.3351, 59.9343); err != nil {
		return err
	}
	if _, err = db.Exec(data, "Казань", 49.1233, 55.7879); err != nil {
		return err
	}

	return nil
}

func addNewUser(user_id int64, city_id int) error {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	defer db.Close()

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
			fmt.Printf("User exist\n")
			return nil
		} else {
			fmt.Printf("User not exist\n")
			break
		}
	}

	data := `INSERT INTO users (user_id, city_id) VALUES($1, $2);`
	_, err = db.Exec(data, user_id, city_id)
	if err != nil {
		fmt.Printf("Can not insert user: user_id=%v, city_id=%v", user_id, city_id)
		return err
	}

	return nil
}

func switchUserCity(user_id int64, new_city_id int) error {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	defer db.Close()

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

func getCityCoordinates(chat_id int64) (Coordinates, error) {
	var coord Coordinates
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return Coordinates{}, err
	}
	defer db.Close()

	row := db.QueryRow(fmt.Sprintf("SELECT longitude, latitude FROM cities WHERE id=(SELECT city_id FROM users WHERE user_id=%v);", chat_id))
	err = row.Scan(&coord.lat, &coord.lon)
	if err != nil {
		return Coordinates{}, err
	}

	return coord, nil
}

func getCityName(chat_id int64) (string, error) {
	var name string
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return "", err
	}
	defer db.Close()

	row := db.QueryRow(fmt.Sprintf("SELECT name FROM cities WHERE id=(SELECT city_id FROM users WHERE user_id=%v);", chat_id))
	err = row.Scan(&name)
	if err != nil {
		return "", err
	}

	return name, nil
}

func getCities() (map[int]string, error) {
	out := make(map[int]string)
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, name FROM cities")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var city string
		if err := rows.Scan(&id, &city); err != nil {
			return nil, err
		}
		out[id] = city
	}

	return out, nil
}
