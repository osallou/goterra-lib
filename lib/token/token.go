package goterratoken

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/fernet/fernet-go"

	terraConfig "github.com/osallou/goterra-lib/lib/config"
)

// TokenTTL is max TTL for a token, can be overiden with env variable GOT_TOKEN_TTL
const TokenTTL time.Duration = 24 * time.Hour

// FernetEncode encode creates a Fernet token from input message
func FernetEncode(msg []byte) (token []byte, err error) {
	config := terraConfig.LoadConfig()
	token = []byte("")
	k, kerr := fernet.DecodeKey(config.Fernet[0])
	if kerr != nil {
		return token, err
	}
	token, err = fernet.EncryptAndSign(msg, k)
	if err != nil {
		return token, err
	}
	//msg := fernet.VerifyAndDecrypt(tok, 60*time.Second, k)
	//fmt.Println(string(msg))
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
		k, kerr := fernet.DecodeKeys(secret)
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
