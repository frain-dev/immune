package immune

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVariableMap_ProcessResponse(t *testing.T) {
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
					"app_name":   "data.apps[1].name",
				},
				values: M{
					"status":  true,
					"message": "fetched application successfully",
					"data": map[string]interface{}{
						"uid": "138-343-12132-4245",
						"metadata": map[string]interface{}{
							"num_events": 23230,
						},
						"apps": []interface{}{
							map[string]interface{}{"name": "danny"},
							map[string]interface{}{"name": "temi"},
						},
						"title": "retro-app",
					},
				},
			},
			wantVariableMap: M{
				"num_events": 23230,
				"message":    "fetched application successfully",
				"app_id":     "138-343-12132-4245",
				"app_name":   "temi",
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
		{
			name: "should_error_for_field_not_found",
			args: args{
				ctx: context.Background(),
				variableToField: S{
					"app_id": "data.uid",
				},
				values: M{
					"status":  true,
					"message": "fetched application successfully",
					"data": map[string]interface{}{
						"title": "retro-app",
					},
				},
			},
			wantVariableMap: M{},
			wantErrMsg:      "field data.uid: not found",
			wantErr:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewVariableMap()

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
			v := &VariableMap{
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
			v := &VariableMap{
				VariableToValue: tt.fields.VariableToValue,
			}
			got, exists := v.Get(tt.args.key)
			require.Equal(t, tt.want, got)
			require.Equal(t, tt.wantExists, exists)
		})
	}
}

func Test_getM(t *testing.T) {
	type args struct {
		m     M
		parts []string
	}
	tests := []struct {
		name       string
		args       args
		want       M
		wantErrMsg string
		wantErr    bool
	}{
		{
			name: "should_get_ref_map",
			args: args{
				m: M{
					"data": map[string]interface{}{
						"ref": map[string]interface{}{
							"status": 1234,
							"marvel": "DC",
						},
					},
				},
				parts: []string{"data", "ref"},
			},
			want: map[string]interface{}{
				"status": 1234,
				"marvel": "DC",
			},
			wantErr: false,
		},
		{
			name: "should_get_array_ref_map",
			args: args{
				m: M{
					"data": map[string]interface{}{
						"ref": []interface{}{
							map[string]interface{}{
								"status": 1234,
								"marvel": "DC",
							},
							map[string]interface{}{
								"status": 320,
								"marvel": "stark",
							},
						},
					},
				},
				parts: []string{"data", "ref[1]"},
			},
			want: map[string]interface{}{
				"status": 320,
				"marvel": "stark",
			},
			wantErr: false,
		},
		{
			name: "should_get_error_for_ref_not_found",
			args: args{
				m: M{
					"data": map[string]interface{}{},
				},
				parts: []string{"data", "ref[1]"},
			},
			wantErr:    true,
			wantErrMsg: "field data.ref[1]: not found",
		},
		{
			name: "should_get_data_map",
			args: args{
				m: M{
					"data": map[string]interface{}{
						"chef": map[string]interface{}{
							"status": 1234,
						},
					},
				},
				parts: []string{"data"},
			},
			want: map[string]interface{}{
				"chef": map[string]interface{}{
					"status": 1234,
				},
			},
			wantErr: false,
		},
		{
			name: "should_error_for_field_not_exists",
			args: args{
				m: M{
					"data": map[string]interface{}{
						"ref": map[string]interface{}{},
					},
				},
				parts: []string{"data", "ref", "status"},
			},
			want:       nil,
			wantErrMsg: "field data.ref.status: not found",
			wantErr:    true,
		},
		{
			name: "should_error_for_non_object_type",
			args: args{
				m: M{
					"data": map[string]interface{}{
						"ref": map[string]interface{}{
							"status": 1234,
						},
					},
				},
				parts: []string{"data", "ref", "status"},
			},
			want:       nil,
			wantErrMsg: "field data.ref.status: required type is object but got int",
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getM(tt.args.m, tt.args.parts)
			if tt.wantErr {
				require.Error(t, err)
				require.Equal(t, tt.wantErrMsg, err.Error())
				return
			}

			require.NoError(t, err)
			require.Equal(t, got, tt.want)
		})
	}
}

