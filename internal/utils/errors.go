package utils

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidTimeFormat   = errors.New("invalid time format")
	ErrInvalidEventFormat  = errors.New("invalid event format")
	ErrInvalidCompetitorID = errors.New("invalid competitor ID")
	ErrInvalidEventID      = errors.New("invalid event ID")
	ErrConfigNotFound      = errors.New("config file not found")
	ErrEventsNotFound      = errors.New("events file not found")
)

type ProcessingError struct {
	CompetitorID int
	EventID      int
	Message      string
}

func (e ProcessingError) Error() string {
	return fmt.Sprintf("error processing event %d for competitor %d: %s",
		e.EventID, e.CompetitorID, e.Message)
}

func NewProcessingError(competitorID, eventID int, message string) error {
	return &ProcessingError{
		CompetitorID: competitorID,
		EventID:      eventID,
		Message:      message,
	}
}
