package service

import (
	"errors"
	"strings"
	"testing"
)

func TestRegisterAndLogin(t *testing.T) {
	f := newFixture()
	session, err := f.auth.Register(ctx, Credentials{Email: "  Alice@Example.com ", Password: "correct horse"})
	if err != nil {
		t.Fatal(err)
	}
	if session.Token == "" || session.User.Email != "alice@example.com" {
		t.Fatalf("session = %+v", session)
	}
	if session.User.PasswordHash == "correct horse" {
		t.Fatal("password stored in plain text")
	}

	// Logging in from a second device returns a session for the same user
	second, err := f.auth.Login(ctx, Credentials{Email: "ALICE@example.com", Password: "correct horse"})
	if err != nil {
		t.Fatal(err)
	}
	if second.User.ID != session.User.ID {
		t.Fatal("login returned a different user")
	}

	me, err := f.auth.Me(ctx, session.User.ID)
	if err != nil || me.Email != "alice@example.com" {
		t.Fatalf("Me = %+v, %v", me, err)
	}
}

func TestRegisterValidation(t *testing.T) {
	f := newFixture()
	_, err := f.auth.Register(ctx, Credentials{Email: "not-an-email", Password: "long enough"})
	assertValidation(t, err, "email")
	_, err = f.auth.Register(ctx, Credentials{Email: "Bob <bob@example.com>", Password: "long enough"})
	assertValidation(t, err, "email")
	_, err = f.auth.Register(ctx, Credentials{Email: "bob@example.com", Password: "short"})
	assertValidation(t, err, "password")
	_, err = f.auth.Register(ctx, Credentials{Email: "bob@example.com", Password: strings.Repeat("a", 73)})
	assertValidation(t, err, "password")

	if _, err := f.auth.Register(ctx, Credentials{Email: "bob@example.com", Password: "long enough"}); err != nil {
		t.Fatal(err)
	}
	if _, err := f.auth.Register(ctx, Credentials{Email: "BOB@example.com", Password: "long enough"}); !errors.Is(err, ErrEmailTaken) {
		t.Fatalf("duplicate err = %v, want ErrEmailTaken", err)
	}
}

func TestLoginRejectsBadCredentials(t *testing.T) {
	f := newFixture()
	if _, err := f.auth.Register(ctx, Credentials{Email: "a@example.com", Password: "right password"}); err != nil {
		t.Fatal(err)
	}
	for _, c := range []Credentials{
		{Email: "a@example.com", Password: "wrong password"},
		{Email: "nobody@example.com", Password: "right password"},
	} {
		if _, err := f.auth.Login(ctx, c); !errors.Is(err, ErrInvalidCredentials) {
			t.Errorf("Login(%s) err = %v, want ErrInvalidCredentials", c.Email, err)
		}
	}
}
