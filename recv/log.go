package recv

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/frain-dev/immune"

	"github.com/frain-dev/immune/util"
)

type Log struct {
	AuthFailures      int64  `json:"auth_failures"`      // caused by authentication failure
	SignatureFailures int64  `json:"signature_failures"` // caused by signature failure
	ErrorRate         string `json:"error_rate"`
	SuccessRate       string `json:"success_rate"`

	captureLock      sync.Mutex        // protects the fields below
	EventsReceived   int               `json:"events_received"`
	EventRecvTime    map[string]string `json:"event_recv_time"`
	EventsDeliveries map[string]MI     `json:"events_deliveries"`
}

type MI map[string]int

func NewLog() *Log {
	return &Log{
		EventsReceived:    0,
		AuthFailures:      0,
		SignatureFailures: 0,
		EventRecvTime:     map[string]string{},
		EventsDeliveries:  map[string]MI{},
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

func (l *Log) CaptureHeaders(r *http.Request, now *time.Time) {
	l.captureLock.Lock()
	l.EventsReceived++

	eventID := r.Header.Get(immune.DefaultEventIDHeader)
	deliveryID := r.Header.Get(immune.DefaultEventDeliveryIDHeader)

	l.EventRecvTime[eventID] = now.UTC().Format(time.RFC3339)

	v := l.EventsDeliveries[eventID]
	if v == nil {
		v = MI{}
	}

	v[deliveryID]++

	l.EventsDeliveries[eventID] = v

	l.captureLock.Unlock()
}

func (l *Log) AddAuthFailure() {
	atomic.AddInt64(&l.AuthFailures, int64(1))
}

func (l *Log) AddSignatureFailure() {
	atomic.AddInt64(&l.SignatureFailures, int64(1))
}