func Test_getKeyInMap(t *testing.T) {
	type args struct {
		field string
		resp  M
	}
	tests := []struct {
		name       string
		args       args
		want       interface{}
		wantErrMsg string
		wantErr    bool
	}{
		{
			name: "should_get_key_value_1",
			args: args{
				field: "data.uid",
				resp: M{
					"status":  false,
					"message": "fetched app",
					"data": map[string]interface{}{
						"uid": "1234243",
					},
				},
			},
			want:    "1234243",
			wantErr: false,
		},
		{
			name: "should_get_array_key_value",
			args: args{
				field: "data[0]",
				resp: M{
					"status":  false,
					"message": "fetched app",
					"data": []interface{}{
						"abc",
						"123",
						"king",
					},
				},
			},
			want:    "abc",
			wantErr: false,
		},
		{
			name: "should_get_array_index_out_of_range",
			args: args{
				field: "data[1]",
				resp: M{
					"status":  false,
					"message": "fetched app",
					"data": []interface{}{
						"abc",
					},
				},
			},
			wantErrMsg: "field data[1]: index out of range with length 1",
			wantErr:    true,
		},
		{
			name: "should_get_key_value_2",
			args: args{
				field: "data.ref.marvel",
				resp: M{
					"status":  false,
					"message": "fetched app",
					"data": map[string]interface{}{
						"ref": map[string]interface{}{
							"marvel": []string{"stark", "steve"},
						},
					},
				},
			},
			want:    []string{"stark", "steve"},
			wantErr: false,
		},
		{
			name: "should_error_for_field_not_exists_1",
			args: args{
				field: "data",
				resp: M{
					"status":  false,
					"message": "fetched app",
				},
			},
			wantErrMsg: "field data: not found",
			wantErr:    true,
		},
		{
			name: "should_error_for_field_not_exists_2",
			args: args{
				field: "data.uid",
				resp: M{
					"status":  false,
					"message": "fetched app",
					"data":    map[string]interface{}{},
				},
			},
			wantErrMsg: "field data.uid: not found",
			wantErr:    true,
		},
		{
			name: "should_error_for_field_not_exists_3",
			args: args{
				field: "data.ref.uid",
				resp: M{
					"status":  false,
					"message": "fetched app",
					"data":    map[string]interface{}{},
				},
			},
			wantErrMsg: "field data.ref: not found",
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getKeyInMap(tt.args.field, tt.args.resp)
			if tt.wantErr {
				require.Error(t, err)
				require.Equal(t, tt.wantErrMsg, err.Error())
				return
			}

			require.NoError(t, err)
			require.Equal(t, got, tt.want)
		})
	}
}

func Test_getArrayValue(t *testing.T) {
	type args struct {
		v string
		m M
	}
	var tests = []struct {
		name       string
		args       args
		want       interface{}
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "should_get_array_value",
			args: args{
				v: "tunnels[0]",
				m: map[string]interface{}{
					"tunnels": []interface{}{
						map[string]interface{}{"url": "https://google.com"},
					},
				},
			},
			want:    map[string]interface{}{"url": "https://google.com"},
			wantErr: false,
		},
		{
			name: "should_error_for_invalid_index_string",
			args: args{
				v: "tunnels[abc]",
				m: map[string]interface{}{
					"tunnels": []interface{}{
						map[string]interface{}{"url": "https://google.com"},
					},
				},
			},
			wantErr:    true,
			wantErrMsg: "invalid index notation: abc",
		},
		{
			name: "should_error_for_invalid_index_range",
			args: args{
				v: "tunnels[-2]",
				m: map[string]interface{}{
					"tunnels": []interface{}{
						map[string]interface{}{"url": "https://google.com"},
					},
				},
			},
			wantErr:    true,
			wantErrMsg: "invalid index range: -2",
		},
		{
			name: "should_error_for_field_not_found",
			args: args{
				v: "data[2]",
				m: map[string]interface{}{
					"tunnels": []interface{}{
						map[string]interface{}{"url": "https://google.com"},
					},
				},
			},
			wantErr:    true,
			wantErrMsg: "not found",
		},
		{
			name: "should_error_for_invalid_type",
			args: args{
				v: "tunnels[2]",
				m: map[string]interface{}{
					"tunnels": map[string]interface{}{"url": "https://google.com"},
				},
			},
			wantErr:    true,
			wantErrMsg: "required type is an array but has type map[string]interface {}",
		},
		{
			name: "should_error_for_index_out_of_range",
			args: args{
				v: "tunnels[1]",
				m: map[string]interface{}{
					"tunnels": []interface{}{
						map[string]interface{}{"url": "https://google.com"},
					},
				},
			},
			wantErr:    true,
			wantErrMsg: "index out of range with length 1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getArrayValue(tt.args.v, tt.args.m)
			if tt.wantErr {
				require.Error(t, err)
				require.Equal(t, tt.wantErrMsg, err.Error())
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
