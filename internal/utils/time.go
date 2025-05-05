package utils

import (
	"fmt"
	"time"
)

func FormatTimeRFC(t time.Time) string {
	return fmt.Sprintf("%02d:%02d:%02d.%03d",
		t.Hour(), t.Minute(), t.Second(),
		t.Nanosecond()/1000000)
}

func FormatDuration(d time.Duration) string {
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	d -= s * time.Second
	ms := d / time.Millisecond

	return fmt.Sprintf("%02d:%02d:%02d.%03d", h, m, s, ms)
}

func ParseCompetitionTime(timeStr string) (time.Time, error) {
	now := time.Now()
	baseDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

	t, err := time.Parse("15:04:05.000", timeStr)
	if err != nil {
		return time.Time{}, err
	}

	return baseDate.Add(time.Duration(t.Hour())*time.Hour +
		time.Duration(t.Minute())*time.Minute +
		time.Duration(t.Second())*time.Second +
		time.Duration(t.Nanosecond())), nil
}
