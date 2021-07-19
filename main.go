package main

import (
	"bufio"
	"etd-transaction/mock"
	"github.com/robfig/cron/v3"
	"log"
	"os"
	"time"
)

func getDay() int64 {
	delta := time.Now().Sub(time.Date(2021, 7, 19, 0, 0, 0, 0, time.Local))
	day := int64(delta.Hours() / 24)
	return day
}

func start() {
	var tsk mock.Task

	day := getDay()
	err := tsk.InitTask(day)
	if err != nil {
		log.Fatal(err.Error())
	}

	ticker := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-ticker.C:
			if getDay() != tsk.Day {
				log.Printf("[day %d] ETD: %v\nTransaction: %v", tsk.Day, tsk.ETD, tsk.Trans)
				return
			}
			go tsk.Run()
		}
	}
}

func main() {
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer file.Close()
	buf := bufio.NewWriter(file)

	log.SetOutput(buf)
	go start()
	cd := cron.New()
	_, err = cd.AddFunc("0 0 */1 * *", start)
	if err != nil {
		log.Fatal(err.Error())
	}
	cd.Start()
	select {}
}
