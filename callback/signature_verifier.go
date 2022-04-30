package callback

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/frain-dev/immune"

	"golang.org/x/crypto/sha3"

	"github.com/pkg/errors"
)

type SignatureVerifier struct {
	ReplayAttacks bool   `json:"replay_attacks"`
	Secret        string `json:"secret"`
	Header        string `json:"header"`
	Hash          string `json:"hash"`
	hashFn        func() hash.Hash
}

func NewSignatureVerifier(replayAttacks bool, secret, header, hash string) (*SignatureVerifier, error) {
	fn, err := getHashFunction(hash)
	if err != nil {
		return nil, err
	}

	return &SignatureVerifier{
		ReplayAttacks: replayAttacks,
		Secret:        secret,
		Header:        header,
		Hash:          hash,
		hashFn:        fn,
	}, nil
}

func (sv *SignatureVerifier) VerifyCallbackSignature(s *immune.Signal) error {
	r := s.Request
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errors.Wrap(err, "unable to read request body")
	}

	signatureHex := []byte(r.Header.Get(sv.Header))
	var signature = make([]byte, hex.DecodedLen(len(signatureHex)))
	_, err = hex.Decode(signature, signatureHex)
	if err != nil {
		return errors.Wrap(err, "unable to hex decode signature body")
	}

	hasher := hmac.New(sv.hashFn, []byte(sv.Secret))

	if sv.ReplayAttacks {
		timestampStr := r.Header.Get("Convoy-Timestamp")
		timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
		if err != nil {
			return errors.Wrap(err, "unable to parse signature timestamp")
		}

		t := time.Unix(timestamp, 0)
		d := time.Since(t)
		if d > time.Minute {
			return errors.Errorf("replay attack timestamp is more than a minute ago")
		}

		hasher.Write([]byte(timestampStr))
		hasher.Write([]byte(","))
	}

	hasher.Write(buf)
	if !hmac.Equal(signature, hasher.Sum(nil)) {
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
