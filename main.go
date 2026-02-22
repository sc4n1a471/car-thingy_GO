package main

import (
	"Go_Thingy_GO/application"
	"Go_Thingy_GO/controllers"
	"log/slog"
	"time"
)

// Schedules a daily cleanup of old query inspections
func scheduleDeleteOldQueryInspections() {
	for {
		now := time.Now()
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())

		duration := next.Sub(now)
		time.Sleep(duration)
		controllers.DeleteOldQueryInspections()
	}
}

func main() {
	slog.Info("Starting API...")
	loc, err := time.LoadLocation("Europe/Budapest")
	if err != nil {
		slog.Error(err.Error())
		return
	}
	time.Local = loc
	now := time.Now().Round(0)
	slog.Info("Current time: " + now.GoString() + ", timezone: " + time.Local.String())
	go scheduleDeleteOldQueryInspections()
	application.Api()
}
