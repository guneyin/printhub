package utils

import (
	"crypto/rand"
	"encoding/base64"
	"testing"
)

var plainText = "plain-text"

func genSecret() string {
	key := make([]byte, 16)
	_, _ = rand.Read(key)
	return base64.StdEncoding.EncodeToString(key)
}

func TestAES(t *testing.T) {
	secret := genSecret()
	enc, err := Encrypt(plainText, []byte(secret))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("enc:", enc)
	dec, err := Decrypt(enc, []byte(secret))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("dec:", dec)
	if dec != plainText {
		t.Fatalf("Decrypt failed, expected %s, got %s", plainText, dec)
	}
}
