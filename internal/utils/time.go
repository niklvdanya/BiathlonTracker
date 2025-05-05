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
