package exec

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"

	"github.com/frain-dev/immune"
)

func TestExecutor_ExecuteSetupTestCase(t *testing.T) {
	ex := NewExecutor(nil, http.DefaultClient, nil, 10, "http://localhost:5005", "data")

	emptyVM := func() *immune.VariableMap {
		return &immune.VariableMap{VariableToValue: immune.M{}}
	}
	type fields struct {
		vm *immune.VariableMap
	}
	type args struct {
		ctx     context.Context
		setupTC *immune.SetupTestCase
	}
	tests := []struct {
		name            string
		fields          fields
		arrangeFn       func() func()
		args            args
		wantVariableMap *immune.VariableMap
		wantErrMsg      string
		wantErr         bool
	}{
		{
			name: "should_execute_setup_test_case",
			fields: fields{
				vm: emptyVM(),
			},
			args: args{
				ctx: context.Background(),
				setupTC: &immune.SetupTestCase{
					StoreResponseVariables: immune.S{
						"user_id": "user_id",
					},
					RequestBody: immune.M{
						"username": "daniel",
						"email":    "daniel@gmail.com",
						"phone":    113234294,
					},
					ResponseBody: true,
					Endpoint:     "/create_user",
					HTTPMethod:   "POST",
					Position:     1,
				},
			},
			arrangeFn: func() func() {
				httpmock.Activate()

				httpmock.RegisterResponder(http.MethodPost, "http://localhost:5005/create_user",
					httpmock.NewStringResponder(http.StatusOK, `{"user_id":"1223-242-2322"}`))

				return func() {
					httpmock.DeactivateAndReset()
				}
			},
			wantVariableMap: &immune.VariableMap{
				VariableToValue: immune.M{"user_id": "1223-242-2322"},
			},
			wantErrMsg: "",
			wantErr:    false,
		},
		{
			name: "should_error_for_url_variable_not_found_in_variable_map",
			fields: fields{
				vm: emptyVM(),
			},
			args: args{
				ctx: context.Background(),
				setupTC: &immune.SetupTestCase{
					StoreResponseVariables: nil,
					RequestBody: immune.M{
						"username": "dan",
						"email":    "daniel@gmail.com",
						"phone":    113234294,
					},
					ResponseBody: true,
					Endpoint:     "/update_user/{user_id}",
					HTTPMethod:   "PUT",
					Position:     1,
				},
			},
			arrangeFn:       nil,
			wantVariableMap: nil,
			wantErrMsg:      "setup_test_case 1: failed to process parsed url with variable map: variable user_id not found in variable map",
			wantErr:         true,
		},
		{
			name: "should_error_for_request_body_variable_not_found_in_variable_map",
			fields: fields{
				vm: emptyVM(),
			},
			args: args{
				ctx: context.Background(),
				setupTC: &immune.SetupTestCase{
					StoreResponseVariables: nil,
					RequestBody: immune.M{
						"user_details": map[string]interface{}{
							"username": "{user_id}",
							"email":    "daniel@gmail.com",
							"phone":    113234294,
						}},
					ResponseBody: true,
					Endpoint:     "/create_user",
					HTTPMethod:   "POST",
					Position:     1,
				},
			},
			arrangeFn:       nil,
			wantVariableMap: nil,
			wantErrMsg:      "setup_test_case 1: failed to process request body with variable map: variable user_id does not exist in variable map",
			wantErr:         true,
		},
		{
			name: "should_error_for_no_response_body",
			fields: fields{
				vm: emptyVM(),
			},
			args: args{
				ctx: context.Background(),
				setupTC: &immune.SetupTestCase{
					RequestBody: immune.M{
						"username": "daniel",
						"email":    "daniel@gmail.com",
						"phone":    113234294,
					},
					ResponseBody: true,
					Endpoint:     "/create_user",
					HTTPMethod:   "POST",
					Position:     1,
				},
			},
			arrangeFn: func() func() {
				httpmock.Activate()

				httpmock.RegisterResponder(http.MethodPost, "http://localhost:5005/create_user",
					httpmock.NewStringResponder(http.StatusOK, ""))

				return func() {
					httpmock.DeactivateAndReset()
				}
			},
			wantErrMsg: "setup_test_case 1: wants response body but got no response body",
			wantErr:    true,
		},
		{
			name: "should_error_for_invalid_response_body",
			fields: fields{
				vm: emptyVM(),
			},
			args: args{
				ctx: context.Background(),
				setupTC: &immune.SetupTestCase{
					RequestBody: immune.M{
						"username": "daniel",
						"email":    "daniel@gmail.com",
						"phone":    113234294,
					},
					ResponseBody: true,
					Endpoint:     "/create_user",
					HTTPMethod:   "POST",
					Position:     1,
				},
			},
			arrangeFn: func() func() {
				httpmock.Activate()

				httpmock.RegisterResponder(http.MethodPost, "http://localhost:5005/create_user",
					httpmock.NewStringResponder(http.StatusOK, "3242424"))

				return func() {
					httpmock.DeactivateAndReset()
				}
			},
			wantErrMsg: "setup_test_case 1: failed to decode response body: json: cannot unmarshal number into Go value of type immune.M",
			wantErr:    true,
		},
		{
			name: "should_error_for_variable_field_not_found",
			fields: fields{
				vm: emptyVM(),
			},
			args: args{
				ctx: context.Background(),
				setupTC: &immune.SetupTestCase{
					StoreResponseVariables: immune.S{
						"user_id": "user_id",
					},
					RequestBody: immune.M{
						"username": "daniel",
						"email":    "daniel@gmail.com",
						"phone":    113234294,
					},
					ResponseBody: true,
					Endpoint:     "/create_user",
					HTTPMethod:   "POST",
					Position:     1,
				},
			},
			arrangeFn: func() func() {
				httpmock.Activate()

				httpmock.RegisterResponder(http.MethodPost, "http://localhost:5005/create_user",
					httpmock.NewStringResponder(http.StatusOK, `{"user":{"username":"daniel"}}`))

				return func() {
					httpmock.DeactivateAndReset()
				}
			},
			wantErrMsg: "setup_test_case 1: failed to process response body: field user_id does not exist",
			wantErr:    true,
		},
		{
			name: "should_error_for_has_response_body",
			fields: fields{
				vm: emptyVM(),
			},
			args: args{
				ctx: context.Background(),
				setupTC: &immune.SetupTestCase{
					RequestBody: immune.M{
						"username": "daniel",
						"email":    "daniel@gmail.com",
						"phone":    113234294,
					},
					ResponseBody: false,
					Endpoint:     "/create_user",
					HTTPMethod:   "POST",
					Position:     1,
				},
			},
			arrangeFn: func() func() {
				httpmock.Activate()

				httpmock.RegisterResponder(http.MethodPost, "http://localhost:5005/create_user",
					httpmock.NewStringResponder(http.StatusOK, "{}"))

				return func() {
					httpmock.DeactivateAndReset()
				}
			},
			wantErrMsg: "setup_test_case 1: does not want a response body but got a response body: '{}'",
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.arrangeFn != nil {
				deferFn := tt.arrangeFn()
				defer deferFn()
			}

			ex.vm = tt.fields.vm
			err := ex.ExecuteSetupTestCase(tt.args.ctx, tt.args.setupTC)
			if tt.wantErr {
				require.Error(t, err)
				require.Equal(t, tt.wantErrMsg, err.Error())
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantVariableMap, ex.vm)
		})
	}
}

