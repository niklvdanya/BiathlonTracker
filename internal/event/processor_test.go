package event

import (
	"context"
	"testing"
	"time"

	"github.com/niklvdanya/BiathlonTracker/internal/config"
	"github.com/niklvdanya/BiathlonTracker/internal/model"
)

func TestRegistrationEvent(t *testing.T) {
	ctx := context.Background()

	baseTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	cfg := config.Config{
		Laps:        2,
		LapLen:      3500,
		PenaltyLen:  150,
		FiringLines: 2,
		Start:       "10:00:00.000",
		StartDelta:  "00:01:30",
	}

	events := []model.Event{
		{
			Time:         baseTime.Add(9*time.Hour + 5*time.Minute),
			EventID:      model.EventRegistration,
			CompetitorID: 1,
		},
	}

	processor := &DefaultEventProcessor{}
	competitors := processor.Process(ctx, events, cfg)

	competitor, exists := competitors[1]
	if !exists {
		t.Fatalf("competitor 1 not found")
	}

	if competitor.Status != model.StatusNotStarted {
		t.Errorf("expected status %s, got %s", model.StatusNotStarted, competitor.Status)
	}

	if !competitor.RegisteredTime.Equal(events[0].Time) {
		t.Errorf("expected registered time %v, got %v", events[0].Time, competitor.RegisteredTime)
	}
}

func TestLostInForestEvent(t *testing.T) {
	ctx := context.Background()

	baseTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	cfg := config.Config{
		Laps:        2,
		LapLen:      3500,
		PenaltyLen:  150,
		FiringLines: 2,
		Start:       "10:00:00.000",
		StartDelta:  "00:01:30",
	}

	events := []model.Event{
		{
			Time:         baseTime.Add(10 * time.Hour),
			EventID:      model.EventLostInForest,
			CompetitorID: 1,
			ExtraParams:  model.LostInForestText,
		},
	}

	processor := &DefaultEventProcessor{}
	competitors := processor.Process(ctx, events, cfg)

	competitor, exists := competitors[1]
	if !exists {
		t.Fatalf("competitor 1 not found")
	}

	if competitor.Status != model.StatusNotFinished {
		t.Errorf("expected status %s, got %s", model.StatusNotFinished, competitor.Status)
	}

	if competitor.StatusComment != model.LostInForestText {
		t.Errorf("expected status comment %s, got %s", model.LostInForestText, competitor.StatusComment)
	}
}

func TestDisqualificationEvent(t *testing.T) {
	ctx := context.Background()

	baseTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	cfg := config.Config{
		Laps:        2,
		LapLen:      3500,
		PenaltyLen:  150,
		FiringLines: 2,
		Start:       "10:00:00.000",
		StartDelta:  "00:01:30",
	}

	events := []model.Event{
		{
			Time:         baseTime.Add(9*time.Hour + 5*time.Minute),
			EventID:      model.EventRegistration,
			CompetitorID: 1,
		},
		{
			Time:         baseTime.Add(9*time.Hour + 15*time.Minute),
			EventID:      model.EventSetStartTime,
			CompetitorID: 1,
			ExtraParams:  "10:00:00.000",
		},
		{
			Time:         baseTime.Add(10*time.Hour + 2*time.Minute),
			EventID:      model.EventStarted,
			CompetitorID: 1,
		},
	}

	processor := &DefaultEventProcessor{}
	competitors := processor.Process(ctx, events, cfg)

	competitor, exists := competitors[1]
	if !exists {
		t.Fatalf("competitor 1 not found")
	}

	if competitor.Status != model.StatusDisqualified {
		t.Errorf("expected status %s, got %s", model.StatusDisqualified, competitor.Status)
	}
}

