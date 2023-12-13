package data

import "time"

type OverallStatistics struct {
	MediumDistance       int `json:"medium_distance"`
	MediumDuration       int `json:"medium_time"`
	MaxDistance          int `json:"max_distance"`
	MaxDuration          int `json:"max_time"`
	MediumWeeklyDistance int `json:"medium_weekly_distance"`
	MediumWeeklyDuration int `json:"medium_weekly_time"`
	MaxWeeklyDistance    int `json:"max_weekly_distance"`
	MaxWeeklyDuration    int `json:"max_weekly_time"`
}

type WeeklyData struct {
	StartDate time.Time
	EndDate   time.Time
	Distance  int
	Duration  int
}
