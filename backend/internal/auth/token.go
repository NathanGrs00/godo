package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ErrInvalidToken means a token is malformed, expired or wrongly signed.
var ErrInvalidToken = errors.New("invalid token")

// TokenManager issues and verifies signed login tokens (JWTs).
// A token identifies a user, so any device holding it sees that user's data.
type TokenManager struct {
	secret []byte
	ttl    time.Duration
	now    func() time.Time
}

// NewTokenManager creates a TokenManager that signs with secret and issues
// tokens valid for ttl.
func NewTokenManager(secret string, ttl time.Duration) *TokenManager {
	return &TokenManager{secret: []byte(secret), ttl: ttl, now: time.Now}
}

// Issue creates a token for the user.
func (m *TokenManager) Issue(userID primitive.ObjectID) (string, error) {
	now := m.now()
	claims := jwt.RegisteredClaims{
		Subject:   userID.Hex(),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(m.ttl)),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(m.secret)
}

// Verify checks a token and returns the user it was issued for.
func (m *TokenManager) Verify(token string) (primitive.ObjectID, error) {
	var claims jwt.RegisteredClaims
	_, err := jwt.ParseWithClaims(token, &claims, func(*jwt.Token) (any, error) {
		return m.secret, nil
	},
		// Only accept the algorithm we sign with
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithTimeFunc(m.now),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return primitive.NilObjectID, ErrInvalidToken
	}

	userID, err := primitive.ObjectIDFromHex(claims.Subject)
	if err != nil {
		return primitive.NilObjectID, ErrInvalidToken
	}
	return userID, nil
}
