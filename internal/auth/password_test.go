package auth

import "testing"

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("password")
	if err != nil {
		t.Fatalf("error hashing password: %v", err)
	}

	if hash == "" {
		t.Fatal("expected hash should not be empty")
	}

	if hash == "password" {
		t.Fatal("expected hash should not be password")
	}
}

func TestComparePasswords(t *testing.T) {
	hash, err := HashPassword("password")
	if err != nil {
		t.Fatalf("error hashing password: %v", err)
	}

	if !ComparePasswords(hash, []byte("password")) {
		t.Fatal("expected password to match hash")
	}

	if ComparePasswords(hash, []byte("notpassword")) {
		t.Fatalf("expected password to not match hash")
	}
}
