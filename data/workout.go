package data

import "time"

type Workout struct {
	Distance  int    `json:"distance"`
	Duration  int    `json:"time"`
	Timestamp string `json:"timestamp"`
}

func (workout Workout) Time() (time.Time, error) {
	return time.Parse(time.RFC3339Nano, workout.Timestamp)
}
