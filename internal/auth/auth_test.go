package auth_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/praneeth-ayla/Chirpy/internal/auth"
)

func TestMakeAndValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "mysecret"

	token, err := auth.MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	gotID, err := auth.ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("ValidateJWT failed: %v", err)
	}

	if gotID != userID {
		t.Errorf("expected %v, got %v", userID, gotID)
	}
}

func TestExpiredJWT(t *testing.T) {
	userID := uuid.New()
	secret := "mysecret"

	token, err := auth.MakeJWT(userID, secret, -time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	_, err = auth.ValidateJWT(token, secret)
	if err == nil {
		t.Fatal("expected error for expired JWT but got none")
	}
}

func TestWrongSecret(t *testing.T) {
	userID := uuid.New()

	token, err := auth.MakeJWT(userID, "correctsecret", time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	_, err = auth.ValidateJWT(token, "wrongsecret")
	if err == nil {
		t.Fatal("expected error for wrong secret but got none")
	}
}
