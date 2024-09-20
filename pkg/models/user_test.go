package models

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewUser(t *testing.T) {
	p := &NewUserParams{
		UserName: "test",
		Email:    "user@example.com",
		Password: []byte("password"),
	}

	u, err := NewUser(p)

	if err != nil {
		t.Fatalf("Expected error to be nil, got %s", err)
	}

	if u.ID == uuid.Nil {
		t.Fatal("was expecting uuid to be populated, instead got uuid.nil")
	}

	if string(u.HashedPassword) == string(p.Password) {
		t.Fatal("expected password to be hashed")
	}
}

func TestUserCheckPassword(t *testing.T) {
	p := &NewUserParams{
		UserName: "test",
		Email:    "user@example.com",
		Password: []byte("password"),
	}

	u, err := NewUser(p)

	if err != nil {
		t.Fatalf("Expected error to be nil, got %s", err)
	}

	ok, err := u.CheckPassword([]byte("password"))
	if err != nil {
		t.Fatalf("expected error to be nil, got %s", err)
	}

	if !ok {
		t.Fatal("expected password check to return true, reutrned false")
	}

	ok, err = u.CheckPassword([]byte("wrongpassword"))
	if err != nil {
		t.Fatalf("expected error to be nil, got %s", err)
	}

	if ok {
		t.Fatal("expected password check to return false, returned true")
	}
}
