package compute

import (
	"encoding/json"
	"io"
	"math"
	"os"
	"time"

	"github.com/frain-dev/immune/fire"
	"github.com/frain-dev/immune/recv"
)

type Compute struct {
	FireLog *fire.Log
	RecvLog *recv.Log
	l       *Log
}

func NewCompute(fireLogPath, recvLogPath string) (*Compute, error) {
	fireFile, err := os.Open(fireLogPath)
	if err != nil {
		return nil, err
	}
	defer fireFile.Close()

	recvFile, err := os.Open(recvLogPath)
	if err != nil {
		return nil, err
	}
	defer recvFile.Close()

	var fireLog fire.Log
	if err := decodeJSON(fireFile, &fireLog); err != nil {
		return nil, err
	}

	var recvLog recv.Log
	if err := decodeJSON(recvFile, &recvLog); err != nil {
		return nil, err
	}

	return &Compute{
		FireLog: &fireLog,
		RecvLog: &recvLog,
	}, nil
}

func decodeJSON(file io.Reader, v interface{}) error {
	decoder := json.NewDecoder(file)
	return decoder.Decode(v)
}

type Log struct {
	UnreceivedEvents    map[string]string `json:"unknown_events"`
	timeToReceiveEvents []int64

	MinimumEventRecvTime string `json:"minimum_event_recv_time"`
	MaximumEventRecvTime string `json:"maximum_event_recv_time"`
	AverageEventRecvTime string `json:"average_event_recv_time"`

	// Fire fields
	FireMinimumRequestTime string `json:"fire_minimum_request_time"`
	FireMaximumRequestTime string `json:"fire_maximum_request_time"`
	FireAverageRequestTime string `json:"fire_average_request_time"`
	FireErrorRate          string `json:"fire_error_rate"`
	FireSuccessRate        string `json:"fire_success_rate"`
	FireEventsSent         int    `json:"fire_events_sent"`

	RecvEventsReceived int    `json:"recv_events_received"`
	RecvErrorRate      string `json:"recv_error_rate"`
	RecvSuccessRate    string `json:"recv_success_rate"`
}

func (c *Compute) ComputeDiffLog() error {
	for eventID, t := range c.FireLog.EventTime {
		sendTime, err := time.Parse(t, time.RFC3339)
		if err != nil {
			return err
		}

		rt, ok := c.RecvLog.EventRecvTime[eventID]
		if !ok {
			c.l.UnreceivedEvents[eventID] = t // record the unreceived event and the time it was sent
		}

		recvTime, err := time.Parse(rt, time.RFC3339)
		if err != nil {
			return err
		}

		interval := recvTime.Sub(sendTime).Milliseconds()
		c.l.timeToReceiveEvents = append(c.l.timeToReceiveEvents, interval)
	}

	min, max, avg := fire.AnalyzeDurations(c.l.timeToReceiveEvents)

	min = math.Ceil()
}
