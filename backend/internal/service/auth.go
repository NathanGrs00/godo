package service

import (
	"context"
	"errors"
	"net/mail"
	"strings"

	"github.com/NathanGrs00/godo/internal/auth"
	"github.com/NathanGrs00/godo/internal/models"
	"github.com/NathanGrs00/godo/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Password limits. bcrypt ignores everything after 72 bytes.
const (
	minPasswordLength = 8
	maxPasswordLength = 72
)

// Credentials is what a client sends to register or log in.
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Session is returned after registering or logging in. The client sends the
// token with every request; logging in on a second device gives that device
// its own token for the same user.
type Session struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

// AuthService registers users and logs them in.
type AuthService struct {
	users  repository.UserRepository
	tokens *auth.TokenManager
	now    clock
}

// NewAuthService creates an AuthService.
func NewAuthService(users repository.UserRepository, tokens *auth.TokenManager) *AuthService {
	return &AuthService{users: users, tokens: tokens, now: systemClock}
}

// Register creates an account and logs it in.
func (s *AuthService) Register(ctx context.Context, in Credentials) (*Session, error) {
	email, err := normalizeEmail(in.Email)
	if err != nil {
		return nil, err
	}
	if len(in.Password) < minPasswordLength {
		return nil, invalid("password", "must be at least 8 characters")
	}
	if len(in.Password) > maxPasswordLength {
		return nil, invalid("password", "must be at most 72 bytes")
	}

	hash, err := auth.HashPassword(in.Password)
	if err != nil {
		return nil, err
	}
	user := &models.User{ID: primitive.NewObjectID(), Email: email, PasswordHash: hash, CreatedAt: s.now()}
	if err := s.users.Create(ctx, user); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return nil, ErrEmailTaken
		}
		return nil, err
	}
	return s.session(user)
}

// Login checks the credentials and returns a new session.
func (s *AuthService) Login(ctx context.Context, in Credentials) (*Session, error) {
	email := strings.ToLower(strings.TrimSpace(in.Email))
	user, err := s.users.GetByEmail(ctx, email)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}
	if !auth.CheckPassword(user.PasswordHash, in.Password) {
		return nil, ErrInvalidCredentials
	}
	return s.session(user)
}

// Me returns the logged-in user.
func (s *AuthService) Me(ctx context.Context, userID primitive.ObjectID) (*models.User, error) {
	return s.users.GetByID(ctx, userID)
}

// session issues a token for user.
func (s *AuthService) session(user *models.User) (*Session, error) {
	token, err := s.tokens.Issue(user.ID)
	if err != nil {
		return nil, err
	}
	return &Session{Token: token, User: user}, nil
}

// normalizeEmail trims and lower-cases an email and checks its format.
func normalizeEmail(email string) (string, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	addr, err := mail.ParseAddress(email)
	if err != nil || addr.Address != email {
		return "", invalid("email", "must be a valid email address")
	}
	return email, nil
}
