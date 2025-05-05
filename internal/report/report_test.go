package report

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/niklvdanya/BiathlonTracker/internal/config"
	"github.com/niklvdanya/BiathlonTracker/internal/model"
)

func TestOutputLog(t *testing.T) {
	now := time.Now()
	events := []model.Event{
		{
			Time:         now,
			EventID:      model.EventRegistration,
			CompetitorID: 1,
		},
		{
			Time:         now.Add(10 * time.Minute),
			EventID:      model.EventSetStartTime,
			CompetitorID: 1,
			ExtraParams:  "10:00:00.000",
		},
		{
			Time:         now.Add(20 * time.Minute),
			EventID:      model.EventStarted,
			CompetitorID: 1,
		},
	}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	OutputLog(events)

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	expectedStrings := []string{
		"registered",
		"start time",
		"started",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected output to contain '%s', got: %s", expected, output)
		}
	}
}

func TestOutputFinalReport(t *testing.T) {
	now := time.Now()
	baseTime := now.Truncate(24 * time.Hour)
	startTime := baseTime.Add(10 * time.Hour)

	competitors := map[int]*model.Competitor{
		1: {
			ID:          1,
			Status:      model.StatusFinished,
			CurrentLap:  3,
			ActualStart: startTime,
			LapTimes: []model.LapInfo{
				{
					Time:   15 * time.Minute,
					Speed:  3.89,
					Finish: startTime.Add(15 * time.Minute),
				},
				{
					Time:   14 * time.Minute,
					Speed:  4.17,
					Finish: startTime.Add(29 * time.Minute),
				},
			},
			ShotsHit:   4,
			TotalShots: 5,
		},
		2: {
			ID:            2,
			Status:        model.StatusNotFinished,
			StatusComment: "Lost in the forest",
			CurrentLap:    2,
			ActualStart:   startTime.Add(1*time.Minute + 30*time.Second),
			LapTimes: []model.LapInfo{
				{
					Time:   15 * time.Minute,
					Speed:  3.89,
					Finish: startTime.Add(16*time.Minute + 30*time.Second),
				},
			},
			ShotsHit:   4,
			TotalShots: 5,
		},
		3: {
			ID:          3,
			Status:      model.StatusDisqualified,
			ActualStart: startTime.Add(3 * time.Minute),
		},
	}

	cfg := config.Config{
		Laps:       2,
		LapLen:     3500,
		PenaltyLen: 150,
	}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	OutputFinalReport(competitors, cfg)

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	expectedStrings := []string{
		"Final Report",
		"00:29:00.000",
		"[NotFinished]",
		"[Disqualified]",
		"4/5",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected output to contain '%s', got: %s", expected, output)
		}
	}
}

func TestSortingCompetitors(t *testing.T) {
	competitors := map[int]*model.Competitor{
		1: {ID: 1, Status: model.StatusFinished, LapTimes: []model.LapInfo{{Time: 30 * time.Minute}}},
		2: {ID: 2, Status: model.StatusFinished, LapTimes: []model.LapInfo{{Time: 25 * time.Minute}}},
		3: {ID: 3, Status: model.StatusNotFinished},
		4: {ID: 4, Status: model.StatusDisqualified},
	}

	sorted := sortCompetitors(competitors)

	if len(sorted) != 4 {
		t.Fatalf("Expected 4 competitors, got %d", len(sorted))
	}

	if sorted[0].ID != 2 {
		t.Errorf("Expected first competitor to be ID 2, got %d", sorted[0].ID)
	}

	if sorted[1].ID != 1 {
		t.Errorf("Expected second competitor to be ID 1, got %d", sorted[1].ID)
	}

	if sorted[2].ID != 3 || sorted[2].Status != model.StatusNotFinished {
		t.Errorf("Expected third competitor to be NotFinished, got %s", sorted[2].Status)
	}

	if sorted[3].ID != 4 || sorted[3].Status != model.StatusDisqualified {
		t.Errorf("Expected fourth competitor to be Disqualified, got %s", sorted[3].Status)
	}
}
