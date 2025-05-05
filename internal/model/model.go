package model

import (
	"fmt"
	"time"
)

const (
	EventRegistration = 1
	EventSetStartTime = 2
	EventStartLine    = 3
	EventStarted      = 4
	EventFiringRange  = 5
	EventShot         = 6
	EventLeaveFiring  = 7
	EventEnterPenalty = 8
	EventLeavePenalty = 9
	EventLapEnd       = 10
	EventLostInForest = 11
	EventDisqualified = 32
	EventFinished     = 33
)

const (
	StatusFinished     = "Finished"
	StatusNotFinished  = "NotFinished"
	StatusNotStarted   = "NotStarted"
	StatusRunning      = "Running"
	StatusDisqualified = "Disqualified"
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

func (c *Competitor) TotalTime() time.Duration {
	var totalTime time.Duration
	for _, lap := range c.LapTimes {
		totalTime += lap.Time
	}
	return totalTime
}

func (c *Competitor) IsFinished() bool {
	return c.Status == StatusFinished
}

func (c *Competitor) IsRunning() bool {
	return c.Status == StatusRunning
}

func (c *Competitor) IsDisqualified() bool {
	return c.Status == StatusDisqualified
}

func (c *Competitor) ShotAccuracy() string {
	return fmt.Sprintf("%d/%d", c.ShotsHit, c.TotalShots)
}
