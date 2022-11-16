package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

var db *sql.DB = nil
var dbErr error = nil

func openDBConnection() {
	db, dbErr = sql.Open("mysql", os.Getenv("USER_NAME")+":"+os.Getenv("PASSWORD")+"@tcp("+os.Getenv("SERVER_NAME")+")/"+os.Getenv("DB_NAME"))
	if db.Ping() != nil && dbErr != nil {
		log.Default().Println("DB Connection not successful")
	} else {
		log.Default().Println("DB Connection established successfully")
	}
}

func closeDBConnection() {
	err := db.Close()
	if err != nil {
		log.Default().Println("Not able to close the db connection")
	} else {
		log.Default().Println("DB closed successfully")
	}
}

func syncDB(netflixData []NetflixData) {

	openDBConnection()

	dbRes, _ := db.Query("select count(*) from NetflixData")
	var count int
	dbRes.Next()
	dbRes.Scan(&count)

	if count < len(netflixData) {

		dbRes2, err := db.Query("select show_id from NetflixData")
		showMap := make(map[string]int)

		if err != nil {
			i := 0
			var showId string
			for dbRes2.Next() {
				dbRes2.Scan(&showId)
				showMap[showId] = 0
				i++
			}
		}
		count := 0
		for _, temp := range netflixData {
			_, flag := showMap[temp.showId]
			if flag {
				continue
			} else {
				db.Exec("insert into NetflixData (show_id, movie_type, title, date_added, release_year, rating, duration, description) values (?, ?, ?, ?, ?, ?, ?, ?) ",
					temp.showId, temp.movieType, temp.title, temp.dateAdded, temp.releaseYear, temp.rating, temp.duration, temp.description)
				for _, cast := range temp.cast {
					db.Exec("insert into cast (cast_name, show_id) values (?, ?)",
						cast, temp.showId)

				}
				for _, country := range temp.country {
					db.Exec("insert into country (country_name, show_id) values (?, ?)", country, temp.showId)
				}
				for _, director := range temp.director {
					res, err := db.Exec("insert into director (director_name, show_id) values (?, ?);", director, temp.showId)
					fmt.Println(res, " ----", err)
				}
				for _, listedIn := range temp.listedIn {
					db.Exec("insert into listed_in (listed_in_name, show_id) values (?, ?);", listedIn, temp.showId)
				}
				count++
			}
		}
		log.Default().Print("Data sync from excel to database sucessful. ", count, " synced successfully")
	}
	closeDBConnection()
}
