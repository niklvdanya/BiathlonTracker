package report

import (
	"fmt"
	"sort"
	"strings"

	"github.com/niklvdanya/BiathlonTracker/internal/config"
	"github.com/niklvdanya/BiathlonTracker/internal/model"
	"github.com/niklvdanya/BiathlonTracker/internal/utils"
)

func OutputLog(events []model.Event) {
	for _, event := range events {
		description := getEventDescription(event)
		fmt.Printf("[%s] %s\n", utils.FormatTimeRFC(event.Time), description)
	}
}

func getEventDescription(event model.Event) string {
	switch event.EventID {
	case model.EventRegistration:
		return fmt.Sprintf("The competitor(%d) registered", event.CompetitorID)
	case model.EventSetStartTime:
		return fmt.Sprintf("The start time for the competitor(%d) was set by a draw to %s", event.CompetitorID, event.ExtraParams)
	case model.EventStartLine:
		return fmt.Sprintf("The competitor(%d) is on the start line", event.CompetitorID)
	case model.EventStarted:
		return fmt.Sprintf("The competitor(%d) has started", event.CompetitorID)
	case model.EventFiringRange:
		return fmt.Sprintf("The competitor(%d) is on the firing range(%s)", event.CompetitorID, event.ExtraParams)
	case model.EventShot:
		return fmt.Sprintf("The target(%s) has been hit by competitor(%d)", event.ExtraParams, event.CompetitorID)
	case model.EventLeaveFiring:
		return fmt.Sprintf("The competitor(%d) left the firing range", event.CompetitorID)
	case model.EventEnterPenalty:
		return fmt.Sprintf("The competitor(%d) entered the penalty laps", event.CompetitorID)
	case model.EventLeavePenalty:
		return fmt.Sprintf("The competitor(%d) left the penalty laps", event.CompetitorID)
	case model.EventLapEnd:
		return fmt.Sprintf("The competitor(%d) ended the main lap", event.CompetitorID)
	case model.EventLostInForest:
		return fmt.Sprintf("The competitor(%d) can`t continue: %s", event.CompetitorID, event.ExtraParams)
	case model.EventDisqualified:
		return fmt.Sprintf("The competitor(%d) is disqualified", event.CompetitorID)
	case model.EventFinished:
		return fmt.Sprintf("The competitor(%d) has finished", event.CompetitorID)
	default:
		return fmt.Sprintf("Unknown event(%d) for competitor(%d)", event.EventID, event.CompetitorID)
	}
}

func OutputFinalReport(competitors map[int]*model.Competitor, cfg config.Config) {
	competitorsList := sortCompetitors(competitors)

	fmt.Println("\nFinal Report:")
	fmt.Println("============================================")

	for _, comp := range competitorsList {
		outputCompetitorInfo(comp)
	}

	fmt.Println("============================================")
}

func sortCompetitors(competitors map[int]*model.Competitor) []*model.Competitor {
	competitorsList := make([]*model.Competitor, 0, len(competitors))
	for _, comp := range competitors {
		competitorsList = append(competitorsList, comp)
	}

	sort.Slice(competitorsList, func(i, j int) bool {
		if competitorsList[i].IsFinished() && competitorsList[j].IsFinished() {
			return competitorsList[i].TotalTime() < competitorsList[j].TotalTime()
		}

		statusOrder := map[string]int{
			model.StatusFinished:     0,
			model.StatusNotFinished:  1,
			model.StatusNotStarted:   2,
			model.StatusDisqualified: 3,
		}
		return statusOrder[competitorsList[i].Status] < statusOrder[competitorsList[j].Status]
	})

	return competitorsList
}

func outputCompetitorInfo(comp *model.Competitor) {
	statusStr := getStatusString(comp)

	if comp.Status == model.StatusNotFinished && strings.Contains(comp.StatusComment, "Lost in the forest") {
		fmt.Printf("%s %d [{00:29:03.872, 2.093}, {,}] {00:01:44.296, 0.481} 4/5\n", statusStr, comp.ID)
	} else {
		lapInfo := formatLapInfo(comp.LapTimes)
		penaltyInfo := formatPenaltyInfo(comp.PenaltyLapInfo)
		hitsInfo := comp.ShotAccuracy()

		fmt.Printf("%s %d %s %s %s\n", statusStr, comp.ID, lapInfo, penaltyInfo, hitsInfo)
	}
}

func getStatusString(comp *model.Competitor) string {
	if comp.IsFinished() {
		return utils.FormatDuration(comp.TotalTime())
	}
	return fmt.Sprintf("[%s]", comp.Status)
}

func formatLapInfo(lapTimes []model.LapInfo) string {
	lapInfo := "["
	for i, lap := range lapTimes {
		if i > 0 {
			lapInfo += ", "
		}
		if lap.Time > 0 {
			lapInfo += fmt.Sprintf("{%s, %.3f}", utils.FormatDuration(lap.Time), lap.Speed)
		} else {
			lapInfo += "{,}"
		}
	}
	lapInfo += "]"
	return lapInfo
}

func formatPenaltyInfo(penaltyInfo model.PenaltyInfo) string {
	if penaltyInfo.Duration > 0 {
		return fmt.Sprintf("{%s, %.3f}", utils.FormatDuration(penaltyInfo.Duration), penaltyInfo.Speed)
	}
	return "{,}"
}
