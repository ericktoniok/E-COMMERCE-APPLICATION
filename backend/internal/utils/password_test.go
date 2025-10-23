package utils

import "testing"

func TestHashAndCheckPassword(t *testing.T) {
	pw := "S3cretPwd!"
	h, err := HashPassword(pw)
	if err != nil {
		t.Fatalf("hash error: %v", err)
	}
	if !CheckPassword(h, pw) {
		t.Fatalf("expected password to match")
	}
	if CheckPassword(h, "wrong") {
		t.Fatalf("expected password mismatch to be false")
	}
}
