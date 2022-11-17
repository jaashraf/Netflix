package main

import (
	"database/sql"
	"errors"
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
					db.Exec("insert into listed_in (listed_in_name, show_id) values (?, ?)", listedIn, temp.showId)
				}
				count++
			}
		}
		log.Default().Print("Data sync from excel to database sucessful. ", count, " synced successfully")
	}
	closeDBConnection()
}

func filterByTypeAndCount(movieType string, count int) ([]NetflixData, error) {
	openDBConnection()
	resNetflix, errNetflix := db.Query("select * from NetflixData where movie_type=? limit ?", movieType, count)
	if errNetflix == nil {
		netflixDataArray := make([]NetflixData, 0, 10)
		for resNetflix.Next() {
			var netflixData NetflixData
			resNetflix.Scan(&netflixData.showId, &netflixData.movieType,
				&netflixData.title, &netflixData.dateAdded, &netflixData.releaseYear,
				&netflixData.rating, &netflixData.duration, &netflixData.description)

			res, _ := db.Query("select cast_name from cast where show_id = ?", netflixData.showId)
			for res.Next() {
				var castName string
				res.Scan(&castName)
				netflixData.cast = append(netflixData.cast, castName)
			}
			res, _ = db.Query("select director_name from director where show_id = ?", netflixData.showId)
			for res.Next() {
				var directorName string
				res.Scan(&directorName)
				netflixData.director = append(netflixData.director, directorName)
			}
			res, _ = db.Query("select country_name from country where show_id = ?", netflixData.showId)
			for res.Next() {
				var countryName string
				res.Scan(&countryName)
				netflixData.country = append(netflixData.country, countryName)
			}
			res, _ = db.Query("select listed_in_name from listed_in where show_id = ?", netflixData.showId)
			for res.Next() {
				var listedIn string
				res.Scan(&listedIn)
				netflixData.listedIn = append(netflixData.listedIn, listedIn)
			}
			netflixDataArray = append(netflixDataArray, netflixData)
		}
		return netflixDataArray, nil
	}
	return nil, errors.New("Cannot fetch data by Type")
}

func filterByTypeAndMovieType(movieType string) []NetflixData {
	openDBConnection()
	res, err := db.Query("SELECT n.show_id, n.movie_type, n.title, n.date_added, n.release_year, n.rating, n.duration, n.description from Netflix.NetflixData n join Netflix.listed_in l on n.show_id=l.show_id where n.movie_type like \"%TV%\" and l.listed_in_name like \"%Horror%\"")
	var netflixData = make([]NetflixData, 0)
	fmt.Println(err)
	for res.Next() {
		var temp NetflixData
		res.Scan(&temp.showId, &temp.movieType, &temp.title, &temp.dateAdded, &temp.releaseYear, &temp.rating, &temp.duration, &temp.description)
		dbListedInRes, _ := db.Query("select listed_in_name from listed_in where show_id=?", temp.showId)
		dbCastRes, _ := db.Query("select cast_name from cast where show_id=?", temp.showId)
		dbCountryRes, _ := db.Query("select country_name from country where  show_id=?", temp.showId)
		dbdirectorRes, _ := db.Query("select director_name from director where show_id=?", temp.showId)
		flag := true
		for flag {
			if dbListedInRes.Next() {
				var listedIn string
				dbListedInRes.Scan(&listedIn)
				temp.listedIn = append(temp.listedIn, listedIn)
			}
			if dbCastRes.Next() {
				var cast string
				dbCastRes.Scan(&cast)
				temp.cast = append(temp.cast, cast)
			}
			if dbCountryRes.Next() {
				var country string
				dbCountryRes.Scan(&country)
				temp.country = append(temp.country, country)
			}
			if dbdirectorRes.Next() {
				var director string
				dbdirectorRes.Scan(&director)
				temp.director = append(temp.director, director)
			}

			if !dbListedInRes.Next() && !dbdirectorRes.Next() && !dbCastRes.Next() && !dbCountryRes.Next() {
				flag = false
			}
		}
		netflixData = append(netflixData, temp)
	}
	closeDBConnection()
	return netflixData
}
