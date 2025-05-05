package event

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/niklvdanya/BiathlonTracker/internal/model"
	"github.com/niklvdanya/BiathlonTracker/internal/utils"
)

func LoadEvents(filename string) ([]model.Event, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	events := make([]model.Event, 0, len(lines))

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

	timeStr, remaining, err := extractTimeString(line)
	if err != nil {
		return event, err
	}

	eventTime, err := parseTimeString(timeStr)
	if err != nil {
		return event, err
	}

	event.Time = eventTime

	eventID, competitorID, extraParams, err := parseEventDetails(remaining)
	if err != nil {
		return event, err
	}

	event.EventID = eventID
	event.CompetitorID = competitorID
	event.ExtraParams = extraParams

	return event, nil
}

func extractTimeString(line string) (string, string, error) {
	timeStart := strings.Index(line, "[")
	timeEnd := strings.Index(line, "]")
	if timeStart == -1 || timeEnd == -1 || timeStart >= timeEnd {
		return "", "", utils.ErrInvalidTimeFormat
	}

	timeStr := line[timeStart+1 : timeEnd]
	remaining := strings.TrimSpace(line[timeEnd+1:])

	return timeStr, remaining, nil
}

func parseTimeString(timeStr string) (time.Time, error) {
	now := time.Now()
	baseDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

	eventTime, err := time.Parse(model.TimeFormat, timeStr)
	if err != nil {
		return time.Time{}, err
	}

	eventTime = baseDate.Add(
		time.Duration(eventTime.Hour())*time.Hour +
			time.Duration(eventTime.Minute())*time.Minute +
			time.Duration(eventTime.Second())*time.Second +
			time.Duration(eventTime.Nanosecond()))

	return eventTime, nil
}

func parseEventDetails(text string) (int, int, string, error) {
	parts := strings.Fields(text)
	if len(parts) < 2 {
		return 0, 0, "", utils.ErrInvalidEventFormat
	}

	eventID, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, "", err
	}

	competitorID, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, "", err
	}

	var extraParams string
	if len(parts) > 2 {
		extraParams = strings.Join(parts[2:], " ")
	}

	return eventID, competitorID, extraParams, nil
}
