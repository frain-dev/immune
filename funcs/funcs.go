package funcs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/frain-dev/immune"
	"github.com/frain-dev/immune/exec"
	"github.com/google/uuid"
)

func SetupGroup(ctx context.Context, ex *exec.Executor) error {
	req := `{
                "config": {
                    "disableEndpoint": true,
                    "signature": {
                        "hash": "SHA256",
                        "header": "X-Retro-Signature"
                    },
                    "strategy": {
                        "default": {
                            "intervalSeconds": 30,
                            "retryLimit": 4
                        },
                        "type": "default"
                    }
                },
                "logo_url": "",
                "name": "immune-group-%s"
            }`

	req = fmt.Sprintf(req, uuid.New().String())
	mapper := map[string]interface{}{}
	err := json.Unmarshal([]byte(req), &mapper)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %v", err)
	}

	tc := &immune.SetupTestCase{
		Name: "setup_group",
		StoreResponseVariables: immune.S{
			"group_id": "data.uid",
		},
		RequestBody:  mapper,
		ResponseBody: true,
		Endpoint:     "/groups",
		HTTPMethod:   "POST",
		StatusCode:   201,
	}

	return ex.ExecuteSetupTestCase(ctx, tc)
}

func SetupApp(ctx context.Context, ex *exec.Executor) error {
	const req = `{
             "name": "retro-app-%s",
			 "support_email": "retro_app-%s@gmail.com"
            }`

	uid := uuid.New().String()
	r := fmt.Sprintf(req, uid, uid)

	mapper := map[string]interface{}{}
	err := json.Unmarshal([]byte(r), &mapper)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %v", err)
	}

	tc := &immune.SetupTestCase{
		Name: "setup_app",
		StoreResponseVariables: immune.S{
			"app_id": "data.uid",
		},
		RequestBody:  mapper,
		ResponseBody: true,
		Endpoint:     "/applications?groupId={group_id}",
		HTTPMethod:   "POST",
		StatusCode:   201,
	}

	return ex.ExecuteSetupTestCase(ctx, tc)
}

func SetupAppEndpoint(ctx context.Context, targetURL string, ex *exec.Executor) error {
	req := `{
             "url": "%s",
                "secret": "12345",
                "description": "Local ngrok endpoint",
                "events": [
                    "payment.failed"
                ]
            }`

	req = fmt.Sprintf(req, targetURL)
	mapper := map[string]interface{}{}
	err := json.Unmarshal([]byte(req), &mapper)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %v", err)
	}

	tc := &immune.SetupTestCase{
		Name: "setup_endpoint",
		StoreResponseVariables: immune.S{
			"endpoint_id": "data.uid",
		},
		RequestBody:  mapper,
		ResponseBody: true,
		Endpoint:     "/applications/{app_id}/endpoints?groupId={group_id}",
		HTTPMethod:   "POST",
		StatusCode:   201,
	}

	return ex.ExecuteSetupTestCase(ctx, tc)
}

func SetupEvent(ctx context.Context, ex *exec.Executor) error {
	// uses a random event type so convoy will discard the event
	req := `{
                "app_id": "{app_id}",
                "event_type": "%s",
                "data": {
                    "sc": "gene",
                    "marvel": "stark"
                }
            }`

	req = fmt.Sprintf(req, uuid.New().String())
	mapper := map[string]interface{}{}
	err := json.Unmarshal([]byte(req), &mapper)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %v", err)
	}

	tc := &immune.SetupTestCase{
		Name: "setup_event",
		StoreResponseVariables: immune.S{
			"event_id": "data.uid",
		},
		RequestBody:  mapper,
		ResponseBody: true,
		Endpoint:     "/events?groupId={group_id}",
		HTTPMethod:   "POST",
		StatusCode:   201,
	}

	return ex.ExecuteSetupTestCase(ctx, tc)
}
