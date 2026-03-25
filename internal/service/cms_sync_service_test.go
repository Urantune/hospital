package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"testing"
)

func TestValidateCMSSignatureAcceptsHexAndPrefixedHex(t *testing.T) {
	body := []byte(`{"event_id":"evt-1"}`)
	secret := "cms-secret"
	t.Setenv("CMS_SYNC_SECRET", secret)

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	signature := hex.EncodeToString(mac.Sum(nil))

	if !ValidateCMSSignature(body, signature) {
		t.Fatal("expected raw hex signature to be accepted")
	}

	if !ValidateCMSSignature(body, "sha256="+signature) {
		t.Fatal("expected sha256-prefixed signature to be accepted")
	}
}

func TestValidateCMSSignatureRejectsWhenSecretMissing(t *testing.T) {
	body := []byte(`{"event_id":"evt-1"}`)
	_ = os.Unsetenv("CMS_SYNC_SECRET")

	if ValidateCMSSignature(body, "anything") {
		t.Fatal("expected signature validation to fail without secret")
	}
}
