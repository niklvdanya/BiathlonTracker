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

	sort.Slice(events, func(i, j int) bool {
		return events[i].Time.Before(events[j].Time)
	})

	processedEvents := make([]model.Event, 0, len(events))

	for i := range events {
		event := events[i]
		if event.Processed {
			continue
		}
		event.Processed = true

		competitor, exists := competitors[event.CompetitorID]
		if !exists {
			competitor = &model.Competitor{
				ID:         event.CompetitorID,
				Status:     "NotStarted",
				LapTimes:   make([]model.LapInfo, cfg.Laps),
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
			// Участник на стартовой линии - ничего не делаем

		case 4:
			competitor.ActualStart = event.Time
			competitor.Status = "Running"

			if event.Time.Sub(competitor.PlannedStart) > startDeltaDuration {
				competitor.Status = "Disqualified"
				disqEvent := model.Event{
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
				speed := float64(cfg.PenaltyLen) / penaltyDuration.Seconds()
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

				speed := float64(cfg.LapLen) / lapTime.Seconds()
				competitor.LapTimes[competitor.CurrentLap-1] = model.LapInfo{
					Time:   lapTime,
					Speed:  speed,
					Finish: event.Time,
				}

				competitor.CurrentLap++

				if competitor.CurrentLap > cfg.Laps {
					competitor.Status = "Finished"
					finishEvent := model.Event{
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

	// Обновляем список событий с обработанными
	for i := range events {
		events[i] = processedEvents[i]
	}

	return competitors
}
