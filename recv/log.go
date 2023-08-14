package recv

import (
	"fmt"

	"github.com/frain-dev/immune/util"
)

type Log struct {
	EventsReceived    int               `json:"events_sent"`
	AuthFailures      int               `json:"auth_failures"`      // caused by authentication failure
	SignatureFailures int               `json:"signature_failures"` // caused by signature failure
	EventRecvTime     map[string]string `json:"event_recv_time"`
	EventCount        map[string]int    `json:"event_count"`
	ErrorRate         string            `json:"error_rate"`
	SuccessRate       string            `json:"success_rate"`
}

func NewLog() *Log {
	return &Log{
		EventsReceived:    0,
		AuthFailures:      0,
		SignatureFailures: 0,
		EventRecvTime:     map[string]string{},
		EventCount:        map[string]int{},
		ErrorRate:         "",
		SuccessRate:       "",
	}
}

func (l *Log) WriteToFile(path string) error {
	return util.WriteJSONToFile(path, l)
}

func (l *Log) CalculateStats() {
	fails := l.AuthFailures + l.SignatureFailures
	errRate := calculatePercentage(float64(fails), float64(l.EventsReceived))

	l.ErrorRate = fmt.Sprintf("%.2f%%", errRate)
	l.SuccessRate = fmt.Sprintf("%.2f%%", 100-errRate)
}

func calculatePercentage(part, whole float64) float64 {
	return (part / whole) * 100
}
