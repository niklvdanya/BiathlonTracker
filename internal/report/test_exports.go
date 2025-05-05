package report

import (
	"github.com/niklvdanya/BiathlonTracker/internal/model"
)

func ExportedSortCompetitors(competitors map[int]*model.Competitor) []*model.Competitor {
	return sortCompetitors(competitors)
}
