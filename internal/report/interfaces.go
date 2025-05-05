package report

import (
	"github.com/niklvdanya/BiathlonTracker/internal/config"
	"github.com/niklvdanya/BiathlonTracker/internal/model"
)

type Reporter interface {
	OutputLog(events []model.Event)
	OutputFinalReport(competitors map[int]*model.Competitor, cfg config.Config)
}

type DefaultReporter struct{}

func (r *DefaultReporter) OutputLog(events []model.Event) {
	OutputLog(events)
}

func (r *DefaultReporter) OutputFinalReport(competitors map[int]*model.Competitor, cfg config.Config) {
	OutputFinalReport(competitors, cfg)
}
