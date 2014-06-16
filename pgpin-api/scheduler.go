package main

import (
	"fmt"
	"github.com/jrallison/go-workers"
	"log"
	"time"
)

func SchedulerEnqueue(pin *Pin) error {
	log.Printf("scheduler.enqueue pin_id=%s", pin.Id)
	pin.ScheduledAt = time.Now()
	err := DataPinUpdate(pin)
	if err != nil {
		return err
	}
	err = workers.Enqueue("pins", "", pin.Id)
	if err != nil {
		return err
	}
	return nil
}

func SchedulerTick() error {
	log.Printf("scheduler.tick")
	queryFrag := fmt.Sprintf("scheduled_at <= now()-'%f seconds'::interval", ConfigPinRefreshInterval.Seconds())
	ready, err := DataPinList(queryFrag)
	if err != nil {
		return err
	} else {
		for _, pin := range ready {
			err = SchedulerEnqueue(pin)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func SchedulerStart() {
	log.Printf("scheduler.start")
	DataStart()
	QueueStart()
	for {
		err := SchedulerTick()
		if err != nil {
			log.Printf("scheduler.error %+s", err.Error())
		}
		time.Sleep(ConfigSchedulerTickInterval)
	}
}