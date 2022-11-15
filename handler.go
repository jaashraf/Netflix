package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

var netflixDataGlobal = []NetflixData{}
var errGlobal error = nil

func getTVShowsApiHandler(w http.ResponseWriter, r *http.Request) {
	exeStartTime := time.Now()
	if r.URL.Query().Has("n") {
		netflixDataGlobal, errGlobal = getTVShowsAsPerCountHandler(r)
	} else if r.URL.Query().Has("movieType") {
		netflixDataGlobal, errGlobal = getTVShowsByMovieTypeHandler(r)
	} else if r.URL.Query().Has("country") {
		netflixDataGlobal, errGlobal = getTVShowsByCountryHandler(r)
	} else if r.URL.Query().Has("startDate") && r.URL.Query().Has("endDate") {
		netflixDataGlobal, errGlobal = getTVShowsBetweenDatesHandler(r)
	} else {
		json.NewEncoder(w).Encode("Invalid Input")
	}
	exeEndTime := time.Since(exeStartTime)
	fmt.Println(exeEndTime)
	w.Header().Set("X-Time-Execute", exeEndTime.String())
	w.Header().Set("Content-Type", "application/json")
	if errGlobal != nil {
		json.NewEncoder(w).Encode(errGlobal)
	} else {
		json.NewEncoder(w).Encode(netflixDataGlobal)
	}
}

func getTVShowsAsPerCountHandler(r *http.Request) ([]NetflixData, error) {
	count, _ := strconv.ParseInt(r.URL.Query().Get("n"), 10, 64)
	netflixData, err := readCSVToObject(filepath)
	if err == nil {
		netflixData := filterBYType("TV Show", netflixData)[0:count]
		log.Default().Print("TV Shows as per the given count ", count, " : \n\n\n", netflixData)
		if netflixData != nil {
			return netflixData, nil
		}
	}
	return nil, errors.New("cannot fetch TV shows by count")
}

func getTVShowsByMovieTypeHandler(r *http.Request) ([]NetflixData, error) {
	movieType := r.URL.Query().Get("movieType")
	netflixData, err := readCSVToObject(filepath)
	if err == nil {
		netflixData := filterByListedIn(movieType, filterBYType("TV Show", netflixData))
		log.Default().Print("TV Shows of Movie Type ", movieType, " : \n\n\n", netflixData)
		if netflixData != nil {
			return netflixData, nil
		}
	}
	return nil, errors.New("Cannot fetch TV shows by Movie Type")
}

func getTVShowsByCountryHandler(r *http.Request) ([]NetflixData, error) {
	countryName := r.URL.Query().Get("country")
	netflixData, err := readCSVToObject(filepath)
	if err == nil {
		netflixData := filterByCountry(countryName, netflixData)
		log.Default().Print("TV Shows of Country ", countryName, " : \n\n\n", netflixData)
		if netflixData != nil {
			return netflixData, nil
		}
	}
	return nil, errors.New("Cannot fetch TV shows by country")
}

func getTVShowsBetweenDatesHandler(r *http.Request) ([]NetflixData, error) {
	startDate := r.URL.Query().Get("startDate")
	endDate := r.URL.Query().Get("endDate")
	netflixData, err := readCSVToObject(filepath)
	if err == nil {
		netflixData, _ := filterByAddedDate(startDate, endDate, netflixData)
		log.Default().Print("TV Shows between ", startDate, "and ", endDate, " : ", netflixData)
		if netflixData != nil {
			return netflixData, nil
		}
	}
	return nil, errors.New("Cannot fetch TV Shows between given dates")
}