func TestExecutor_ExecuteTestCase(t *testing.T) {
	type fields struct {
		callbackIDLocation     string
		baseURL                string
		maxCallbackWaitSeconds uint
		s                      immune.CallbackServer
		client                 *http.Client
		vm                     *immune.VariableMap
	}
	type args struct {
		ctx context.Context
		tc  *immune.TestCase
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		arrangeFn func() func()
		wantErr   bool
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ex := &Executor{
				callbackIDLocation:     tt.fields.callbackIDLocation,
				baseURL:                tt.fields.baseURL,
				maxCallbackWaitSeconds: tt.fields.maxCallbackWaitSeconds,
				s:                      tt.fields.s,
				client:                 tt.fields.client,
				vm:                     tt.fields.vm,
			}
			if err := ex.ExecuteTestCase(tt.args.ctx, tt.args.tc); (err != nil) != tt.wantErr {
				t.Errorf("Executor.ExecuteTestCase() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExecutor_sendRequest(t *testing.T) {
	type fields struct {
		s  immune.CallbackServer
		vm *immune.VariableMap
	}

	ex := NewExecutor(nil, http.DefaultClient, nil, 10, "", "")

	type args struct {
		ctx context.Context
		r   *request
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		arrangeFn func() func()
		want      *response
		wantErr   bool
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ex.vm = tt.fields.vm
			ex.s = tt.fields.s
			got, err := ex.sendRequest(tt.args.ctx, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Executor.sendRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Executor.sendRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
