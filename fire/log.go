package fire

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

type Log struct {
	EventsSent     int               `json:"events_sent"`
	Failures       int               `json:"failures"`
	FailureCodes   map[int]int       `json:"failure_codes"`
	EventTime      map[string]string `json:"event_time"`
	TotalTimeTaken string            `json:"total_time_taken"`
}

func NewLog() *Log {
	return &Log{
		EventsSent:   0,
		Failures:     0,
		FailureCodes: map[int]int{},
		EventTime:    map[string]string{},
	}
}

func (l *Log) WriteToFile(path string) error {
	b, err := json.Marshal(l)
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if dir != "." && dir != ".." {
		err = os.Mkdir(dir, 0o777)
		if err != nil && !os.IsExist(err) {
			return fmt.Errorf("failed to create log directory: %v", err)
		}
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.WithError(err).Errorf("Unable to open file %q, %v", path, err)
		return err
	}

	defer file.Close()

	_, err = file.Write(b)
	if err != nil {
		return err
	}

	return nil
}
