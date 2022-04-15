package immune

import (
	"context"
	"fmt"
	"strings"
)

type CallbackConfiguration struct {
	MaxWaitSeconds uint                   `json:"max_wait_seconds"`
	Port           uint                   `json:"port"`
	Route          string                 `json:"route"`
	SSL            bool                   `json:"ssl" envconfig:"IMMUNE_SSL"`
	SSLKeyFile     string                 `json:"ssl_key_file" envconfig:"IMMUNE_SSL_KEY_FILE"`
	SSLCertFile    string                 `json:"ssl_cert_file" envconfig:"IMMUNE_SSL_CERT_FILE"`
	IDLocation     string                 `json:"id_location"`
	Signature      SignatureConfiguration `json:"signature"`
}

type SignatureConfiguration struct {
	ReplayAttacks bool   `json:"replay_attacks" envconfig:"IMMUNE_REPLAY_ATTACKS"`
	Secret        string `json:"secret" envconfig:"IMMUNE_SIGNATURE_SECRET"`
	Header        string `json:"header" envconfig:"IMMUNE_SIGNATURE_HEADER"`
	Hash          string `json:"hash" envconfig:"IMMUNE_SIGNATURE_HASH"`
}

const CallbackIDFieldName = "immune_callback_id"

// InjectCallbackID injects a callback id into field(expected to be a map[string]interface{} in r)
// in r, using CallbackIDFieldName as the key and value as the value.
func InjectCallbackID(field string, value interface{}, r M) error {
	// we may have separators referencing deeper fields in r e.g data.uid
	parts := strings.Split(field, ".")
	if len(parts) < 2 { // if it's less than 2, then there is no '.' in field
		v, ok := r[field]
		if !ok {
			return fmt.Errorf("the field %s, does not exist", field)
		}

		m, ok := v.(map[string]interface{})
		if !ok {
			return fmt.Errorf("the field %s, is not an object in the request body", field)
		}
		m[CallbackIDFieldName] = value
		return nil
	}

	nextLevel, err := getM(r, parts)
	if err != nil {
		return err
	}

	nextLevel[CallbackIDFieldName] = value

	return nil
}

type CallbackServer interface {
	ReceiveCallback(rc chan<- *Signal)
	Start(ctx context.Context) error
	Stop()
}
