package main

import (
	"github.com/robfig/cron"
	"log"
)

func excelToDBSyncScheduler() {
	c := cron.New()
	c.AddFunc("@every 1m", func() {
		log.Default().Println("Running scheduler")
		errGlobal = nil
		netflixDataGlobal, errGlobal = readCSVToObject(filepath)
		if errGlobal != nil {
			log.Default().Println("Cannot fetch data from Netflix CSV")
		} else {
			syncDB(netflixDataGlobal[1:])
		}
	})
	c.Start()
}
