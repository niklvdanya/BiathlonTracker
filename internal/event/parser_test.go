package event

import (
	"os"
	"testing"

	"github.com/niklvdanya/BiathlonTracker/internal/model"
)

func TestEventParsing(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected model.Event
		wantErr  bool
	}{
		{
			name:  "RegisterEvent",
			input: "[09:05:59.867] 1 1",
			expected: model.Event{
				EventID:      model.EventRegistration,
				CompetitorID: 1,
				ExtraParams:  "",
			},
			wantErr: false,
		},
		{
			name:  "SetStartTimeEvent",
			input: "[09:15:00.841] 2 1 09:30:00.000",
			expected: model.Event{
				EventID:      model.EventSetStartTime,
				CompetitorID: 1,
				ExtraParams:  "09:30:00.000",
			},
			wantErr: false,
		},
		{
			name:  "LostInForestEvent",
			input: "[09:59:03.872] 11 1 Lost in the forest",
			expected: model.Event{
				EventID:      model.EventLostInForest,
				CompetitorID: 1,
				ExtraParams:  "Lost in the forest",
			},
			wantErr: false,
		},
		{
			name:     "InvalidFormat",
			input:    "Invalid format",
			expected: model.Event{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseEvent(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result.EventID != tt.expected.EventID {
				t.Errorf("EventID mismatch: got %d, want %d", result.EventID, tt.expected.EventID)
			}

			if result.CompetitorID != tt.expected.CompetitorID {
				t.Errorf("CompetitorID mismatch: got %d, want %d", result.CompetitorID, tt.expected.CompetitorID)
			}

			if result.ExtraParams != tt.expected.ExtraParams {
				t.Errorf("ExtraParams mismatch: got %s, want %s", result.ExtraParams, tt.expected.ExtraParams)
			}
		})
	}
}

func TestLoadEvents(t *testing.T) {
	tempFile, err := os.CreateTemp("", "events_test.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	testData := `[09:05:59.867] 1 1
[09:15:00.841] 2 1 09:30:00.000
[09:29:45.734] 3 1
[09:30:01.005] 4 1
`

	if _, err := tempFile.Write([]byte(testData)); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}

	if err := tempFile.Close(); err != nil {
		t.Fatalf("failed to close temp file: %v", err)
	}

	events, err := LoadEvents(tempFile.Name())
	if err != nil {
		t.Fatalf("LoadEvents failed: %v", err)
	}

	if len(events) != 4 {
		t.Errorf("expected 4 events, got %d", len(events))
	}

	expectedIDs := []int{model.EventRegistration, model.EventSetStartTime, model.EventStartLine, model.EventStarted}
	for i, e := range events {
		if e.EventID != expectedIDs[i] {
			t.Errorf("event[%d]: expected EventID %d, got %d", i, expectedIDs[i], e.EventID)
		}
	}
}
