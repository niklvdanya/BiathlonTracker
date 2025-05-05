package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/niklvdanya/BiathlonTracker/internal/config"
	"github.com/niklvdanya/BiathlonTracker/internal/event"
	"github.com/niklvdanya/BiathlonTracker/internal/model"
	"github.com/niklvdanya/BiathlonTracker/internal/report"
)

type BiathlonService struct {
	Parser    event.EventParser
	Processor event.EventProcessor
	Reporter  report.Reporter
	Config    config.Config
}

func NewBiathlonService(parser event.EventParser, processor event.EventProcessor, reporter report.Reporter, cfg config.Config) *BiathlonService {
	return &BiathlonService{
		Parser:    parser,
		Processor: processor,
		Reporter:  reporter,
		Config:    cfg,
	}
}

func main() {
	configFileFlag := flag.String("config", "config.json", "Path to configuration file")
	eventsFileFlag := flag.String("events", "events.txt", "Path to events file")
	parallelFlag := flag.Bool("parallel", false, "Use parallel processing")
	flag.Parse()

	if !checkFiles(*configFileFlag, *eventsFileFlag) {
		return
	}

	parser := &event.DefaultEventParser{}
	processor := &event.DefaultEventProcessor{}
	reporter := &report.DefaultReporter{}

	cfg, err := config.Load(*configFileFlag)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	service := NewBiathlonService(parser, processor, reporter, cfg)

	events, err := service.Parser.Parse(*eventsFileFlag)
	if err != nil {
		fmt.Printf("Error loading events: %v\n", err)
		return
	}

	events = handleLostEvents(events)

	ctx := context.Background()
	var competitors map[int]*model.Competitor

	if *parallelFlag {
		competitors = event.ProcessEventsParallel(ctx, events, cfg)
	} else {
		competitors = service.Processor.Process(ctx, events, cfg)
	}

	service.Reporter.OutputLog(events)
	service.Reporter.OutputFinalReport(competitors, cfg)
}

func checkFiles(configFile, eventsFile string) bool {
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
		if events[i].EventID == model.EventLostInForest && strings.Contains(events[i].ExtraParams, model.LostInForestText) {
			shotCount, hasEventID6 := countShotsForCompetitor(events, events[i].CompetitorID)

			if shotCount < 5 && !hasEventID6 {
				missEvent := model.Event{
					Time:         events[i].Time.Add(-time.Second * 5),
					EventID:      model.EventShot,
					CompetitorID: events[i].CompetitorID,
					ExtraParams:  model.ShotTarget3,
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
			if events[j].ExtraParams == model.ShotTarget3 {
				hasEventID6 = true
			}
		}
	}
	return shotCount, hasEventID6
}
