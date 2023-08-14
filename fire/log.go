package fire

import (
	"fmt"

	"github.com/frain-dev/immune/util"
)

type Log struct {
	EventsSent         int               `json:"events_sent"`
	RequestDurations   []int64           `json:"request_durations,omitempty"`
	Failures           int               `json:"failures"`
	FailureCodes       map[int]int       `json:"failure_codes"`
	EventTime          map[string]string `json:"event_time"`
	MinimumRequestTime string            `json:"minimum_request_time"`
	MaximumRequestTime string            `json:"maximum_request_time"`
	AverageRequestTime string            `json:"average_request_time"`
	ErrorRate          string            `json:"error_rate"`
	SuccessRate        string            `json:"success_rate"`
	TotalTimeTaken     string            `json:"total_time_taken"`
}

func NewLog(n int64) *Log {
	return &Log{
		EventsSent:       0,
		Failures:         0,
		RequestDurations: make([]int64, 0, n),
		FailureCodes:     map[int]int{},
		EventTime:        map[string]string{},
	}
}

func (l *Log) WriteToFile(path string) error {
	return util.WriteJSONToFile(path, l)
}

func (l *Log) CalculateStats() {
	errRate := calculatePercentage(float64(l.Failures), float64(l.EventsSent))

	l.ErrorRate = fmt.Sprintf("%.2f%%", errRate)
	l.SuccessRate = fmt.Sprintf("%.2f%%", 100-errRate)

	min, max, avg := analyzeDurations(l.RequestDurations)
	l.MinimumRequestTime = fmt.Sprintf("%d milliseconds", min)
	l.MaximumRequestTime = fmt.Sprintf("%d milliseconds", max)
	l.AverageRequestTime = fmt.Sprintf("%.2f milliseconds", avg)
	l.RequestDurations = nil
}

func calculatePercentage(part, whole float64) float64 {
	return (part / whole) * 100
}

func analyzeDurations(durations []int64) (int64, int64, float64) {
	if len(durations) == 0 {
		return 0, 0, 0
	}

	min := durations[0]
	max := durations[0]
	sum := int64(0)

	for _, num := range durations {
		if num < min {
			min = num
		}
		if num > max {
			max = num
		}
		sum += num
	}

	average := float64(sum) / float64(len(durations))

	return min, max, average
}
