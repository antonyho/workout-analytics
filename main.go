package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/antonyho/workout-analytics/data"
	log "github.com/antonyho/workout-analytics/logger"
)

const (
	HealthCheckPath = "/health"
	AnalyserPath    = "/analyse"

	oneDay  = 24 * time.Hour
	oneWeek = 7 * oneDay

	// Health check response
	healthCheckHealthyMessage = "{'status': 'available'}"
)

func main() {
	http.HandleFunc(HealthCheckPath, healthCheck)
	http.HandleFunc(AnalyserPath, analyse)
	err := http.ListenAndServe(":8080", nil)
	if !errors.Is(err, http.ErrServerClosed) {
		log.Error("HTTP Server error. %v", err)
	}
}

func healthCheck(response http.ResponseWriter, _ *http.Request) {
	response.Header().Set("Cache-Control", "no-cache")
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	_, err := response.Write([]byte(healthCheckHealthyMessage))
	if err != nil {
		log.Error("Failed to write healthCheck HTTP response stream. %v", err)
	}
}

func analyse(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		// This endpoint only handle POST request
		log.Info("Non POST request got rejected. %v", request)
		response.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	numberOfWeeks, err := strconv.Atoi(request.FormValue("nweeks"))
	if err != nil {
		log.Error("Failed to read 'nweeks' parameter from request. %v", err)
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Info("New request for %d week(s) of analysis.", numberOfWeeks)

	workouts, err := read(request.Body)
	if err != nil {
		log.Error("Failed to read request JSON body. %v", err)
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	today := time.Now()
	overallStatistics, err := figure(today, numberOfWeeks, workouts)
	if err != nil {
		log.Error("Malformed Request JSON body. %v", err)
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(response).Encode(overallStatistics); err != nil {
		log.Error("Failed to encode analysis result to JSON. %v", err)
	}
}

func read(requestBody io.ReadCloser) ([]data.Workout, error) {
	workouts := make([]data.Workout, 0)
	entryIndex := 0
	jsonDecoder := json.NewDecoder(requestBody)
	defer requestBody.Close()
	if _, err := jsonDecoder.Token(); err != nil {
		return nil, fmt.Errorf("invalid json format. %w", err)
	}
	for jsonDecoder.More() {
		var workout data.Workout
		if err := jsonDecoder.Decode(&workout); err != nil {
			return nil, fmt.Errorf("error reading #%d object. %w", entryIndex, err)
		}
		workouts = append(workouts, workout)
		entryIndex++
	}

	return workouts, nil
}

// getStartingDate returns the date back of numberOfWeeks from the closest Monday before the endingDate.
// If endingDate is Monday, it simply returns the date back of numberOfWeeks from endingDate.
func getStartingDate(endingDate time.Time, numberOfWeeks int) time.Time {
	offsetToMonday := endingDate.Weekday() - time.Monday
	if offsetToMonday < 0 {
		offsetToMonday -= 7
	}
	closestMonday := endingDate.AddDate(0, 0, -int(offsetToMonday))

	return closestMonday.Add(time.Duration(-numberOfWeeks) * oneWeek).Truncate(oneDay)
}

func figure(endingDate time.Time, numberOfWeeks int, workouts []data.Workout) (data.OverallStatistics, error) {
	var statistics data.OverallStatistics

	startingDate := getStartingDate(endingDate, numberOfWeeks)

	allDistances := make([]int, 0)
	allDurations := make([]int, 0)
	allWeeklyData := make([]data.WeeklyData, 0)
	for day := startingDate; day.Before(endingDate); day = day.Add(7 * oneDay) {
		weeklyData := data.WeeklyData{
			StartDate: day,
			EndDate:   day.Add(7 * oneDay).Add(-time.Nanosecond),
		}
		allWeeklyData = append(allWeeklyData, weeklyData)
	}

	for _, workout := range workouts {
		for idx, weeklyData := range allWeeklyData {
			workoutDate, err := workout.Time()
			if err != nil {
				return statistics, fmt.Errorf("invalid date format at #%d JSON object. %w", idx, err)
			}
			if workoutDate.After(weeklyData.StartDate) && workoutDate.Before(weeklyData.EndDate) {
				allDistances = append(allDistances, workout.Distance)
				allDurations = append(allDurations, workout.Duration)

				allWeeklyData[idx].Distance += workout.Distance
				allWeeklyData[idx].Duration += workout.Duration
			}
		}
	}

	allWeeklyDistances := make([]int, len(allWeeklyData))
	allWeeklyDurations := make([]int, len(allWeeklyData))

	statistics.MaxDistance, statistics.MediumDistance = findMaxAndMedian(allDistances)
	statistics.MaxDuration, statistics.MediumDuration = findMaxAndMedian(allDurations)
	for idx, weeklyData := range allWeeklyData {
		if statistics.MaxWeeklyDistance < weeklyData.Distance {
			statistics.MaxWeeklyDistance = weeklyData.Distance
		}
		if statistics.MaxWeeklyDuration < weeklyData.Duration {
			statistics.MaxWeeklyDuration = weeklyData.Duration
		}

		allWeeklyDistances[idx] = weeklyData.Distance
		allWeeklyDurations[idx] = weeklyData.Duration
	}
	statistics.MaxWeeklyDistance, statistics.MediumWeeklyDistance = findMaxAndMedian(allWeeklyDistances)
	statistics.MaxWeeklyDuration, statistics.MediumWeeklyDuration = findMaxAndMedian(allWeeklyDurations)

	return statistics, nil
}

func findMaxAndMedian(figures []int) (max int, median int) {
	sort.Sort(sort.Reverse(sort.IntSlice(figures)))

	max = figures[0]

	dataSize := len(figures)
	if dataSize%2 == 1 {
		median = figures[dataSize/2]
	} else {
		median = int(math.Round(float64(figures[dataSize/2-1]+figures[dataSize/2]) / 2))
	}

	return
}
