package exec

import (
	"testing"

	"github.com/frain-dev/immune"
	"github.com/stretchr/testify/require"
)

func Test_request_processWithVariableMap(t *testing.T) {
	r := &request{
		contentType: "application/json",
		url:         "http://localhost:5005",
		method:      immune.MethodGet,
	}

	type args struct {
		vm *immune.VariableMap
	}
	tests := []struct {
		name       string
		body       immune.M
		args       args
		wantErr    bool
		wantErrMsg string
		wantBody   immune.M
	}{
		{
			name: "should_process_request_body",
			body: immune.M{
				"data": map[string]interface{}{
					"name":  "daniel",
					"phone": "{phone}",
					"ref":   "{}",
					"email": "dan@gmail.com",
					"user": map[string]interface{}{
						"group_name": "{group_name}",
						"group_id":   "{group_id}",
					},
				},
			},
			args: args{
				vm: &immune.VariableMap{
					VariableToValue: immune.M{
						"phone":      90324242,
						"group_name": "red_house",
						"group_id":   "123-454-655",
					},
				},
			},
			wantBody: immune.M{
				"data": map[string]interface{}{
					"name":  "daniel",
					"phone": 90324242,
					"ref":   "{}",
					"email": "dan@gmail.com",
					"user": map[string]interface{}{
						"group_name": "red_house",
						"group_id":   "123-454-655",
					},
				},
			},
		},
		{
			name: "should_process_request_body",
			body: immune.M{
				"data": map[string]interface{}{
					"name":  "daniel",
					"email": "dan@gmail.com",
					"user": map[string]interface{}{
						"status": "{status}",
					},
				},
			},
			args: args{
				vm: &immune.VariableMap{
					VariableToValue: immune.M{
						"group_name": "red_house",
						"group_id":   "123-454-655",
					},
				},
			},
			wantErr:    true,
			wantErrMsg: "variable status does not exist in variable map",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r.body = tt.body
			err := r.processWithVariableMap(tt.args.vm)
			if tt.wantErr {
				require.Error(t, err)
				require.Equal(t, tt.wantErrMsg, err.Error())
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantBody, r.body)
		})
	}
}
