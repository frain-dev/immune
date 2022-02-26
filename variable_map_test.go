package immune

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVariableMap_ProcessResponse(t *testing.T) {
	type fields struct {
		VariableToValue M
	}
	type args struct {
		ctx             context.Context
		variableToField S
		values          M
	}
	tests := []struct {
		name            string
		args            args
		wantVariableMap M
		wantErr         bool
		wantErrMsg      string
	}{
		{
			name: "should_process_response_successfully",
			args: args{
				ctx: context.Background(),
				variableToField: S{
					"num_events": "data.metadata.num_events",
					"message":    "message",
					"app_id":     "data.uid",
				},
				values: M{
					"status":  true,
					"message": "fetched application successfully",
					"data": map[string]interface{}{
						"uid": "138-343-12132-4245",
						"metadata": map[string]interface{}{
							"num_events": 23230,
						},
						"title": "retro-app",
					},
				},
			},
			wantVariableMap: M{
				"num_events": 23230,
				"message":    "fetched application successfully",
				"app_id":     "138-343-12132-4245",
			},
			wantErr: false,
		},
		{
			name: "should_error_for_unsupported_type",
			args: args{
				ctx: context.Background(),
				variableToField: S{
					"status":  "status",
					"message": "message",
					"app_id":  "data.uid",
				},
				values: M{
					"status":  true,
					"message": "fetched application successfully",
					"data": map[string]interface{}{
						"uid":   "138-343-12132-4245",
						"title": "retro-app",
					},
				},
			},
			wantVariableMap: M{},
			wantErrMsg:      "variable status is of type bool in the response body, only string & integer is currently supported",
			wantErr:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &VariableMap{VariableToValue: M{}}

			err := v.ProcessResponse(tt.args.ctx, tt.args.variableToField, tt.args.values)
			if tt.wantErr {
				require.Error(t, err)
				require.Equal(t, err.Error(), tt.wantErrMsg)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantVariableMap, v.VariableToValue)
		})
	}
}

func TestVariableMap_GetString(t *testing.T) {
	type fields struct {
		VariableToValue M
	}
	type args struct {
		key string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		want       string
		wantExists bool
	}{
		{
			name: "should_get_key_value_successfully",
			fields: fields{
				VariableToValue: M{
					"app_id": "12345678",
				},
			},
			args: args{
				key: "app_id",
			},
			want:       "12345678",
			wantExists: true,
		},
		{
			name: "should_get_key_not_exists",
			fields: fields{
				VariableToValue: M{
					"app_id": "12345678",
				},
			},
			args: args{
				key: "group_id",
			},
			want:       "",
			wantExists: false,
		},
		{
			name: "should_get_int_value_in_string_format",
			fields: fields{
				VariableToValue: M{
					"app_id": 12345678,
				},
			},
			args: args{
				key: "app_id",
			},
			want:       "12345678",
			wantExists: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := VariableMap{
				VariableToValue: tt.fields.VariableToValue,
			}
			got, exists := v.GetString(tt.args.key)
			require.Equal(t, tt.want, got)
			require.Equal(t, tt.wantExists, exists)
		})
	}
}

func TestVariableMap_Get(t *testing.T) {
	type fields struct {
		VariableToValue M
	}
	type args struct {
		key string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		want       interface{}
		wantExists bool
	}{
		{
			name: "should_get_value",
			fields: fields{
				VariableToValue: M{
					"group_id": "12345678",
				},
			},
			args: args{
				key: "group_id",
			},
			want:       "12345678",
			wantExists: true,
		},
		{
			name: "should_get_key_not_exists",
			fields: fields{
				VariableToValue: M{
					"group_id": "12345678",
				},
			},
			args: args{
				key: "app_id",
			},
			want:       nil,
			wantExists: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := VariableMap{
				VariableToValue: tt.fields.VariableToValue,
			}
			got, exists := v.Get(tt.args.key)
			require.Equal(t, tt.want, got)
			require.Equal(t, tt.wantExists, exists)
		})
	}
}
