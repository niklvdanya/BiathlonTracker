package model

import (
	"time"
)

type Competitor struct {
	ID             int
	RegisteredTime time.Time
	PlannedStart   time.Time
	ActualStart    time.Time
	Status         string
	LapTimes       []LapInfo
	PenaltyLapInfo PenaltyInfo
	CurrentLap     int
	CurrentFiring  int
	ShotsHit       int
	TotalShots     int
	InPenalty      bool
	OnFiringRange  bool
	StatusComment  string
	MissedShot     bool
}

type LapInfo struct {
	Time   time.Duration
	Speed  float64
	Finish time.Time
}

type PenaltyInfo struct {
	StartTime time.Time
	Duration  time.Duration
	Speed     float64
}

type Event struct {
	Time         time.Time
	EventID      int
	CompetitorID int
	ExtraParams  string
	Processed    bool
}
