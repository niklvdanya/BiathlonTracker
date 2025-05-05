package event

import (
	"github.com/niklvdanya/BiathlonTracker/internal/config"
	"github.com/niklvdanya/BiathlonTracker/internal/model"
)

type EventParser interface {
	Parse(filename string) ([]model.Event, error)
}

type EventProcessor interface {
	Process(events []model.Event, cfg config.Config) map[int]*model.Competitor
}

type DefaultEventParser struct{}

func (p *DefaultEventParser) Parse(filename string) ([]model.Event, error) {
	return LoadEvents(filename)
}

type DefaultEventProcessor struct{}

func (p *DefaultEventProcessor) Process(events []model.Event, cfg config.Config) map[int]*model.Competitor {
	return ProcessEvents(events, cfg)
}
