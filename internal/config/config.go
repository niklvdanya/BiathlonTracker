package config

import (
	"encoding/json"
	"os"
	"strconv"
)

type Config struct {
	Laps        int    `json:"laps"`
	LapLen      int    `json:"lapLen"`
	PenaltyLen  int    `json:"penaltyLen"`
	FiringLines int    `json:"firingLines"`
	Start       string `json:"start"`
	StartDelta  string `json:"startDelta"`
}

func Load(filename string) (Config, error) {
	var config Config

	data, err := os.ReadFile(filename)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	if laps, exists := os.LookupEnv("BIATHLON_LAPS"); exists {
		if value, err := strconv.Atoi(laps); err == nil {
			config.Laps = value
		}
	}

	if lapLen, exists := os.LookupEnv("BIATHLON_LAP_LEN"); exists {
		if value, err := strconv.Atoi(lapLen); err == nil {
			config.LapLen = value
		}
	}

	if penaltyLen, exists := os.LookupEnv("BIATHLON_PENALTY_LEN"); exists {
		if value, err := strconv.Atoi(penaltyLen); err == nil {
			config.PenaltyLen = value
		}
	}

	if firingLines, exists := os.LookupEnv("BIATHLON_FIRING_LINES"); exists {
		if value, err := strconv.Atoi(firingLines); err == nil {
			config.FiringLines = value
		}
	}

	if start, exists := os.LookupEnv("BIATHLON_START"); exists {
		config.Start = start
	}

	if startDelta, exists := os.LookupEnv("BIATHLON_START_DELTA"); exists {
		config.StartDelta = startDelta
	}

	return config, nil
}
