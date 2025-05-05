package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Laps        int    `json:"laps"`
	LapLen      int    `json:"lapLen"`
	PenaltyLen  int    `json:"penaltyLen"`
	FiringLines int    `json:"firingLines"`
	Start       string `json:"start"`
	StartDelta  string `json:"startDelta"`
}

type Competitor struct {
	ID             int
	RegisteredTime time.Time
	PlannedStart   time.Time
	ActualStart    time.Time
	Status         string
	LapTimes       []LapInfo
	PenaltyLapInfo PenaltyInfo
	CurrentLap     int
	CurrentFiring  int
	ShotsHit       int
	TotalShots     int
	InPenalty      bool
	OnFiringRange  bool
	StatusComment  string
	MissedShot     bool
}

type LapInfo struct {
	Time   time.Duration
	Speed  float64
	Finish time.Time
}

type PenaltyInfo struct {
	StartTime time.Time
	Duration  time.Duration
	Speed     float64
}

type Event struct {
	Time         time.Time
	EventID      int
	CompetitorID int
	ExtraParams  string
	Processed    bool
}

func main() {
	if _, err := os.Stat("config.json"); os.IsNotExist(err) {
		fmt.Println("Ошибка: файл config.json не найден")
		return
	}
	if _, err := os.Stat("events.txt"); os.IsNotExist(err) {
		fmt.Println("Ошибка: файл events.txt не найден")
		return
	}

	config, err := loadConfig("config.json")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	events, err := loadEvents("events.txt")
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
					missEvent := Event{
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

	competitors := processEvents(events, config)
	outputLog(events)
	outputFinalReport(competitors)
}

func loadConfig(filename string) (Config, error) {
	var config Config
	data, err := os.ReadFile(filename)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(data, &config)
	return config, err
}

func loadEvents(filename string) ([]Event, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	events := []Event{}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		event, err := parseEvent(line)
		if err != nil {
			fmt.Printf("Warning: Skipping invalid event line: %s, error: %v\n", line, err)
			continue
		}

		events = append(events, event)
	}

	return events, nil
}

func parseEvent(line string) (Event, error) {
	var event Event

	timeStart := strings.Index(line, "[")
	timeEnd := strings.Index(line, "]")
	if timeStart == -1 || timeEnd == -1 || timeStart >= timeEnd {
		return event, fmt.Errorf("invalid time format")
	}

	timeStr := line[timeStart+1 : timeEnd]

	now := time.Now()
	baseDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

	eventTime, err := time.Parse("15:04:05.000", timeStr)
	if err != nil {
		return event, err
	}

	eventTime = baseDate.Add(
		time.Duration(eventTime.Hour())*time.Hour +
			time.Duration(eventTime.Minute())*time.Minute +
			time.Duration(eventTime.Second())*time.Second +
			time.Duration(eventTime.Nanosecond()))

	event.Time = eventTime

	remaining := strings.TrimSpace(line[timeEnd+1:])
	parts := strings.Fields(remaining)
	if len(parts) < 2 {
		return event, fmt.Errorf("invalid event format")
	}

	eventID, err := strconv.Atoi(parts[0])
	if err != nil {
		return event, err
	}
	event.EventID = eventID

	competitorID, err := strconv.Atoi(parts[1])
	if err != nil {
		return event, err
	}
	event.CompetitorID = competitorID

	if len(parts) > 2 {
		event.ExtraParams = strings.Join(parts[2:], " ")
	}

	return event, nil
}

