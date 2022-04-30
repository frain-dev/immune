package callback

import (
	"crypto/hmac"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/frain-dev/immune"
	"github.com/stretchr/testify/require"
)

func TestSignatureVerifier_VerifyCallbackSignature(t *testing.T) {
	type fields struct {
		ReplayAttacks bool
		Secret        string
		Header        string
		Hash          string
		hashFn        func() hash.Hash
	}
	type args struct {
		s *immune.Signal
	}
	tests := []struct {
		name                string
		fields              fields
		args                args
		bodyStr             string
		t                   int64
		mismatchRequestBody bool // so we can make cases where the signature should not match, see it's usage
		wantErr             bool
		wantErrMsg          string
	}{
		{
			name:    "should_verify_signature_header_SHA512",
			bodyStr: `{"name":"Daniel"}`,
			t:       time.Now().Unix(),
			fields: fields{
				ReplayAttacks: true,
				Secret:        "1234",
				Header:        "X-Test",
				Hash:          "SHA512",
			},
			args:    args{s: &immune.Signal{}},
			wantErr: false,
		},
		{
			name:    "should_verify_signature_header_MD5",
			bodyStr: `{"name":"Daniel"}`,
			t:       time.Now().Unix(),
			fields: fields{
				ReplayAttacks: true,
				Secret:        "1234",
				Header:        "X-Test",
				Hash:          "MD5",
			},
			args:    args{s: &immune.Signal{}},
			wantErr: false,
		},
		{
			name:    "should_verify_signature_header_SHA1",
			bodyStr: `{"name":"Daniel"}`,
			t:       time.Now().Unix(),
			fields: fields{
				ReplayAttacks: true,
				Secret:        "1234",
				Header:        "X-Test",
				Hash:          "SHA1",
			},
			args:    args{s: &immune.Signal{}},
			wantErr: false,
		},
		{
			name:    "should_verify_signature_header_SHA224",
			bodyStr: `{"name":"Daniel"}`,
			t:       time.Now().Unix(),
			fields: fields{
				ReplayAttacks: true,
				Secret:        "1234",
				Header:        "X-Test",
				Hash:          "SHA224",
			},
			args:    args{s: &immune.Signal{}},
			wantErr: false,
		},
		{
			name:    "should_verify_signature_header_SHA256",
			bodyStr: `{"name":"Daniel"}`,
			t:       time.Now().Unix(),
			fields: fields{
				ReplayAttacks: true,
				Secret:        "1234",
				Header:        "X-Test",
				Hash:          "SHA256",
			},
			args:    args{s: &immune.Signal{}},
			wantErr: false,
		},
		{
			name:    "should_verify_signature_header_SHA384",
			bodyStr: `{"name":"Daniel"}`,
			t:       time.Now().Unix(),
			fields: fields{
				ReplayAttacks: true,
				Secret:        "1234",
				Header:        "X-Test",
				Hash:          "SHA384",
			},
			args:    args{s: &immune.Signal{}},
			wantErr: false,
		},
		{
			name:    "should_verify_signature_header_SHA3_224",
			bodyStr: `{"name":"Daniel"}`,
			t:       time.Now().Unix(),
			fields: fields{
				ReplayAttacks: true,
				Secret:        "1234",
				Header:        "X-Test",
				Hash:          "SHA3_224",
			},
			args:    args{s: &immune.Signal{}},
			wantErr: false,
		},
		{
			name:    "should_verify_signature_header_SHA3_256",
			bodyStr: `{"name":"Daniel"}`,
			t:       time.Now().Unix(),
			fields: fields{
				ReplayAttacks: true,
				Secret:        "1234",
				Header:        "X-Test",
				Hash:          "SHA3_256",
			},
			args:    args{s: &immune.Signal{}},
			wantErr: false,
		},
		{
			name:    "should_verify_signature_header_SHA3_384",
			bodyStr: `{"name":"Daniel"}`,
			t:       time.Now().Unix(),
			fields: fields{
				ReplayAttacks: true,
				Secret:        "1234",
				Header:        "X-Test",
				Hash:          "SHA3_384",
			},
			args:    args{s: &immune.Signal{}},
			wantErr: false,
		},
		{
			name:    "should_verify_signature_header_SHA3_512",
			bodyStr: `{"name":"Daniel"}`,
			t:       time.Now().Unix(),
			fields: fields{
				ReplayAttacks: true,
				Secret:        "1234",
				Header:        "X-Test",
				Hash:          "SHA3_512",
			},
			args:    args{s: &immune.Signal{}},
			wantErr: false,
		},
		{
			name:    "should_verify_signature_header_SHA512_224",
			bodyStr: `{"name":"Daniel"}`,
			t:       time.Now().Unix(),
			fields: fields{
				ReplayAttacks: true,
				Secret:        "1234",
				Header:        "X-Test",
				Hash:          "SHA512_224",
			},
			args:    args{s: &immune.Signal{}},
			wantErr: false,
		},
		{
			name:    "should_verify_signature_header_SHA512_256",
			bodyStr: `{"name":"Daniel"}`,
			t:       time.Now().Unix(),
			fields: fields{
				ReplayAttacks: true,
				Secret:        "1234",
				Header:        "X-Test",
				Hash:          "SHA512_256",
			},
			args:    args{s: &immune.Signal{}},
			wantErr: false,
		},
		{
			name:    "should_error_for_old_timestamp",
			bodyStr: `{"name":"Daniel"}`,
			t:       time.Now().Add(-time.Hour).Unix(),
			fields: fields{
				ReplayAttacks: true,
				Secret:        "1234",
				Header:        "X-Test",
				Hash:          "SHA512_256",
			},
			args:       args{s: &immune.Signal{}},
			wantErr:    true,
			wantErrMsg: "replay attack timestamp is more than a minute ago",
		},
		{
			name:    "should_error_for_invalid_signature",
			bodyStr: `{"name":"Daniel"}`,
			t:       time.Now().Unix(),
			fields: fields{
				ReplayAttacks: true,
				Secret:        "1234",
				Header:        "X-Test",
				Hash:          "SHA512",
			},
			mismatchRequestBody: true,
			args:                args{s: &immune.Signal{}},
			wantErr:             true,
			wantErrMsg:          "signature invalid",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sv, err := NewSignatureVerifier(
				tt.fields.ReplayAttacks,
				tt.fields.Secret,
				tt.fields.Header,
				tt.fields.Hash,
			)
			require.NoError(t, err)

			r, err := http.NewRequest(http.MethodPost, "/", nil)
			require.NoError(t, err)

			tt.args.s.Request = r

			generateSignatureHeader(sv.(*SignatureVerifier), tt.bodyStr, tt.t, tt.mismatchRequestBody, tt.args.s.Request)
			err = sv.VerifyCallbackSignature(tt.args.s)
			if tt.wantErr {
				require.Error(t, err)
				require.Equal(t, tt.wantErrMsg, err.Error())
				return
			}

			require.NoError(t, err)
		})
	}
}

