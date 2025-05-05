package event

import (
	"github.com/niklvdanya/BiathlonTracker/internal/model"
)

func ExportedParseEvent(line string) (model.Event, error) {
	return parseEvent(line)
}
