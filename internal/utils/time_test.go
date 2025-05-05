package utils

import (
	"testing"
	"time"
)

func TestFormatTimeRFC(t *testing.T) {
	cases := []struct {
		input    time.Time
		expected string
	}{
		{
			input:    time.Date(2025, 1, 1, 12, 34, 56, 789*1000000, time.UTC),
			expected: "12:34:56.789",
		},
		{
			input:    time.Date(2025, 1, 1, 9, 5, 0, 1*1000000, time.UTC),
			expected: "09:05:00.001",
		},
		{
			input:    time.Date(2025, 1, 1, 23, 59, 59, 999*1000000, time.UTC),
			expected: "23:59:59.999",
		},
	}

	for _, c := range cases {
		result := FormatTimeRFC(c.input)
		if result != c.expected {
			t.Errorf("FormatTimeRFC(%v): expected %s, got %s", c.input, c.expected, result)
		}
	}
}

func TestFormatDuration(t *testing.T) {
	cases := []struct {
		input    time.Duration
		expected string
	}{
		{
			input:    30*time.Minute + 15*time.Second + 500*time.Millisecond,
			expected: "00:30:15.500",
		},
		{
			input:    1*time.Hour + 15*time.Minute + 30*time.Second + 250*time.Millisecond,
			expected: "01:15:30.250",
		},
		{
			input:    2*time.Hour + 59*time.Minute + 59*time.Second + 999*time.Millisecond,
			expected: "02:59:59.999",
		},
		{
			input:    123 * time.Millisecond,
			expected: "00:00:00.123",
		},
	}

	for _, c := range cases {
		result := FormatDuration(c.input)
		if result != c.expected {
			t.Errorf("FormatDuration(%v): expected %s, got %s", c.input, c.expected, result)
		}
	}
}

func TestParseCompetitionTime(t *testing.T) {
	cases := []struct {
		input    string
		expected string
		wantErr  bool
	}{
		{
			input:    "10:15:30.500",
			expected: "10:15:30.500",
			wantErr:  false,
		},
		{
			input:    "09:05:00.001",
			expected: "09:05:00.001",
			wantErr:  false,
		},
		{
			input:    "23:59:59.999",
			expected: "23:59:59.999",
			wantErr:  false,
		},
		{
			input:    "invalid",
			expected: "",
			wantErr:  true,
		},
	}

	for _, c := range cases {
		result, err := ParseCompetitionTime(c.input)

		if c.wantErr {
			if err == nil {
				t.Errorf("ParseCompetitionTime(%s): expected error, got nil", c.input)
			}
			continue
		}

		if err != nil {
			t.Errorf("ParseCompetitionTime(%s): unexpected error: %v", c.input, err)
			continue
		}

		formatted := FormatTimeRFC(result)
		if formatted != c.expected {
			t.Errorf("ParseCompetitionTime(%s): expected %s, got %s", c.input, c.expected, formatted)
		}
	}
}
