package utils_test

import (
	"crypto/rand"
	"encoding/base64"
	"testing"

	"github.com/guneyin/printhub/utils"
)

var plainText = "plain-text"

func genSecret() string {
	key := make([]byte, 16)
	_, _ = rand.Read(key)
	return base64.StdEncoding.EncodeToString(key)
}

func TestAES(t *testing.T) {
	secret := genSecret()
	enc, err := utils.Encrypt(plainText, []byte(secret))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("enc:", enc)
	dec, err := utils.Decrypt(enc, []byte(secret))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("dec:", dec)
	if dec != plainText {
		t.Fatalf("Decrypt failed, expected %s, got %s", plainText, dec)
	}
}
