package immune

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/sha3"

	"github.com/pkg/errors"
)

type SignatureVerifier struct {
	ReplayAttacks bool   `json:"replay_attacks"`
	Secret        string `json:"secret"`
	Header        string `json:"header"`
	Hash          string `json:"hash"`
}

func (sv *SignatureVerifier) VerifySignatureHeader(r *http.Request) error {
	signature := r.Header.Get(sv.Header)

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errors.Wrap(err, "unable to read request body")
	}

	fn, err := getHashFunction(sv.Hash)
	if err != nil {
		return err
	}

	hasher := hmac.New(fn, []byte(sv.Secret))
	hasher.Write(buf)
	sum := hasher.Sum(nil)

	var actual = make([]byte, len(sum))
	_, err = hex.Decode(actual, []byte(signature))
	if err != nil {
		return errors.Wrap(err, "unable to hex decode signature body")
	}

	if sv.ReplayAttacks {
		parts := bytes.Split(actual, []byte(","))
		if len(parts) < 2 {
			return errors.Errorf(`replay attack signature header must have 2 parts seperated by ","`)
		}

		timestamp, err := strconv.ParseInt(string(parts[0]), 10, 0)
		if err != nil {
			return errors.Wrap(err, "unable to parse signature timestamp")
		}

		t := time.Unix(timestamp, 0)
		d := time.Since(t)
		if d > time.Minute {
			return errors.Errorf("replay attack timestamp is more than a minute ago")
		}
		actual = parts[1]
	}

	if !hmac.Equal(actual, sum) {
		return errors.New("signature invalid")
	}
	return nil
}

func getHashFunction(algorithm string) (func() hash.Hash, error) {
	switch algorithm {
	case "MD5":
		return md5.New, nil
	case "SHA1":
		return sha1.New, nil
	case "SHA224":
		return sha256.New224, nil
	case "SHA256":
		return sha256.New, nil
	case "SHA384":
		return sha512.New384, nil
	case "SHA512":
		return sha512.New, nil
	case "SHA3_224":
		return sha3.New224, nil
	case "SHA3_256":
		return sha3.New256, nil
	case "SHA3_384":
		return sha3.New384, nil
	case "SHA3_512":
		return sha3.New512, nil
	case "SHA512_224":
		return sha512.New512_224, nil
	case "SHA512_256":
		return sha512.New512_256, nil
	}
	return nil, errors.New("unknown hash algorithm")
}
