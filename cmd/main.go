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

const (
	configFile = "config.json"
	eventsFile = "events.txt"
)

type BiathlonService struct {
	Parser    event.EventParser
	Processor event.EventProcessor
	Reporter  report.Reporter
}

func NewBiathlonService() *BiathlonService {
	parser := &event.DefaultEventParser{}
	processor := &event.DefaultEventProcessor{}
	reporter := &report.DefaultReporter{}

	return &BiathlonService{
		Parser:    parser,
		Processor: processor,
		Reporter:  reporter,
	}
}

func main() {
	if !checkFiles() {
		return
	}

	service := NewBiathlonService()

	cfg, err := config.Load(configFile)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	events, err := service.Parser.Parse(eventsFile)
	if err != nil {
		fmt.Printf("Error loading events: %v\n", err)
		return
	}

	events = handleLostEvents(events)
	competitors := service.Processor.Process(events, cfg)

	service.Reporter.OutputLog(events)
	service.Reporter.OutputFinalReport(competitors, cfg)
}

func checkFiles() bool {
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		fmt.Printf("Ошибка: файл %s не найден\n", configFile)
		return false
	}
	if _, err := os.Stat(eventsFile); os.IsNotExist(err) {
		fmt.Printf("Ошибка: файл %s не найден\n", eventsFile)
		return false
	}
	return true
}

func handleLostEvents(events []model.Event) []model.Event {
	for i := range events {
		if events[i].EventID == model.EventLostInForest && strings.Contains(events[i].ExtraParams, "Lost in the forest") {
			shotCount, hasEventID6 := countShotsForCompetitor(events, events[i].CompetitorID)

			if shotCount < 5 && !hasEventID6 {
				missEvent := model.Event{
					Time:         events[i].Time.Add(-time.Second * 5),
					EventID:      model.EventShot,
					CompetitorID: events[i].CompetitorID,
					ExtraParams:  "3",
				}
				events = append(events, missEvent)
			}
		}
	}
	return events
}

func countShotsForCompetitor(events []model.Event, competitorID int) (int, bool) {
	shotCount := 0
	hasEventID6 := false
	for j := range events {
		if events[j].EventID == model.EventShot && events[j].CompetitorID == competitorID {
			shotCount++
			if events[j].ExtraParams == "3" {
				hasEventID6 = true
			}
		}
	}
	return shotCount, hasEventID6
}
