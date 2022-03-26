package exec

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/frain-dev/immune"
	"github.com/frain-dev/immune/mocks"
	"github.com/golang/mock/gomock"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestExecutor_ExecuteSetupTestCase(t *testing.T) {
	ex := NewExecutor(nil, http.DefaultClient, nil, 10, "http://localhost:5005", "data", nil, nil)

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
				vm: immune.NewVariableMap(),
			},
			args: args{
				ctx: context.Background(),
				setupTC: &immune.SetupTestCase{
					Name: "abc",
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
					StatusCode:   http.StatusOK,
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
			name: "should_execute_setup_test_case_with_no_request_body",
			fields: fields{
				vm: immune.NewVariableMap(),
			},
			args: args{
				ctx: context.Background(),
				setupTC: &immune.SetupTestCase{
					Name: "abc",
					StoreResponseVariables: immune.S{
						"user_id": "user_id",
					},
					ResponseBody: true,
					Endpoint:     "/create_user",
					HTTPMethod:   "POST",
					StatusCode:   http.StatusOK,
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
			name: "should_error_for_wrong_status_code",
			fields: fields{
				vm: immune.NewVariableMap(),
			},
			args: args{
				ctx: context.Background(),
				setupTC: &immune.SetupTestCase{
					Name: "abc",
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
					StatusCode:   http.StatusUnauthorized,
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
			wantErrMsg: `setup_test_case abc: wants status code 401 but got status code 200, response body: {"user_id":"1223-242-2322"}`,
			wantErr:    true,
		},
		{
			name: "should_error_for_url_variable_not_found_in_variable_map",
			fields: fields{
				vm: immune.NewVariableMap(),
			},
			args: args{
				ctx: context.Background(),
				setupTC: &immune.SetupTestCase{
					Name:                   "abc",
					StoreResponseVariables: nil,
					RequestBody: immune.M{
						"username": "dan",
						"email":    "daniel@gmail.com",
						"phone":    113234294,
					},
					ResponseBody: true,
					Endpoint:     "/update_user/{user_id}",
					HTTPMethod:   "PUT",
					StatusCode:   http.StatusOK,
				},
			},
			arrangeFn:       nil,
			wantVariableMap: nil,
			wantErrMsg:      "setup_test_case abc: failed to process parsed url with variable map: variable user_id not found in variable map",
			wantErr:         true,
		},
		{
			name: "should_error_for_request_body_variable_not_found_in_variable_map",
			fields: fields{
				vm: immune.NewVariableMap(),
			},
			args: args{
				ctx: context.Background(),
				setupTC: &immune.SetupTestCase{
					Name:                   "abc",
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
					StatusCode:   http.StatusOK,
				},
			},
			arrangeFn:       nil,
			wantVariableMap: nil,
			wantErrMsg:      "setup_test_case abc: failed to process request body with variable map: variable user_id does not exist in variable map",
			wantErr:         true,
		},
		{
			name: "should_error_for_no_response_body",
			fields: fields{
				vm: immune.NewVariableMap(),
			},
			args: args{
				ctx: context.Background(),
				setupTC: &immune.SetupTestCase{
					Name: "abc",
					RequestBody: immune.M{
						"username": "daniel",
						"email":    "daniel@gmail.com",
						"phone":    113234294,
					},
					ResponseBody: true,
					Endpoint:     "/create_user",
					HTTPMethod:   "POST",
					StatusCode:   http.StatusOK,
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
			wantErrMsg: "setup_test_case abc: wants response body but got no response body",
			wantErr:    true,
		},
		{
			name: "should_error_for_invalid_response_body",
			fields: fields{
				vm: immune.NewVariableMap(),
			},
			args: args{
				ctx: context.Background(),
				setupTC: &immune.SetupTestCase{
					Name: "abc",
					RequestBody: immune.M{
						"username": "daniel",
						"email":    "daniel@gmail.com",
						"phone":    113234294,
					},
					ResponseBody: true,
					Endpoint:     "/create_user",
					HTTPMethod:   "POST",
					StatusCode:   http.StatusOK,
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
			wantErrMsg: "setup_test_case abc: failed to decode response body: response body: : json: cannot unmarshal number into Go value of type immune.M",
			wantErr:    true,
		},
		{
			name: "should_error_for_variable_field_not_found",
			fields: fields{
				vm: immune.NewVariableMap(),
			},
			args: args{
				ctx: context.Background(),
				setupTC: &immune.SetupTestCase{
					Name: "abc",
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
					StatusCode:   http.StatusOK,
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
			wantErrMsg: "setup_test_case abc: failed to process response body: response body: : field user_id: not found",
			wantErr:    true,
		},
		{
			name: "should_error_for_has_response_body",
			fields: fields{
				vm: immune.NewVariableMap(),
			},
			args: args{
				ctx: context.Background(),
				setupTC: &immune.SetupTestCase{
					Name: "abc",
					RequestBody: immune.M{
						"username": "daniel",
						"email":    "daniel@gmail.com",
						"phone":    113234294,
					},
					ResponseBody: false,
					Endpoint:     "/create_user",
					HTTPMethod:   "POST",
					StatusCode:   http.StatusOK,
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
			wantErrMsg: "setup_test_case abc: does not want a response body but got a response body: '{}'",
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
	ex := NewExecutor(nil, http.DefaultClient, nil, 10, "http://localhost:5005", "data", nil, nil)
	type fields struct {
		vm *immune.VariableMap
	}

	type args struct {
		ctx context.Context
		tc  *immune.TestCase
	}
	tests := []struct {
		name       string
		idFn       func() string
		fields     fields
		args       args
		arrangeFn  func(server *mocks.MockCallbackServer, tr *mocks.MockTruncator) func()
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "should_execute_test_case",
			fields: fields{
				vm: &immune.VariableMap{
					VariableToValue: immune.M{
						"user_id":      "1234",
						"company_name": "abc",
					},
				},
			},
			args: args{
				ctx: context.Background(),
				tc: &immune.TestCase{
					Name:         "abc",
					Setup:        nil,
					StatusCode:   200,
					HTTPMethod:   "POST",
					Endpoint:     "/update_user/{user_id}",
					ResponseBody: true,
					Callback: immune.Callback{
						Enabled: true,
						Times:   2,
					},
					RequestBody: immune.M{
						"email":        "dan@gmail.com",
						"phone":        23453530833,
						"data":         map[string]interface{}{},
						"company_name": "{company_name}",
					},
				},
			},
			idFn: func() string {
				return "12345"
			},
			arrangeFn: func(server *mocks.MockCallbackServer, tr *mocks.MockTruncator) func() {
				var rc chan<- *immune.Signal
				server.EXPECT().ReceiveCallback(gomock.AssignableToTypeOf(rc)).Times(2).DoAndReturn(func(c chan<- *immune.Signal) {
					c <- &immune.Signal{ImmuneCallBackID: "12345"}
				})

				tr.EXPECT().Truncate(gomock.Any()).Times(1)
				httpmock.Activate()

				httpmock.RegisterResponder(http.MethodPost, "http://localhost:5005/update_user/1234",
					httpmock.NewStringResponder(http.StatusOK, `{"user":{"username":"daniel"}}`))

				return func() {
					httpmock.DeactivateAndReset()
				}
			},
			wantErr: false,
		},
		{
			name: "should_execute_test_case_with_no_request_body",
			fields: fields{
				vm: &immune.VariableMap{
					VariableToValue: immune.M{
						"user_id":      "1234",
						"company_name": "abc",
					},
				},
			},
			args: args{
				ctx: context.Background(),
				tc: &immune.TestCase{
					Name:         "abc",
					Setup:        nil,
					StatusCode:   200,
					HTTPMethod:   "POST",
					Endpoint:     "/update_user/{user_id}",
					ResponseBody: true,
				},
			},
			idFn: func() string {
				return "12345"
			},
			arrangeFn: func(server *mocks.MockCallbackServer, tr *mocks.MockTruncator) func() {
				tr.EXPECT().Truncate(gomock.Any()).Times(1)
				httpmock.Activate()

				httpmock.RegisterResponder(http.MethodPost, "http://localhost:5005/update_user/1234",
					httpmock.NewStringResponder(http.StatusOK, `{"user":{"username":"daniel"}}`))

				return func() {
					httpmock.DeactivateAndReset()
				}
			},
			wantErr: false,
		},
		{
			name: "should_error_for_invalid_callback_body",
			fields: fields{
				vm: &immune.VariableMap{
					VariableToValue: immune.M{
						"company_name": "abc",
					},
				},
			},
			args: args{
				ctx: context.Background(),
				tc: &immune.TestCase{
					Name:         "abc",
					Setup:        nil,
					StatusCode:   200,
					HTTPMethod:   "POST",
					Endpoint:     "/update",
					ResponseBody: true,
					Callback: immune.Callback{
						Enabled: true,
						Times:   2,
					},
					RequestBody: immune.M{
						"email":        "dan@gmail.com",
						"phone":        23453530833,
						"data":         map[string]interface{}{},
						"company_name": "{company_name}",
					},
				},
			},
			idFn: func() string {
				return "12345"
			},
			arrangeFn: func(server *mocks.MockCallbackServer, tr *mocks.MockTruncator) func() {
				var rc chan<- *immune.Signal
				server.EXPECT().ReceiveCallback(gomock.AssignableToTypeOf(rc)).Times(1).DoAndReturn(func(c chan<- *immune.Signal) {
					c <- &immune.Signal{Err: errors.New("failed to decode callback body")}
				})
				httpmock.Activate()

				httpmock.RegisterResponder(http.MethodPost, "http://localhost:5005/update",
					httpmock.NewStringResponder(http.StatusOK, `{"user":{"username":"daniel"}}`))

				return func() {
					httpmock.DeactivateAndReset()
				}
			},
			wantErr:    true,
			wantErrMsg: "test_case abc: callback error: failed to decode callback body",
		},
		{
			name: "should_error_for_url_variable_not_found_in_variable_map",
			fields: fields{
				vm: &immune.VariableMap{
					VariableToValue: immune.M{
						"company_name": "abc",
					},
				},
			},
			args: args{
				ctx: context.Background(),
				tc: &immune.TestCase{
					Name:         "abc",
					Setup:        nil,
					StatusCode:   200,
					HTTPMethod:   "POST",
					Endpoint:     "/update_user/{user_id}",
					ResponseBody: true,
					Callback: immune.Callback{
						Enabled: true,
						Times:   2,
					},
					RequestBody: immune.M{
						"email":        "dan@gmail.com",
						"phone":        23453530833,
						"data":         map[string]interface{}{},
						"company_name": "{company_name}",
					},
				},
			},
			arrangeFn: func(server *mocks.MockCallbackServer, tr *mocks.MockTruncator) func() {
				httpmock.Activate()

				httpmock.RegisterResponder(http.MethodPost, "http://localhost:5005/update_user/1234",
					httpmock.NewStringResponder(http.StatusOK, `{"user":{"username":"daniel"}}`))

				return func() {
					httpmock.DeactivateAndReset()
				}
			},
			wantErr:    true,
			wantErrMsg: "test_case abc: failed to process parsed url with variable map: variable user_id not found in variable map",
		},
		{
			name: "should_error_for_callback_id_field_not_found",
			fields: fields{
				vm: &immune.VariableMap{
					VariableToValue: immune.M{
						"company_name": "abc",
					},
				},
			},
			args: args{
				ctx: context.Background(),
				tc: &immune.TestCase{
					Name:         "abc",
					Setup:        nil,
					StatusCode:   200,
					HTTPMethod:   "POST",
					Endpoint:     "/update_user",
					ResponseBody: true,
					Callback: immune.Callback{
						Enabled: true,
						Times:   2,
					},
					RequestBody: immune.M{
						"email":        "dan@gmail.com",
						"phone":        23453530833,
						"company_name": "{company_name}",
					},
				},
			},
			idFn: func() string {
				return "12345"
			},
			arrangeFn: func(server *mocks.MockCallbackServer, tr *mocks.MockTruncator) func() {
				httpmock.Activate()

				httpmock.RegisterResponder(http.MethodPost, "http://localhost:5005/update_user",
					httpmock.NewStringResponder(http.StatusOK, `{"user":{"username":"daniel"}}`))

				return func() {
					httpmock.DeactivateAndReset()
				}
			},
			wantErr:    true,
			wantErrMsg: "test_case abc: failed to inject callback id into request body: the field data, does not exist",
		},
		{
			name: "should_error_for_company_name_variable_not_found",
			fields: fields{
				vm: &immune.VariableMap{
					VariableToValue: immune.M{
						"user_id": "1234",
					},
				},
			},
			args: args{
				ctx: context.Background(),
				tc: &immune.TestCase{
					Name:         "abc",
					Setup:        nil,
					StatusCode:   200,
					HTTPMethod:   "POST",
					Endpoint:     "/update_user",
					ResponseBody: true,
					Callback: immune.Callback{
						Enabled: true,
						Times:   2,
					},
					RequestBody: immune.M{
						"email":        "dan@gmail.com",
						"phone":        23453530833,
						"data":         map[string]interface{}{},
						"company_name": "{company_name}",
					},
				},
			},
			idFn: func() string {
				return "12345"
			},
			arrangeFn: func(server *mocks.MockCallbackServer, tr *mocks.MockTruncator) func() {
				httpmock.Activate()

				httpmock.RegisterResponder(http.MethodPost, "http://localhost:5005/update_user",
					httpmock.NewStringResponder(http.StatusOK, `{"user":{"username":"daniel"}}`))

				return func() {
					httpmock.DeactivateAndReset()
				}
			},
			wantErr:    true,
			wantErrMsg: "test_case abc: failed to process request body with variable map: variable company_name does not exist in variable map",
		},
		{
			name: "should_error_for_wrong_status_code",
			fields: fields{
				vm: &immune.VariableMap{
					VariableToValue: immune.M{
						"user_id": "1234",
					},
				},
			},
			args: args{
				ctx: context.Background(),
				tc: &immune.TestCase{
					Name:         "abc",
					Setup:        nil,
					StatusCode:   201,
					HTTPMethod:   "POST",
					Endpoint:     "/update_user",
					ResponseBody: true,
					Callback: immune.Callback{
						Enabled: true,
						Times:   2,
					},
					RequestBody: immune.M{
						"email": "dan@gmail.com",
						"phone": 23453530833,
						"data":  map[string]interface{}{},
					},
				},
			},
			idFn: func() string {
				return "12345"
			},
			arrangeFn: func(server *mocks.MockCallbackServer, tr *mocks.MockTruncator) func() {
				httpmock.Activate()

				httpmock.RegisterResponder(http.MethodPost, "http://localhost:5005/update_user",
					httpmock.NewStringResponder(http.StatusOK, `{"user":{"username":"daniel"}}`))

				return func() {
					httpmock.DeactivateAndReset()
				}
			},
			wantErr:    true,
			wantErrMsg: "test_case abc: wants status code 201 but got status code 200",
		},
		{
			name: "should_error_for_got_no_response_body",
			fields: fields{
				vm: &immune.VariableMap{
					VariableToValue: immune.M{
						"user_id": "1234",
					},
				},
			},
			args: args{
				ctx: context.Background(),
				tc: &immune.TestCase{
					Name:         "abc",
					Setup:        nil,
					StatusCode:   405,
					HTTPMethod:   "POST",
					Endpoint:     "/update_user",
					ResponseBody: true,
					Callback: immune.Callback{
						Enabled: true,
						Times:   2,
					},
					RequestBody: immune.M{
						"email": "dan@gmail.com",
						"phone": 23453530833,
						"data":  map[string]interface{}{},
					},
				},
			},
			idFn: func() string {
				return "12345"
			},
			arrangeFn: func(server *mocks.MockCallbackServer, tr *mocks.MockTruncator) func() {
				httpmock.Activate()

				httpmock.RegisterResponder(http.MethodPost, "http://localhost:5005/update_user",
					httpmock.NewStringResponder(http.StatusMethodNotAllowed, ``))

				return func() {
					httpmock.DeactivateAndReset()
				}
			},
			wantErr:    true,
			wantErrMsg: "test_case abc: wants response body but got no response body: status_code: 405",
		},
		{
			name: "should_error_for_invalid_response_body",
			fields: fields{
				vm: &immune.VariableMap{
					VariableToValue: immune.M{
						"user_id": "1234",
					},
				},
			},
			args: args{
				ctx: context.Background(),
				tc: &immune.TestCase{
					Name:         "abc",
					Setup:        nil,
					StatusCode:   200,
					HTTPMethod:   "POST",
					Endpoint:     "/update_user",
					ResponseBody: true,
					Callback: immune.Callback{
						Enabled: true,
						Times:   2,
					},
					RequestBody: immune.M{
						"email": "dan@gmail.com",
						"phone": 23453530833,
						"data":  map[string]interface{}{},
					},
				},
			},
			idFn: func() string {
				return "12345"
			},
			arrangeFn: func(server *mocks.MockCallbackServer, tr *mocks.MockTruncator) func() {
				httpmock.Activate()

				httpmock.RegisterResponder(http.MethodPost, "http://localhost:5005/update_user",
					httpmock.NewStringResponder(http.StatusOK, `123456`))

				return func() {
					httpmock.DeactivateAndReset()
				}
			},
			wantErr:    true,
			wantErrMsg: "test_case abc: failed to decode response body: 123456: json: cannot unmarshal number into Go value of type immune.M",
		},
		{
			name: "should_error_for_got_response_body",
			fields: fields{
				vm: &immune.VariableMap{
					VariableToValue: immune.M{
						"user_id": "1234",
					},
				},
			},
			args: args{
				ctx: context.Background(),
				tc: &immune.TestCase{
					Name:         "abc",
					Setup:        nil,
					StatusCode:   200,
					HTTPMethod:   "POST",
					Endpoint:     "/update_user",
					ResponseBody: false,
					Callback: immune.Callback{
						Enabled: true,
						Times:   2,
					},
					RequestBody: immune.M{
						"email": "dan@gmail.com",
						"phone": 23453530833,
						"data":  map[string]interface{}{},
					},
				},
			},
			idFn: func() string {
				return "12345"
			},
			arrangeFn: func(server *mocks.MockCallbackServer, tr *mocks.MockTruncator) func() {
				httpmock.Activate()

				httpmock.RegisterResponder(http.MethodPost, "http://localhost:5005/update_user",
					httpmock.NewStringResponder(http.StatusOK, `123456`))

				return func() {
					httpmock.DeactivateAndReset()
				}
			},
			wantErr:    true,
			wantErrMsg: "test_case abc: does not want a response body but got a response body: '123456'",
		},
		{
			name: "should_error_for_incorrect_callback_id",
			fields: fields{
				vm: &immune.VariableMap{
					VariableToValue: immune.M{
						"user_id": "1234",
					},
				},
			},
			args: args{
				ctx: context.Background(),
				tc: &immune.TestCase{
					Name:         "abc",
					Setup:        nil,
					StatusCode:   200,
					HTTPMethod:   "POST",
					Endpoint:     "/update_user",
					ResponseBody: true,
					Callback: immune.Callback{
						Enabled: true,
						Times:   2,
					},
					RequestBody: immune.M{
						"email": "dan@gmail.com",
						"phone": 23453530833,
						"data":  map[string]interface{}{},
					},
				},
			},
			idFn: func() string {
				return "12345"
			},
			arrangeFn: func(server *mocks.MockCallbackServer, tr *mocks.MockTruncator) func() {
				var rc chan<- *immune.Signal
				server.EXPECT().ReceiveCallback(gomock.AssignableToTypeOf(rc)).Times(1).Do(func(rc chan<- *immune.Signal) {
					rc <- &immune.Signal{ImmuneCallBackID: "1234"}
				})
				httpmock.Activate()

				httpmock.RegisterResponder(http.MethodPost, "http://localhost:5005/update_user",
					httpmock.NewStringResponder(http.StatusOK, `{"immune_callback_id":"12345"}`))

				return func() {
					httpmock.DeactivateAndReset()
				}
			},
			wantErr:    true,
			wantErrMsg: "test_case abc: incorrect callback_id: expected_callback_id '12345', got_callback_id '1234'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockDBTruncator := mocks.NewMockTruncator(ctrl)

			mockCallbackServer := mocks.NewMockCallbackServer(ctrl)
			if tt.arrangeFn != nil {
				deferFn := tt.arrangeFn(mockCallbackServer, mockDBTruncator)
				defer deferFn()
			}

			ex.s = mockCallbackServer
			ex.dbTruncator = mockDBTruncator
			ex.vm = tt.fields.vm
			ex.idFn = tt.idFn
			err := ex.ExecuteTestCase(tt.args.ctx, tt.args.tc)
			if tt.wantErr {
				require.Error(t, err)
				require.Equal(t, tt.wantErrMsg, err.Error())
				return
			}

			require.NoError(t, err)
		})
	}
}