func generateSignatureHeader(sv *SignatureVerifier, bodyStr string, t int64, mismatchRequestBody bool, r *http.Request) {
	body := strings.NewReader(bodyStr)

	var signedPayload strings.Builder
	var timestamp string

	if sv.ReplayAttacks {
		timestamp = fmt.Sprint(t)
		r.Header.Set(ConvoyTimestampHeader, timestamp)
		signedPayload.WriteString(timestamp)
		signedPayload.WriteString(",")
	}
	signedPayload.WriteString(bodyStr)

	if mismatchRequestBody {
		r.Body = io.NopCloser(strings.NewReader(bodyStr + bodyStr)) // scramble the request body
	} else {
		r.Body = io.NopCloser(body) // set the normal value
	}

	h := hmac.New(sv.hashFn, []byte(sv.Secret))
	h.Write([]byte(signedPayload.String()))
	e := hex.EncodeToString(h.Sum(nil))
	r.Header.Set(sv.Header, e)
}

func Test_getHashFunction(t *testing.T) {
	type args struct {
		algorithm string
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:    "should_verify_signature_header_SHA512",
			args:    args{algorithm: "SHA512"},
			wantErr: false,
		},
		{
			name:    "should_verify_signature_header_MD5",
			args:    args{algorithm: "MD5"},
			wantErr: false,
		},
		{
			name:    "should_verify_signature_header_SHA1",
			args:    args{algorithm: "SHA1"},
			wantErr: false,
		},
		{
			name:    "should_verify_signature_header_SHA224",
			args:    args{algorithm: "SHA224"},
			wantErr: false,
		},
		{
			name:    "should_verify_signature_header_SHA256",
			args:    args{algorithm: "SHA256"},
			wantErr: false,
		},
		{
			name:    "should_verify_signature_header_SHA384",
			args:    args{algorithm: "SHA384"},
			wantErr: false,
		},
		{
			name:    "should_verify_signature_header_SHA3_224",
			args:    args{algorithm: "SHA3_224"},
			wantErr: false,
		},
		{
			name:    "should_verify_signature_header_SHA3_256",
			args:    args{algorithm: "SHA3_256"},
			wantErr: false,
		},
		{
			name:    "should_verify_signature_header_SHA3_384",
			args:    args{algorithm: "SHA3_384"},
			wantErr: false,
		},
		{
			name:    "should_verify_signature_header_SHA3_512",
			args:    args{algorithm: "SHA3_512"},
			wantErr: false,
		},
		{
			name:    "should_verify_signature_header_SHA512_224",
			args:    args{algorithm: "SHA512_224"},
			wantErr: false,
		},
		{
			name:    "should_verify_signature_header_SHA512_256",
			args:    args{algorithm: "SHA512_256"},
			wantErr: false,
		},
		{
			name:       "should_error_for_unknown_hash_algorithm",
			args:       args{algorithm: "abc"},
			wantErr:    true,
			wantErrMsg: "unknown hash algorithm",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := getHashFunction(tt.args.algorithm)
			if tt.wantErr {
				require.Error(t, err)
				require.Equal(t, tt.wantErrMsg, err.Error())
				return
			}

			require.NoError(t, err)
		})
	}
}
