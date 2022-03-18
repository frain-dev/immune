package callback

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/frain-dev/immune"
	"github.com/stretchr/testify/require"
)

func Test_handleCallback(t *testing.T) {
	type args struct {
		outbound chan *immune.Signal
	}
	tests := []struct {
		name       string
		args       args
		request    *http.Request
		wantSignal *immune.Signal
	}{
		{
			name: "should_receive_signal",
			args: args{
				outbound: make(chan *immune.Signal, 1),
			},
			request: httptest.NewRequest(http.MethodGet, "/", strings.NewReader(`{"immune_callback_id":"123-4242-13429-4221"}`)),
			wantSignal: &immune.Signal{
				ImmuneCallBackID: "123-4242-13429-4221",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handleFunc := handleCallback(tt.args.outbound)

			recorder := httptest.NewRecorder()
			handleFunc(recorder, tt.request)

			require.Equal(t, http.StatusOK, recorder.Code)
			require.Equal(t, tt.wantSignal, <-tt.args.outbound)
		})
	}
}

func Test_server_ReceiveCallback(t *testing.T) {
	type args struct {
		rc chan *immune.Signal
	}
	tests := []struct {
		name       string
		arrangeFn  func(*server)
		args       args
		wantSignal *immune.Signal
	}{
		{
			name: "should_receive_callback",
			args: args{
				rc: make(chan *immune.Signal, 1),
			},
			wantSignal: &immune.Signal{ImmuneCallBackID: "abc"},
			arrangeFn: func(s *server) {
				go func() {
					s.outbound <- &immune.Signal{ImmuneCallBackID: "abc"}
				}()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &server{outbound: make(chan *immune.Signal)}

			if tt.arrangeFn != nil {
				tt.arrangeFn(s)
			}
			s.ReceiveCallback(tt.args.rc)

			require.Equal(t, tt.wantSignal, <-tt.args.rc)
		})
	}
}