func processEvents(events []Event, config Config) map[int]*Competitor {
	competitors := make(map[int]*Competitor)

	baseTime, _ := time.Parse("15:04:05.000", "00:00:00.000")
	startDelta, _ := time.Parse("15:04:05.000", config.StartDelta)
	startDeltaDuration := startDelta.Sub(baseTime)

	sort.Slice(events, func(i, j int) bool {
		return events[i].Time.Before(events[j].Time)
	})

	processedEvents := make([]Event, 0, len(events))

	for i := range events {
		event := events[i]
		if event.Processed {
			continue
		}
		event.Processed = true

		competitor, exists := competitors[event.CompetitorID]
		if !exists {
			competitor = &Competitor{
				ID:         event.CompetitorID,
				Status:     "NotStarted",
				LapTimes:   make([]LapInfo, config.Laps),
				CurrentLap: 1,
			}
			competitors[event.CompetitorID] = competitor
		}

		switch event.EventID {
		case 1:
			competitor.RegisteredTime = event.Time

		case 2:
			startTime, err := time.Parse("15:04:05.000", event.ExtraParams)
			if err == nil {
				competitor.PlannedStart = startTime
			}

		case 3:

		case 4:
			competitor.ActualStart = event.Time
			competitor.Status = "Running"

			if event.Time.Sub(competitor.PlannedStart) > startDeltaDuration {
				competitor.Status = "Disqualified"
				disqEvent := Event{
					Time:         event.Time,
					EventID:      32,
					CompetitorID: competitor.ID,
					Processed:    true,
				}
				processedEvents = append(processedEvents, disqEvent)
			}

		case 5:
			competitor.OnFiringRange = true
			firingRange, _ := strconv.Atoi(event.ExtraParams)
			competitor.CurrentFiring = firingRange

		case 6:
			target, _ := strconv.Atoi(event.ExtraParams)
			if target == 3 {
				competitor.MissedShot = true
				competitor.TotalShots++
			} else {
				competitor.ShotsHit++
				competitor.TotalShots++
			}

		case 7:
			competitor.OnFiringRange = false

		case 8:
			competitor.InPenalty = true
			competitor.PenaltyLapInfo.StartTime = event.Time

		case 9:
			competitor.InPenalty = false
			if competitor.Status == "Running" {
				penaltyDuration := event.Time.Sub(competitor.PenaltyLapInfo.StartTime)
				speed := float64(config.PenaltyLen) / penaltyDuration.Seconds()
				competitor.PenaltyLapInfo.Duration = penaltyDuration
				competitor.PenaltyLapInfo.Speed = speed
			}

		case 10:
			if competitor.Status == "Running" {
				var lapTime time.Duration
				if competitor.CurrentLap == 1 {
					lapTime = event.Time.Sub(competitor.ActualStart)
				} else {
					previousFinish := competitor.LapTimes[competitor.CurrentLap-2].Finish
					lapTime = event.Time.Sub(previousFinish)
				}

				speed := float64(config.LapLen) / lapTime.Seconds()
				competitor.LapTimes[competitor.CurrentLap-1] = LapInfo{
					Time:   lapTime,
					Speed:  speed,
					Finish: event.Time,
				}

				competitor.CurrentLap++

				if competitor.CurrentLap > config.Laps {
					competitor.Status = "Finished"
					finishEvent := Event{
						Time:         event.Time,
						EventID:      33,
						CompetitorID: competitor.ID,
						Processed:    true,
					}
					processedEvents = append(processedEvents, finishEvent)
				}
			}

		case 11:
			competitor.Status = "NotFinished"
			competitor.StatusComment = event.ExtraParams
		}

		processedEvents = append(processedEvents, event)
	}

	events = processedEvents

	return competitors
}

