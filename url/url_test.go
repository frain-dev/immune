package url

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/frain-dev/immune"
)

func TestParse(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name       string
		args       args
		want       *URL
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "should_parse_url",
			args: args{
				s: "https://localhost:5005/applications/{app_id}/endpoints/{endpoint_id}/{app_id}",
			},
			want: &URL{
				variables: []string{
					"app_id",
					"endpoint_id",
				},
				url: "https://localhost:5005/applications/{app_id}/endpoints/{endpoint_id}/{app_id}",
			},
			wantErr: false,
		},
		{
			name: "should_parse_url",
			args: args{
				s: "https://localhost:5005/applications/{}/endpoints/{endpoint_id}",
			},
			want: &URL{
				variables: []string{
					"endpoint_id",
				},
				url: "https://localhost:5005/applications/{}/endpoints/{endpoint_id}",
			},
			wantErr: false,
		},
		{
			name: "should_error_for_empty_url",
			args: args{
				s: "",
			},
			wantErr:    true,
			wantErrMsg: "url is empty",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.s)
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

func TestURL_ProcessWithVariableMap(t *testing.T) {
	type fields struct {
		variables []string
		url       string
	}
	type args struct {
		vm *immune.VariableMap
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		want       string
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "should_process_url_with_variable_map",
			fields: fields{
				variables: []string{
					"app_id",
					"endpoint_id",
				},
				url: "https://localhost:5005/applications/{app_id}/endpoints/{endpoint_id}",
			},
			args: args{
				vm: &immune.VariableMap{
					VariableToValue: immune.M{
						"app_id":      "123243-4233",
						"endpoint_id": "123-4343-23422",
					},
				},
			},
			want:    "https://localhost:5005/applications/123243-4233/endpoints/123-4343-23422",
			wantErr: false,
		},
		{
			name: "should_process_url_that_has_empty_segment_with_variable_map",
			fields: fields{
				variables: []string{
					"endpoint_id",
				},
				url: "https://localhost:5005/applications/{}/endpoints/{endpoint_id}",
			},
			args: args{
				vm: &immune.VariableMap{
					VariableToValue: immune.M{
						"endpoint_id": "123-4343-23422",
					},
				},
			},
			want:    "https://localhost:5005/applications/{}/endpoints/123-4343-23422",
			wantErr: false,
		},
		{
			name: "should_skip_processing_because_of_no_variables",
			fields: fields{
				variables: nil,
				url:       "https://localhost:5005/applications",
			},
			args: args{
				vm: nil,
			},
			want:    "https://localhost:5005/applications",
			wantErr: false,
		},
		{
			name: "should_error_for_variable_not_found",
			fields: fields{
				variables: []string{"app_id"},
				url:       "https://localhost:5005/applications/{app_id}",
			},
			args: args{
				vm: &immune.VariableMap{VariableToValue: immune.M{}},
			},
			wantErr:    true,
			wantErrMsg: "variable app_id not found in variable map",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &URL{
				variables: tt.fields.variables,
				url:       tt.fields.url,
			}

			got, err := u.ProcessWithVariableMap(tt.args.vm)
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
