package event

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/niklvdanya/BiathlonTracker/internal/model"
)

func LoadEvents(filename string) ([]model.Event, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	events := []model.Event{}

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

func parseEvent(line string) (model.Event, error) {
	var event model.Event

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
