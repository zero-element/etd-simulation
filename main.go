package main

import (
	"bufio"
	"fmt"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"github.com/zero-element/etd-transaction/config"
	"github.com/zero-element/etd-transaction/mock"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	ex chan bool
	wg sync.WaitGroup
)

func getDay() int64 {
	delta := time.Now().Sub(config.StartTime)
	day := int64(delta.Hours() / 24)
	return day
}

func start() {
	var tsk mock.Task

	wg.Add(1)
	defer wg.Done()

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
				log.Infof("[day %d] ETD: %v\nTransaction: %v", tsk.Day, tsk.ETD, tsk.Trans)
				return
			}
			go tsk.Run()
		case <-ex:
			log.Infof("[day %d] ETD: %v\nTransaction: %v", tsk.Day, tsk.ETD, tsk.Trans)
			return
		}
	}
}

func main() {
	ex = make(chan bool)
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err.Error())
	}
	buf := bufio.NewWriter(file)
	defer func() {
		err := buf.Flush()
		if err != nil {
			fmt.Println(err.Error())
		}
		err = file.Close()
		if err != nil {
			fmt.Printf(err.Error())
		}
	}()

	log.SetOutput(buf)
	go start()
	cd := cron.New()
	_, err = cd.AddFunc("0 0 */1 * *", start)
	if err != nil {
		log.Fatal(err.Error())
	}
	cd.Start()

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	for s := range c {
		switch s {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			ex <- true
			wg.Wait()
			buf.Flush()
			return
		default:
			fmt.Println("other signal", s)
		}
	}
}