func TestShotEvent(t *testing.T) {
	ctx := context.Background()

	baseTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	cfg := config.Config{
		Laps:        2,
		LapLen:      3500,
		PenaltyLen:  150,
		FiringLines: 2,
		Start:       "10:00:00.000",
		StartDelta:  "00:01:30",
	}

	events := []model.Event{
		{
			Time:         baseTime.Add(10 * time.Hour),
			EventID:      model.EventFiringRange,
			CompetitorID: 1,
			ExtraParams:  "1",
		},
		{
			Time:         baseTime.Add(10*time.Hour + time.Second),
			EventID:      model.EventShot,
			CompetitorID: 1,
			ExtraParams:  model.ShotTarget1,
		},
		{
			Time:         baseTime.Add(10*time.Hour + 2*time.Second),
			EventID:      model.EventShot,
			CompetitorID: 1,
			ExtraParams:  model.ShotTarget3,
		},
	}

	processor := &DefaultEventProcessor{}
	competitors := processor.Process(ctx, events, cfg)

	competitor, exists := competitors[1]
	if !exists {
		t.Fatalf("competitor 1 not found")
	}

	if competitor.ShotsHit != 1 {
		t.Errorf("expected 1 shot hit, got %d", competitor.ShotsHit)
	}

	if competitor.TotalShots != 2 {
		t.Errorf("expected 2 total shots, got %d", competitor.TotalShots)
	}

	if !competitor.MissedShot {
		t.Errorf("expected missed shot flag to be true")
	}
}

func TestBasicParallelProcessing(t *testing.T) {
	ctx := context.Background()
	cfg := config.Config{
		Laps:        2,
		LapLen:      3500,
		PenaltyLen:  150,
		FiringLines: 2,
		Start:       "10:00:00.000",
		StartDelta:  "00:01:30",
	}
	baseTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	events := []model.Event{
		{
			Time:         baseTime.Add(9 * time.Hour),
			EventID:      model.EventRegistration,
			CompetitorID: 1,
		},
		{
			Time:         baseTime.Add(9*time.Hour + 5*time.Minute),
			EventID:      model.EventRegistration,
			CompetitorID: 2,
		},
	}

	competitors := ProcessEventsParallel(ctx, events, cfg)

	if len(competitors) != 2 {
		t.Errorf("expected 2 competitors, got %d", len(competitors))
	}
	comp1, exists := competitors[1]
	if !exists {
		t.Fatalf("competitor 1 not found")
	}

	if comp1.Status != model.StatusNotStarted {
		t.Errorf("competitor 1: expected status %s, got %s", model.StatusNotStarted, comp1.Status)
	}

	comp2, exists := competitors[2]
	if !exists {
		t.Fatalf("competitor 2 not found")
	}

	if comp2.Status != model.StatusNotStarted {
		t.Errorf("competitor 2: expected status %s, got %s", model.StatusNotStarted, comp2.Status)
	}
}

func TestParallelEvent(t *testing.T) {
	ctx := context.Background()
	cfg := config.Config{
		Laps:        2,
		LapLen:      3500,
		PenaltyLen:  150,
		FiringLines: 2,
		Start:       "10:00:00.000",
		StartDelta:  "00:01:30",
	}

	baseTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	events := []model.Event{
		{
			Time:         baseTime.Add(10 * time.Hour),
			EventID:      model.EventFiringRange,
			CompetitorID: 1,
			ExtraParams:  "1",
		},
		{
			Time:         baseTime.Add(10*time.Hour + 1*time.Second),
			EventID:      model.EventShot,
			CompetitorID: 1,
			ExtraParams:  model.ShotTarget1,
		},
		{
			Time:         baseTime.Add(10 * time.Hour),
			EventID:      model.EventFiringRange,
			CompetitorID: 2,
			ExtraParams:  "1",
		},
		{
			Time:         baseTime.Add(10*time.Hour + 1*time.Second),
			EventID:      model.EventShot,
			CompetitorID: 2,
			ExtraParams:  model.ShotTarget3,
		},
	}

	competitors := ProcessEventsParallel(ctx, events, cfg)
	comp1, exists := competitors[1]
	if !exists {
		t.Fatalf("competitor 1 not found")
	}

	if comp1.ShotsHit != 1 {
		t.Errorf("competitor 1: expected 1 shot hit, got %d", comp1.ShotsHit)
	}

	comp2, exists := competitors[2]
	if !exists {
		t.Fatalf("competitor 2 not found")
	}

	if comp2.ShotsHit != 0 {
		t.Errorf("competitor 2: expected 0 shots hit, got %d", comp2.ShotsHit)
	}

	if comp2.MissedShot != true {
		t.Errorf("competitor 2: expected missed shot to be true")
	}
}
