package util

import "testing"

func TestGenerateAndParseToken(t *testing.T) {
	secret := "test-secret-for-unit"
	token, err := GenerateToken(secret, 1, 42, "alice", "admin")
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}
	claims, err := ParseToken(secret, token)
	if err != nil {
		t.Fatalf("parse token: %v", err)
	}
	if claims.UserID != 42 || claims.Username != "alice" || claims.Role != "admin" {
		t.Fatalf("unexpected claims: %+v", claims)
	}
}

func TestParseTokenWrongSecret(t *testing.T) {
	token, err := GenerateToken("secret-a", 1, 1, "bob", "user")
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}
	if _, err := ParseToken("secret-b", token); err == nil {
		t.Fatal("expected error parsing token with wrong secret")
	}
}
