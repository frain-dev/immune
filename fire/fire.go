package fire

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/schollz/progressbar/v3"

	"github.com/frain-dev/convoy/api/models"
	"github.com/frain-dev/convoy/datastore"
	"github.com/frain-dev/immune"
	"github.com/frain-dev/immune/config"
	log "github.com/sirupsen/logrus"
)

type Fire struct {
	cfg             *config.Config
	jwtToken        string
	ProjectApiKey   string
	ProjectID       string
	orgID           string
	endpointIDs     []string
	endpointOwnerID string
}

func NewFire(cfg *config.Config) *Fire {
	return &Fire{cfg: cfg}
}

func (f *Fire) Start(ctx context.Context) (*Log, error) {
	err := f.Init()
	if err != nil {
		return nil, fmt.Errorf("initialization failed: %v", err)
	}

	r := Request{
		contentType: contentType,
		url:         fmt.Sprintf("%s/api/v1/projects/%s/events", f.cfg.ConvoyURL, f.ProjectID),
		method:      immune.MethodPost,
		headers:     http.Header{"Authorization": []string{fmt.Sprintf("Bearer %s", f.ProjectApiKey)}},
	}

	err = r.WithJSONBody(f.getRequestBody())

	if err != nil {
		return nil, fmt.Errorf("failed to marshal event: %v", err)
	}

	var (
		l     = NewLog(f.cfg.Events)
		event = &datastore.Event{}
		c     = time.Now()
		resp  *Response
		now   time.Time
	)

	bar := progressbar.Default(f.cfg.Events)
	for i := int64(0); i < f.cfg.Events; i++ {
		_ = bar.Add(1)

		now = time.Now()
		resp, err = r.SendRequest(ctx)
		if err != nil {
			l.Failures++ // the request failed, record the status code
			l.FailureCodes[resp.statusCode]++
			log.WithError(err).Error("send event request failed")
			continue
		}
		l.RequestDurations = append(l.RequestDurations, time.Since(now).Milliseconds()) // record request duration

		l.EventsSent++
		err = resp.DecodeJSON(event)
		if err != nil {
			l.Failures++
			l.FailureCodes[resp.statusCode]++
			log.WithError(err).Error("failed to decode event json body, marked as failure")
			continue
		}

		l.EventTime[event.UID] = now.UTC().Format(time.RFC3339)
	}

	l.TotalTimeTaken = fmt.Sprintf("%f minutes", time.Since(c).Minutes())
	l.CalculateStats()

	return l, nil
}

func (f *Fire) getRequestBody() interface{} {
	switch f.cfg.TestType {
	case config.FanOutTest:
		return &models.FanoutEvent{
			OwnerID:   f.endpointOwnerID,
			EventType: "immune.test",
			Data:      payload,
		}
	default:
		return &models.CreateEvent{
			EndpointID: f.endpointIDs[0],
			EventType:  "immune.test",
			Data:       payload,
		}

	}
}
