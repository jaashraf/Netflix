package main

import (
	"encoding/csv"
	"os"
	"strconv"
	"strings"
	"time"
)

func csvToNetflixDataObject(line []string) NetflixData {

	relYear, err := strconv.Atoi(line[7])

	if err != nil {
		relYear = -1
	}
	data := NetflixData{
		showId:      line[0],
		movieType:   line[1],
		title:       line[2],
		director:    strings.Split(line[3], ","),
		cast:        strings.Split(line[4], ","),
		country:     strings.Split(line[5], ","),
		dateAdded:   line[6],
		releaseYear: relYear,
		rating:      line[8],
		duration:    line[9],
		listedIn:    strings.Split(line[10], ","),
		description: line[11],
	}
	return data
}

func readCSVToObject(filename string) ([]NetflixData, error) {
	content, e := os.Open(filename)
	lines, err := csv.NewReader(content).ReadAll()

	if e != nil {
		return nil, e
	} else if err != nil {
		return nil, err
	}

	netlfixDataSlice := []NetflixData{}
	for _, line := range lines {
		netlfixDataSlice = append(netlfixDataSlice, csvToNetflixDataObject(line))
	}

	return netlfixDataSlice, nil
}

func filterBYType(movieType string, netflixDataArray []NetflixData) []NetflixData {
	resultNetflixData := []NetflixData{}
	for _, data := range netflixDataArray {
		if strings.Contains(strings.ToLower(data.movieType), strings.ToLower(movieType)) {
			resultNetflixData = append(resultNetflixData, data)
		}
	}
	return resultNetflixData
}

func filterByListedIn(listedIn string, netflixDataArray []NetflixData) []NetflixData {
	resultNetflixData := []NetflixData{}
	for _, data := range netflixDataArray {
		for _, listedData := range data.listedIn {
			if strings.Contains(strings.ToLower(listedData), strings.ToLower(listedIn)) {
				resultNetflixData = append(resultNetflixData, data)
			}
		}
	}
	return resultNetflixData
}

func filterByCountry(countryName string, netflixDataArray []NetflixData) []NetflixData {
	resultNetflixData := []NetflixData{}
	for _, data := range netflixDataArray {
		for _, countryData := range data.country {
			if strings.Contains(strings.ToLower(countryData), strings.ToLower(countryName)) {
				resultNetflixData = append(resultNetflixData, data)
			}
		}
	}
	return resultNetflixData
}

func filterByAddedDate(startDate string, endDate string, netflixDataArray []NetflixData) ([]NetflixData, error) {
	layout := "January 02, 2006"
	sTime, er1 := time.Parse(layout, startDate)
	eTime, er2 := time.Parse(layout, endDate)
	resultNetflixData := []NetflixData{}

	if er1 != nil {
		return nil, er1
	}
	if er2 != nil {
		return nil, er2
	}

	for _, data := range netflixDataArray {
		if data.dateAdded == "date_added" {
			continue
		} else {
			curTime, er3 := time.Parse(layout, data.dateAdded)
			if er3 != nil {
				return nil, er3
			}
			if sTime.UnixMilli() <= curTime.UnixMilli() && curTime.UnixMilli() <= eTime.UnixMilli() {
				resultNetflixData = append(resultNetflixData, data)
			}
		}

	}
	return resultNetflixData, nil
}
