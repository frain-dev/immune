package fire

import (
	"context"
	"fmt"
	"net/http"
	"time"

	convoyConfig "github.com/frain-dev/convoy/config"

	"github.com/frain-dev/convoy/api/models"
	"github.com/frain-dev/convoy/datastore"
	"github.com/frain-dev/immune"
	"github.com/frain-dev/immune/config"
	"github.com/oklog/ulid/v2"
)

const contentType = "application/json"

func (f *Fire) Init() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := f.loginConvoyUser(ctx)
	if err != nil {
		return fmt.Errorf("failed to login convoy user: %v", err)
	}

	err = f.createOrg(ctx)
	if err != nil {
		return fmt.Errorf("failed to create org: %v", err)
	}

	err = f.createProject(ctx)
	if err != nil {
		return fmt.Errorf("failed to create project: %v", err)
	}

	err = f.createEndpoint(ctx)
	if err != nil {
		return fmt.Errorf("failed to create endpoint: %v", err)
	}

	err = f.createSubscription(ctx)
	if err != nil {
		return fmt.Errorf("failed to create endpoint: %v", err)
	}

	return nil
}

func (f *Fire) loginConvoyUser(ctx context.Context) error {
	login := &models.LoginUser{
		Username: "superuser@default.com",
		Password: "default",
	}

	r := Request{
		contentType: contentType,
		url:         fmt.Sprintf("%s/ui/auth/login", f.cfg.ConvoyURL),
		method:      immune.MethodPost,
	}

	err := r.WithJSONBody(login)
	if err != nil {
		return fmt.Errorf("failed to add json body: %v", err)
	}

	auth := models.LoginUserResponse{}
	resp, err := r.SendRequest(ctx)
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}

	err = resp.DecodeJSON(&auth)
	if err != nil {
		return fmt.Errorf("failed to decode json body: %v", err)
	}

	f.jwtToken = auth.Token.AccessToken

	return nil
}

func (f *Fire) createOrg(ctx context.Context) error {
	r := Request{
		contentType: contentType,
		url:         fmt.Sprintf("%s/ui/organisations", f.cfg.ConvoyURL),
		method:      immune.MethodPost,
		headers:     http.Header{"Authorization": []string{fmt.Sprintf("Bearer %s", f.jwtToken)}},
	}

	err := r.WithJSONBody(models.Organisation{Name: "immune_org"})
	if err != nil {
		return fmt.Errorf("failed to add json body: %v", err)
	}

	org := datastore.Organisation{}
	resp, err := r.SendRequest(ctx)
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}

	err = resp.DecodeJSON(&org)
	if err != nil {
		return fmt.Errorf("failed to decode json body: %v", err)
	}

	f.orgID = org.UID

	return nil
}

func (f *Fire) createProject(ctx context.Context) error {
	r := Request{
		contentType: contentType,
		url:         fmt.Sprintf("%s/ui/organisations/%s/projects", f.cfg.ConvoyURL, f.orgID),
		method:      immune.MethodPost,
		headers:     http.Header{"Authorization": []string{fmt.Sprintf("Bearer %s", f.jwtToken)}},
	}

	newProject := &models.CreateProject{
		Name: "immune_project",
		Type: f.getProjectType(),
		Config: &models.ProjectConfig{
			MaxIngestSize:            2000,
			ReplayAttacks:            true,
			IsRetentionPolicyEnabled: false,
			DisableEndpoint:          false,
			RetentionPolicy: &models.RetentionPolicyConfiguration{
				Policy: "720h",
			},
			RateLimit: &models.RateLimitConfiguration{
				Count:    1000,
				Duration: 60,
			},
			Strategy: &models.StrategyConfiguration{
				Type:       "linear",
				Duration:   10,
				RetryCount: 3,
			},
			Signature: &models.SignatureConfiguration{
				Header: convoyConfig.SignatureHeaderProvider(immune.DefaultSignatureHeader),
				Versions: []models.SignatureVersion{
					{
						UID:       ulid.Make().String(),
						Hash:      immune.DefaultHash,
						Encoding:  immune.DefaultEncoding,
						CreatedAt: time.Now(),
					},
				},
			},
			MetaEvent: &models.MetaEventConfiguration{IsEnabled: false},
		},
	}

	err := r.WithJSONBody(newProject)
	if err != nil {
		return fmt.Errorf("failed to add json body: %v", err)
	}

	project := &models.CreateProjectResponse{}
	resp, err := r.SendRequest(ctx)
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}

	err = resp.DecodeJSON(project)
	if err != nil {
		return fmt.Errorf("failed to decode json body: %v", err)
	}

	f.ProjectApiKey = project.APIKey.Key
	f.ProjectID = project.Project.UID

	return nil
}

func (f *Fire) getProjectType() string {
	switch f.cfg.TestType {
	case config.IngestTest:
		return string(datastore.IncomingProject)
	case config.SingleTest, config.FanOutTest, config.DynamicTest, config.PubSubTest:
		return string(datastore.OutgoingProject)
	default:
		return ""
	}
}

func (f *Fire) createEndpoint(ctx context.Context) error {
	r := Request{
		contentType: contentType,
		url:         fmt.Sprintf("%s/api/v1/projects/%s/endpoints", f.cfg.ConvoyURL, f.ProjectID),
		method:      immune.MethodPost,
		headers:     http.Header{"Authorization": []string{fmt.Sprintf("Bearer %s", f.jwtToken)}},
	}

	newEndpoint := &models.CreateEndpoint{
		URL:                f.cfg.EndpointURL,
		Secret:             f.cfg.EndpointSecret,
		OwnerID:            "",
		Description:        "immune test endpoint",
		AdvancedSignatures: false,
		Name:               "immune-endpoint",
		SupportEmail:       "support@immune.com",
		HttpTimeout:        "10s",
	}

	err := r.WithJSONBody(newEndpoint)
	if err != nil {
		return fmt.Errorf("failed to add json body: %v", err)
	}

	endpoint := datastore.Endpoint{}

	resp, err := r.SendRequest(ctx)
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}

	err = resp.DecodeJSON(&endpoint)
	if err != nil {
		return fmt.Errorf("failed to decode json body: %v", err)
	}

	f.endpointID = endpoint.UID

	return nil
}

func (f *Fire) createSubscription(ctx context.Context) error {
	r := Request{
		contentType: contentType,
		url:         fmt.Sprintf("%s/api/v1/projects/%s/subscriptions", f.cfg.ConvoyURL, f.ProjectID),
		method:      immune.MethodPost,
		headers:     http.Header{"Authorization": []string{fmt.Sprintf("Bearer %s", f.jwtToken)}},
	}

	newSubscription := &models.CreateSubscription{
		Name:       "immune-subscription",
		EndpointID: f.endpointID,
		RetryConfig: &models.RetryConfiguration{
			Type:            datastore.LinearStrategyProvider,
			Duration:        "10s",
			IntervalSeconds: 0,
			RetryCount:      3,
		},
	}

	err := r.WithJSONBody(newSubscription)
	if err != nil {
		return fmt.Errorf("failed to add json body: %v", err)
	}

	resp, err := r.SendRequest(ctx)
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}

	err = resp.DecodeJSON(&datastore.Subscription{})
	if err != nil {
		return fmt.Errorf("failed to decode json body: %v", err)
	}

	return nil
}
