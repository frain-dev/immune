package immune

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInjectCallbackID(t *testing.T) {
	type args struct {
		field string
		v     interface{}
		r     M
	}
	tests := []struct {
		name           string
		args           args
		hasSeparator   bool
		wantCallbackID string
		wantErrMsg     string
		wantErr        bool
	}{
		{
			name: "should_inject_one_level_callback_id",
			args: args{
				field: "data",
				v:     "123-564-2242",
				r: M{
					"data": map[string]interface{}{},
				},
			},
			wantCallbackID: "123-564-2242",
			hasSeparator:   false,
			wantErr:        false,
		},
		{
			name: "should_inject_two_level_callback_id",
			args: args{
				field: "data.ref",
				v:     "123-564-2242-2394",
				r: M{
					"data": map[string]interface{}{
						"ref": map[string]interface{}{},
					},
				},
			},
			wantCallbackID: "123-564-2242-2394",
			hasSeparator:   true,
			wantErr:        false,
		},
		{
			name: "should_inject_three_level_callback_id",
			args: args{
				field: "data.ref.caller",
				v:     "123-564-2242-32823",
				r: M{
					"data": map[string]interface{}{
						"ref": map[string]interface{}{
							"caller": map[string]interface{}{},
						},
					},
				},
			},
			wantCallbackID: "123-564-2242-32823",
			hasSeparator:   true,
			wantErr:        false,
		},
		{
			name: "should_error_for_incorrect_field_type",
			args: args{
				field: "data",
				v:     "123-564-2242-32823",
				r: M{
					"data": 123,
				},
			},
			wantErr:    true,
			wantErrMsg: "the field data, is not an object in the request body",
		},
		{
			name: "should_error_for_incorrect_field_type",
			args: args{
				field: "data",
				v:     "123-564-2242-32823",
				r: M{
					"ref": 123,
				},
			},
			wantErr:    true,
			wantErrMsg: "the field data, does not exist",
		},
		{
			name: "should_error_for_field_not_found",
			args: args{
				field: "data.ref.marvel",
				v:     "123-564-2242-32823",
				r: M{
					"data": map[string]interface{}{},
				},
			},
			wantErr:    true,
			wantErrMsg: "field data.ref: not found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := InjectCallbackID(tt.args.field, tt.args.v, tt.args.r)
			if tt.wantErr {
				require.Error(t, err)
				require.Equal(t, tt.wantErrMsg, err.Error())
				return
			}

			require.NoError(t, err)

			if !tt.hasSeparator {
				m := tt.args.r[tt.args.field].(map[string]interface{})
				require.Equal(t, tt.wantCallbackID, m[CallbackIDFieldName])
				return
			}

			parts := strings.Split(tt.args.field, ".")
			m, err := getM(tt.args.r, parts)
			require.NoError(t, err)
			require.Equal(t, tt.wantCallbackID, m[CallbackIDFieldName])
		})
	}
}
