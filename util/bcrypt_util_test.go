package util

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "testpassword"
	rounds := 10

	hash, err := HashPassword(password, rounds)
	if err != nil {
		t.Errorf("HashPassword failed: %v", err)
	}
	if hash == "" {
		t.Error("HashPassword returned empty hash")
	}
	if hash == password {
		t.Error("Hash should not be equal to password")
	}
}

func TestHashPasswordWithDifferentRounds(t *testing.T) {
	password := "testpassword"
	hash1, _ := HashPassword(password, 8)
	hash2, _ := HashPassword(password, 12)

	if hash1 == hash2 {
		t.Error("Hashes with different rounds should be different")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "testpassword"
	hash, _ := HashPassword(password, 10)

	if !CheckPasswordHash(password, hash) {
		t.Error("CheckPasswordHash should return true for correct password")
	}

	if CheckPasswordHash("wrongpassword", hash) {
		t.Error("CheckPasswordHash should return false for incorrect password")
	}

	if CheckPasswordHash(password, "wronghash") {
		t.Error("CheckPasswordHash should return false for incorrect hash")
	}
}

func TestCheckPasswordHashEdgeCases(t *testing.T) {
	// Empty password
	hash, _ := HashPassword("", 10)
	if !CheckPasswordHash("", hash) {
		t.Error("CheckPasswordHash should work with empty password")
	}

	// Empty hash
	if CheckPasswordHash("password", "") {
		t.Error("CheckPasswordHash should return false for empty hash")
	}
}
