package auth

import (
	"testing"
	"time"
)

func TestJWTGenerateAndParse(t *testing.T) {
	m := NewJWTManager("test-secret", time.Hour)
	tok, err := m.Generate(42, "admin")
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}
	claims, err := m.Parse(tok)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if claims.UserID != 42 || claims.Role != "admin" {
		t.Fatalf("unexpected claims: %+v", claims)
	}
	if time.Until(claims.ExpiresAt.Time) <= 0 {
		t.Fatalf("token already expired: %v", claims.ExpiresAt.Time)
	}
}

func TestJWTExpired(t *testing.T) {
	m := NewJWTManager("test-secret", -1*time.Second)
	tok, err := m.Generate(1, "customer")
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}
	if _, err := m.Parse(tok); err == nil {
		t.Fatalf("expected parse error for expired token, got nil")
	}
}
