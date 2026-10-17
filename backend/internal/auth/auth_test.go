package auth

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const secret = "test-secret-test-secret-test-secret"

func TestPassword(t *testing.T) {
	hash, err := HashPassword("hunter22")
	if err != nil {
		t.Fatal(err)
	}
	if !CheckPassword(hash, "hunter22") {
		t.Error("correct password rejected")
	}
	if CheckPassword(hash, "hunter23") {
		t.Error("wrong password accepted")
	}
}

func TestTokenRoundTrip(t *testing.T) {
	m := NewTokenManager(secret, time.Hour)
	user := primitive.NewObjectID()
	token, err := m.Issue(user)
	if err != nil {
		t.Fatal(err)
	}
	got, err := m.Verify(token)
	if err != nil || got != user {
		t.Fatalf("Verify = %v, %v; want %v", got, err, user)
	}
}

func TestTokenRejected(t *testing.T) {
	m := NewTokenManager(secret, time.Hour)
	token, _ := m.Issue(primitive.NewObjectID())

	t.Run("expired", func(t *testing.T) {
		later := NewTokenManager(secret, time.Hour)
		later.now = func() time.Time { return time.Now().Add(2 * time.Hour) }
		if _, err := later.Verify(token); err != ErrInvalidToken {
			t.Fatalf("err = %v", err)
		}
	})
	t.Run("wrong secret", func(t *testing.T) {
		other := NewTokenManager("another-secret-another-secret-another", time.Hour)
		if _, err := other.Verify(token); err != ErrInvalidToken {
			t.Fatalf("err = %v", err)
		}
	})
	t.Run("garbage", func(t *testing.T) {
		if _, err := m.Verify("not.a.token"); err != ErrInvalidToken {
			t.Fatalf("err = %v", err)
		}
	})
	t.Run("unsigned", func(t *testing.T) {
		// Header {"alg":"none"} with a subject and no signature
		none := "eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0.eyJzdWIiOiI2NTAwMDAwMDAwMDAwMDAwMDAwMDAwMDAiLCJleHAiOjk5OTk5OTk5OTl9."
		if _, err := m.Verify(none); err != ErrInvalidToken {
			t.Fatalf("err = %v", err)
		}
	})
}
