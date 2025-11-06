package auth

import "testing"

const password = "password"

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("error hashing password: %v", err)
	}

	if hash == "" {
		t.Fatal("expected hash should not be empty")
	}

	if hash == password {
		t.Fatal("expected hash should not be password")
	}
}

func TestComparePasswords(t *testing.T) {
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("error hashing password: %v", err)
	}

	if err := ComparePasswords(hash, password); err != nil {
		t.Fatal("expected password to match hash")
	}

	if err := ComparePasswords(hash, "notpassword"); err == nil {
		t.Fatalf("expected password to not match hash and err to be non-null")
	}
}
