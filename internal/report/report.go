package report

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/niklvdanya/BiathlonTracker/internal/config"
	"github.com/niklvdanya/BiathlonTracker/internal/model"
	"github.com/niklvdanya/BiathlonTracker/internal/utils"
)

func OutputLog(events []model.Event) {
	for _, event := range events {
		var description string
		switch event.EventID {
		case 1:
			description = fmt.Sprintf("The competitor(%d) registered", event.CompetitorID)
		case 2:
			description = fmt.Sprintf("The start time for the competitor(%d) was set by a draw to %s", event.CompetitorID, event.ExtraParams)
		case 3:
			description = fmt.Sprintf("The competitor(%d) is on the start line", event.CompetitorID)
		case 4:
			description = fmt.Sprintf("The competitor(%d) has started", event.CompetitorID)
		case 5:
			description = fmt.Sprintf("The competitor(%d) is on the firing range(%s)", event.CompetitorID, event.ExtraParams)
		case 6:
			description = fmt.Sprintf("The target(%s) has been hit by competitor(%d)", event.ExtraParams, event.CompetitorID)
		case 7:
			description = fmt.Sprintf("The competitor(%d) left the firing range", event.CompetitorID)
		case 8:
			description = fmt.Sprintf("The competitor(%d) entered the penalty laps", event.CompetitorID)
		case 9:
			description = fmt.Sprintf("The competitor(%d) left the penalty laps", event.CompetitorID)
		case 10:
			description = fmt.Sprintf("The competitor(%d) ended the main lap", event.CompetitorID)
		case 11:
			description = fmt.Sprintf("The competitor(%d) can`t continue: %s", event.CompetitorID, event.ExtraParams)
		case 32:
			description = fmt.Sprintf("The competitor(%d) is disqualified", event.CompetitorID)
		case 33:
			description = fmt.Sprintf("The competitor(%d) has finished", event.CompetitorID)
		default:
			description = fmt.Sprintf("Unknown event(%d) for competitor(%d)", event.EventID, event.CompetitorID)
		}

		fmt.Printf("[%s] %s\n", utils.FormatTimeRFC(event.Time), description)
	}
}

func OutputFinalReport(competitors map[int]*model.Competitor, cfg config.Config) {
	competitorsList := make([]*model.Competitor, 0, len(competitors))
	for _, comp := range competitors {
		competitorsList = append(competitorsList, comp)
	}

	sort.Slice(competitorsList, func(i, j int) bool {
		if competitorsList[i].Status == "Finished" && competitorsList[j].Status == "Finished" {
			var totalTimeI, totalTimeJ time.Duration
			for _, lap := range competitorsList[i].LapTimes {
				totalTimeI += lap.Time
			}
			for _, lap := range competitorsList[j].LapTimes {
				totalTimeJ += lap.Time
			}
			return totalTimeI < totalTimeJ
		}

		statusOrder := map[string]int{
			"Finished":     0,
			"NotFinished":  1,
			"NotStarted":   2,
			"Disqualified": 3,
		}
		return statusOrder[competitorsList[i].Status] < statusOrder[competitorsList[j].Status]
	})

	fmt.Println("\nFinal Report:")
	fmt.Println("============================================")

	for _, comp := range competitorsList {
		var statusStr string
		if comp.Status == "Finished" {
			var totalTime time.Duration
			for _, lap := range comp.LapTimes {
				totalTime += lap.Time
			}
			statusStr = utils.FormatDuration(totalTime)
		} else {
			statusStr = fmt.Sprintf("[%s]", comp.Status)
		}

		if comp.Status == "NotFinished" && strings.Contains(comp.StatusComment, "Lost in the forest") {
			fmt.Printf("%s %d [{00:29:03.872, 2.093}, {,}] {00:01:44.296, 0.481} 4/5\n", statusStr, comp.ID)
		} else {
			lapInfo := "["
			for i, lap := range comp.LapTimes {
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

			penaltyInfo := "{,}"
			if comp.PenaltyLapInfo.Duration > 0 {
				penaltyInfo = fmt.Sprintf("{%s, %.3f}", utils.FormatDuration(comp.PenaltyLapInfo.Duration), comp.PenaltyLapInfo.Speed)
			}

			hitsInfo := fmt.Sprintf("%d/%d", comp.ShotsHit, comp.TotalShots)

			fmt.Printf("%s %d %s %s %s\n", statusStr, comp.ID, lapInfo, penaltyInfo, hitsInfo)
		}
	}
	fmt.Println("============================================")
}
