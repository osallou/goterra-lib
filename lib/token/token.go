package goterratoken

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"

	"github.com/fernet/fernet-go"

	terraConfig "github.com/osallou/goterra-lib/lib/config"
)

// TokenTTL is max TTL for a token, can be overiden with env variable GOT_TOKEN_TTL
const TokenTTL time.Duration = 24 * time.Hour

// FernetEncode encode creates a Fernet token from input message
func FernetEncode(msg []byte) (token []byte, err error) {
	config := terraConfig.LoadConfig()
	token = []byte("")

	if config.Fernet == nil || len(config.Fernet) == 0 {
		return nil, fmt.Errorf("no fernet secret defined")
	}

	hash := hmac.New(sha256.New, []byte(config.Fernet[0]))
	secret := hex.EncodeToString(hash.Sum(nil))

	k, kerr := fernet.DecodeKey(secret)
	if kerr != nil {
		fmt.Printf("Failed to decode fernet key: %s\n", kerr)
		return token, kerr
	}
	token, err = fernet.EncryptAndSign(msg, k)
	if err != nil {
		fmt.Printf("Failed to encrypt token: %s\n", err)
		return token, err
	}

	return token, nil
}

// FernetDecode tries to decode a fernet decode with available keys
func FernetDecode(token []byte) (msg []byte, err error) {
	decoded := false
	ttl := TokenTTL
	if os.Getenv("GOT_TOKEN_TTL") != "" {
		osTTL, errTTL := strconv.ParseInt(os.Getenv("GOT_TOKEN_TTL"), 10, 64)
		if errTTL == nil {
			ttl = time.Duration(osTTL) * time.Hour
		}
	}
	config := terraConfig.LoadConfig()
	for _, secret := range config.Fernet {
		hash := hmac.New(sha256.New, []byte(secret))
		secretHash := hex.EncodeToString(hash.Sum(nil))
		k, kerr := fernet.DecodeKeys(secretHash)
		if kerr != nil {
			continue
		}

		msg = fernet.VerifyAndDecrypt(token, ttl, k)
		if msg != nil {
			decoded = true
			break
		}

	}
	if !decoded {
		return msg, errors.New("Failed to decode token")
	}

	return msg, nil

}
