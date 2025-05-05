package event

import (
	"sort"
	"strconv"
	"time"

	"github.com/niklvdanya/BiathlonTracker/internal/config"
	"github.com/niklvdanya/BiathlonTracker/internal/model"
)

func ProcessEvents(events []model.Event, cfg config.Config) map[int]*model.Competitor {
	competitors := make(map[int]*model.Competitor)

	baseTime, _ := time.Parse("15:04:05.000", "00:00:00.000")
	startDelta, _ := time.Parse("15:04:05.000", cfg.StartDelta)
	startDeltaDuration := startDelta.Sub(baseTime)

	sortEvents(events)
	processedEvents := make([]model.Event, 0, len(events))

	for i := range events {
		event := events[i]
		if event.Processed {
			continue
		}
		event.Processed = true

		competitor := getOrCreateCompetitor(competitors, event.CompetitorID, cfg.Laps)
		processEvent(competitor, event, cfg, startDeltaDuration, &processedEvents)
		processedEvents = append(processedEvents, event)
	}

	for i := range events {
		events[i] = processedEvents[i]
	}

	return competitors
}

func sortEvents(events []model.Event) {
	sort.Slice(events, func(i, j int) bool {
		return events[i].Time.Before(events[j].Time)
	})
}

func getOrCreateCompetitor(competitors map[int]*model.Competitor, competitorID int, laps int) *model.Competitor {
	competitor, exists := competitors[competitorID]
	if !exists {
		competitor = &model.Competitor{
			ID:         competitorID,
			Status:     model.StatusNotStarted,
			LapTimes:   make([]model.LapInfo, laps),
			CurrentLap: 1,
		}
		competitors[competitorID] = competitor
	}
	return competitor
}

func processEvent(competitor *model.Competitor, event model.Event, cfg config.Config, startDeltaDuration time.Duration, processedEvents *[]model.Event) {
	switch event.EventID {
	case model.EventRegistration:
		handleRegistrationEvent(competitor, event)
	case model.EventSetStartTime:
		handleSetStartTimeEvent(competitor, event)
	case model.EventStartLine:
		// start line
	case model.EventStarted:
		handleStartedEvent(competitor, event, startDeltaDuration, processedEvents)
	case model.EventFiringRange:
		handleFiringRangeEvent(competitor, event)
	case model.EventShot:
		handleShotEvent(competitor, event)
	case model.EventLeaveFiring:
		handleLeaveFireEvent(competitor, event)
	case model.EventEnterPenalty:
		handleEnterPenaltyEvent(competitor, event)
	case model.EventLeavePenalty:
		handleLeavePenaltyEvent(competitor, event, cfg)
	case model.EventLapEnd:
		handleLapEndEvent(competitor, event, cfg, processedEvents)
	case model.EventLostInForest:
		handleLostEvent(competitor, event)
	}
}

func handleRegistrationEvent(competitor *model.Competitor, event model.Event) {
	competitor.RegisteredTime = event.Time
}

func handleSetStartTimeEvent(competitor *model.Competitor, event model.Event) {
	startTime, err := time.Parse("15:04:05.000", event.ExtraParams)
	if err == nil {
		competitor.PlannedStart = startTime
	}
}

func handleStartedEvent(competitor *model.Competitor, event model.Event, startDeltaDuration time.Duration, processedEvents *[]model.Event) {
	competitor.ActualStart = event.Time
	competitor.Status = model.StatusRunning

	if event.Time.Sub(competitor.PlannedStart) > startDeltaDuration {
		competitor.Status = model.StatusDisqualified
		disqEvent := model.Event{
			Time:         event.Time,
			EventID:      model.EventDisqualified,
			CompetitorID: competitor.ID,
			Processed:    true,
		}
		*processedEvents = append(*processedEvents, disqEvent)
	}
}

func handleFiringRangeEvent(competitor *model.Competitor, event model.Event) {
	competitor.OnFiringRange = true
	firingRange, _ := strconv.Atoi(event.ExtraParams)
	competitor.CurrentFiring = firingRange
}

func handleShotEvent(competitor *model.Competitor, event model.Event) {
	target, _ := strconv.Atoi(event.ExtraParams)
	if target == 3 {
		competitor.MissedShot = true
		competitor.TotalShots++
	} else {
		competitor.ShotsHit++
		competitor.TotalShots++
	}
}

func handleLeaveFireEvent(competitor *model.Competitor, event model.Event) {
	competitor.OnFiringRange = false
}

func handleEnterPenaltyEvent(competitor *model.Competitor, event model.Event) {
	competitor.InPenalty = true
	competitor.PenaltyLapInfo.StartTime = event.Time
}

func handleLeavePenaltyEvent(competitor *model.Competitor, event model.Event, cfg config.Config) {
	competitor.InPenalty = false
	if competitor.IsRunning() {
		penaltyDuration := event.Time.Sub(competitor.PenaltyLapInfo.StartTime)
		speed := float64(cfg.PenaltyLen) / penaltyDuration.Seconds()
		competitor.PenaltyLapInfo.Duration = penaltyDuration
		competitor.PenaltyLapInfo.Speed = speed
	}
}

func handleLapEndEvent(competitor *model.Competitor, event model.Event, cfg config.Config, processedEvents *[]model.Event) {
	if competitor.IsRunning() {
		var lapTime time.Duration
		if competitor.CurrentLap == 1 {
			lapTime = event.Time.Sub(competitor.ActualStart)
		} else {
			previousFinish := competitor.LapTimes[competitor.CurrentLap-2].Finish
			lapTime = event.Time.Sub(previousFinish)
		}

		speed := float64(cfg.LapLen) / lapTime.Seconds()
		competitor.LapTimes[competitor.CurrentLap-1] = model.LapInfo{
			Time:   lapTime,
			Speed:  speed,
			Finish: event.Time,
		}

		competitor.CurrentLap++

		if competitor.CurrentLap > cfg.Laps {
			competitor.Status = model.StatusFinished
			finishEvent := model.Event{
				Time:         event.Time,
				EventID:      model.EventFinished,
				CompetitorID: competitor.ID,
				Processed:    true,
			}
			*processedEvents = append(*processedEvents, finishEvent)
		}
	}
}

func handleLostEvent(competitor *model.Competitor, event model.Event) {
	competitor.Status = model.StatusNotFinished
	competitor.StatusComment = event.ExtraParams
}
