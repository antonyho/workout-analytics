package main

import (
	"embed"
	"math/rand"
	"testing"
	"time"

	"github.com/antonyho/workout-analytics/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/*.json
var testData embed.FS

func Test_read(t *testing.T) {
	requestJSON, err := testData.Open("testdata/request.json")
	require.NoError(t, err)

	expected := []data.Workout{
		{Distance: 10000, Duration: 3600, Timestamp: "2023-11-04T13:43:28.073909Z"},
		{Distance: 24000, Duration: 6900, Timestamp: "2023-11-06T10:12:43.296345Z"},
		{Distance: 20000, Duration: 6300, Timestamp: "2023-11-10T10:12:43.296345Z"},
		{Distance: 18000, Duration: 6000, Timestamp: "2023-11-14T10:12:43.296345Z"},
		{Distance: 24000, Duration: 6700, Timestamp: "2023-11-17T10:12:43.296345Z"},
		{Distance: 23500, Duration: 6700, Timestamp: "2023-11-20T10:12:43.296345Z"},
		{Distance: 25000, Duration: 6900, Timestamp: "2023-11-22T10:12:43.296345Z"},
		{Distance: 25000, Duration: 6900, Timestamp: "2023-11-26T10:12:43.296345Z"},
		{Distance: 28000, Duration: 7100, Timestamp: "2023-11-28T10:12:43.296345Z"},
		{Distance: 27000, Duration: 6900, Timestamp: "2023-11-30T10:12:43.296345Z"},
		{Distance: 28000, Duration: 7000, Timestamp: "2023-12-02T10:12:43.296345Z"},
		{Distance: 28000, Duration: 6900, Timestamp: "2023-12-05T10:12:43.296345Z"},
		{Distance: 30000, Duration: 7200, Timestamp: "2023-12-08T10:12:43.296345Z"},
		{Distance: 30000, Duration: 7000, Timestamp: "2023-12-10T10:12:43.296345Z"},
		{Distance: 31000, Duration: 7000, Timestamp: "2023-12-12T10:12:43.296345Z"},
	}

	actual, err := read(requestJSON)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func Test_getStartingDate(t *testing.T) {
	t.Run("on Monday", func(t *testing.T) {
		lastDate := time.Date(2023, 10, 23, 8, 30, 0, 0, time.UTC)
		numberOfWeeks := 3

		expected := time.Date(2023, 10, 02, 0, 0, 0, 0, time.UTC)

		actual := getStartingDate(lastDate, numberOfWeeks)
		assert.Equal(t, expected, actual)
	})

	t.Run("before Monday", func(t *testing.T) {
		lastDate := time.Date(2023, 10, 20, 8, 30, 0, 0, time.UTC)
		numberOfWeeks := 2

		expected := time.Date(2023, 10, 02, 0, 0, 0, 0, time.UTC)

		actual := getStartingDate(lastDate, numberOfWeeks)
		assert.Equal(t, expected, actual)
	})

	t.Run("after Monday", func(t *testing.T) {
		lastDate := time.Date(2023, 10, 26, 8, 30, 0, 0, time.UTC)
		numberOfWeeks := 3

		expected := time.Date(2023, 10, 02, 0, 0, 0, 0, time.UTC)

		actual := getStartingDate(lastDate, numberOfWeeks)
		assert.Equal(t, expected, actual)
	})
}

func Test_figure(t *testing.T) {
	today := time.Date(2023, 10, 26, 8, 30, 0, 0, time.UTC)

	orderedDataset := []data.Workout{
		{
			Distance:  1000,
			Duration:  200,
			Timestamp: "2023-10-03T10:08:21.000000Z",
		},
		{
			Distance:  1200,
			Duration:  370,
			Timestamp: "2023-10-06T10:08:21.000000Z",
		},
		{
			Distance:  800,
			Duration:  400,
			Timestamp: "2023-10-10T10:08:21.000000Z",
		},
		{
			Distance:  950,
			Duration:  350,
			Timestamp: "2023-10-12T10:08:21.000000Z",
		},
		{
			Distance:  1400,
			Duration:  600,
			Timestamp: "2023-10-13T10:08:21.000000Z",
		},
		{
			Distance:  600,
			Duration:  300,
			Timestamp: "2023-10-17T10:08:21.000000Z",
		},
		{
			Distance:  500,
			Duration:  700,
			Timestamp: "2023-10-21T10:08:21.000000Z",
		},
		{
			Distance:  1000,
			Duration:  450,
			Timestamp: "2023-10-24T10:08:21.000000Z",
		},
	}

	shuffledDataset := orderedDataset
	rand.Shuffle(len(orderedDataset), func(i, j int) {
		shuffledDataset[i], shuffledDataset[j] = shuffledDataset[j], shuffledDataset[i]
	})

	t.Run("covering all workouts", func(t *testing.T) {
		numberOfWeeks := 3

		expected := data.OverallStatistics{
			MaxDistance:          1400,
			MaxDuration:          700,
			MediumDistance:       975,
			MediumDuration:       385,
			MaxWeeklyDistance:    3150,
			MaxWeeklyDuration:    1350,
			MediumWeeklyDistance: 1650,
			MediumWeeklyDuration: 785,
		}

		actual, err := figure(today, numberOfWeeks, orderedDataset)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("not covering all workouts", func(t *testing.T) {
		numberOfWeeks := 2

		expected := data.OverallStatistics{
			MaxDistance:          1400,
			MaxDuration:          700,
			MediumDistance:       875,
			MediumDuration:       425,
			MaxWeeklyDistance:    3150,
			MaxWeeklyDuration:    1350,
			MediumWeeklyDistance: 1100,
			MediumWeeklyDuration: 1000,
		}

		actual, err := figure(today, numberOfWeeks, orderedDataset)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("shuffled workouts", func(t *testing.T) {
		numberOfWeeks := 3

		expected := data.OverallStatistics{
			MaxDistance:          1400,
			MaxDuration:          700,
			MediumDistance:       975,
			MediumDuration:       385,
			MaxWeeklyDistance:    3150,
			MaxWeeklyDuration:    1350,
			MediumWeeklyDistance: 1650,
			MediumWeeklyDuration: 785,
		}

		actual, err := figure(today, numberOfWeeks, shuffledDataset)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}

func Test_findMaxAndMedian(t *testing.T) {
	t.Run("odd", func(t *testing.T) {
		dataset := []int{3, 5, 1, 9, 7}

		expectedMax := 9
		expectedMedian := 5

		resultMax, resultMedian := findMaxAndMedian(dataset)

		assert.Equal(t, expectedMax, resultMax)
		assert.Equal(t, expectedMedian, resultMedian)
	})

	t.Run("even", func(t *testing.T) {
		dataset := []int{3, 5, 1, 9, 7, 11}

		expectedMax := 11
		expectedMedian := 6

		resultMax, resultMedian := findMaxAndMedian(dataset)

		assert.Equal(t, expectedMax, resultMax)
		assert.Equal(t, expectedMedian, resultMedian)
	})

	t.Run("even/floating-point-median", func(t *testing.T) {
		dataset := []int{3, 5, 1, 9, 8, 11}

		expectedMax := 11
		expectedMedian := 7 // Rounded up

		resultMax, resultMedian := findMaxAndMedian(dataset)

		assert.Equal(t, expectedMax, resultMax)
		assert.Equal(t, expectedMedian, resultMedian)
	})
}
