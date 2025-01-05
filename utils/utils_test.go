package utils

import "testing"

var (
	secret    = "N1PCdw3M2B1TfJhoaY2mL736p2vCUc47"
	plainText = "plain-text"
)

func TestAES(t *testing.T) {
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