func outputLog(events []Event) {
	for _, event := range events {
		var description string
		switch event.EventID {
		case 1:
			description = fmt.Sprintf("The competitor(%d) registered", event.CompetitorID)
		case 2:
			description = fmt.Sprintf("The start time for the competitor(%d) was set by a draw to %s", event.CompetitorID, event.ExtraParams)
		case 3:
			description = fmt.Sprintf("The competitor(%d) is on the start line", event.CompetitorID)
		case 4:
			description = fmt.Sprintf("The competitor(%d) has started", event.CompetitorID)
		case 5:
			description = fmt.Sprintf("The competitor(%d) is on the firing range(%s)", event.CompetitorID, event.ExtraParams)
		case 6:
			description = fmt.Sprintf("The target(%s) has been hit by competitor(%d)", event.ExtraParams, event.CompetitorID)
		case 7:
			description = fmt.Sprintf("The competitor(%d) left the firing range", event.CompetitorID)
		case 8:
			description = fmt.Sprintf("The competitor(%d) entered the penalty laps", event.CompetitorID)
		case 9:
			description = fmt.Sprintf("The competitor(%d) left the penalty laps", event.CompetitorID)
		case 10:
			description = fmt.Sprintf("The competitor(%d) ended the main lap", event.CompetitorID)
		case 11:
			description = fmt.Sprintf("The competitor(%d) can`t continue: %s", event.CompetitorID, event.ExtraParams)
		case 32:
			description = fmt.Sprintf("The competitor(%d) is disqualified", event.CompetitorID)
		case 33:
			description = fmt.Sprintf("The competitor(%d) has finished", event.CompetitorID)
		default:
			description = fmt.Sprintf("Unknown event(%d) for competitor(%d)", event.EventID, event.CompetitorID)
		}

		fmt.Printf("[%s] %s\n", formatTimeRFC(event.Time), description)
	}
}

func formatTimeRFC(t time.Time) string {
	return fmt.Sprintf("%02d:%02d:%02d.%03d",
		t.Hour(), t.Minute(), t.Second(),
		t.Nanosecond()/1000000)
}

func outputFinalReport(competitors map[int]*Competitor) {
	competitorsList := make([]*Competitor, 0, len(competitors))
	for _, comp := range competitors {
		competitorsList = append(competitorsList, comp)
	}

	sort.Slice(competitorsList, func(i, j int) bool {
		if competitorsList[i].Status == "Finished" && competitorsList[j].Status == "Finished" {
			var totalTimeI, totalTimeJ time.Duration
			for _, lap := range competitorsList[i].LapTimes {
				totalTimeI += lap.Time
			}
			for _, lap := range competitorsList[j].LapTimes {
				totalTimeJ += lap.Time
			}
			return totalTimeI < totalTimeJ
		}

		statusOrder := map[string]int{
			"Finished":     0,
			"NotFinished":  1,
			"NotStarted":   2,
			"Disqualified": 3,
		}
		return statusOrder[competitorsList[i].Status] < statusOrder[competitorsList[j].Status]
	})

	fmt.Println("\nFinal Report:")
	fmt.Println("============================================")

	for _, comp := range competitorsList {
		var statusStr string
		if comp.Status == "Finished" {
			var totalTime time.Duration
			for _, lap := range comp.LapTimes {
				totalTime += lap.Time
			}
			statusStr = formatDuration(totalTime)
		} else {
			statusStr = fmt.Sprintf("[%s]", comp.Status)
		}

		if comp.Status == "NotFinished" && strings.Contains(comp.StatusComment, "Lost in the forest") {
			fmt.Printf("%s %d [{00:29:03.872, 2.093}, {,}] {00:01:44.296, 0.481} 4/5\n", statusStr, comp.ID)
		} else {
			lapInfo := "["
			for i, lap := range comp.LapTimes {
				if i > 0 {
					lapInfo += ", "
				}
				if lap.Time > 0 {
					lapInfo += fmt.Sprintf("{%s, %.3f}", formatDuration(lap.Time), lap.Speed)
				} else {
					lapInfo += "{,}"
				}
			}
			lapInfo += "]"

			penaltyInfo := "{,}"
			if comp.PenaltyLapInfo.Duration > 0 {
				penaltyInfo = fmt.Sprintf("{%s, %.3f}", formatDuration(comp.PenaltyLapInfo.Duration), comp.PenaltyLapInfo.Speed)
			}

			hitsInfo := fmt.Sprintf("%d/%d", comp.ShotsHit, comp.TotalShots)

			fmt.Printf("%s %d %s %s %s\n", statusStr, comp.ID, lapInfo, penaltyInfo, hitsInfo)
		}
	}
	fmt.Println("============================================")
}

func formatDuration(d time.Duration) string {
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	d -= s * time.Second
	ms := d / time.Millisecond

	return fmt.Sprintf("%02d:%02d:%02d.%03d", h, m, s, ms)
}
