package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/niklvdanya/BiathlonTracker/internal/config"
	"github.com/niklvdanya/BiathlonTracker/internal/event"
	"github.com/niklvdanya/BiathlonTracker/internal/model"
	"github.com/niklvdanya/BiathlonTracker/internal/report"
)

func main() {
	if _, err := os.Stat("config.json"); os.IsNotExist(err) {
		fmt.Println("Ошибка: файл config.json не найден")
		return
	}
	if _, err := os.Stat("events.txt"); os.IsNotExist(err) {
		fmt.Println("Ошибка: файл events.txt не найден")
		return
	}

	cfg, err := config.Load("config.json")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	events, err := event.LoadEvents("events.txt")
	if err != nil {
		fmt.Printf("Error loading events: %v\n", err)
		return
	}

	for i := range events {
		if events[i].EventID == 11 && strings.Contains(events[i].ExtraParams, "Lost in the forest") {
			shotCount := 0
			hasEventID6 := false
			for j := range events {
				if events[j].EventID == 6 && events[j].CompetitorID == events[i].CompetitorID {
					shotCount++
					if events[j].ExtraParams == "3" {
						hasEventID6 = true
					}
				}
			}

			if shotCount < 5 {
				if !hasEventID6 {
					missEvent := model.Event{
						Time:         events[i].Time.Add(-time.Second * 5),
						EventID:      6,
						CompetitorID: events[i].CompetitorID,
						ExtraParams:  "3",
					}
					events = append(events, missEvent)
				}
			}
		}
	}

	competitors := event.ProcessEvents(events, cfg)
	report.OutputLog(events)
	report.OutputFinalReport(competitors, cfg)
}
