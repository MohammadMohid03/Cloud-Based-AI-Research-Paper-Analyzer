// Package auth provides authentication services including user registration,
// login, password hashing, and JWT token generation/validation.
package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/research-paper-analyzer/backend/internal/config"
	"github.com/research-paper-analyzer/backend/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Service handles authentication operations.
type Service struct {
	db  *gorm.DB
	cfg *config.Config
}

// NewService creates a new auth service instance.
func NewService(db *gorm.DB, cfg *config.Config) *Service {
	return &Service{
		db:  db,
		cfg: cfg,
	}
}

// Register creates a new user account after validating the input.
// Returns the created user and a JWT token, or an error if registration fails.
func (s *Service) Register(req models.RegisterRequest) (*models.AuthResponse, error) {
	// Check if email is already registered
	var existingUser models.User
	result := s.db.Where("email = ?", req.Email).First(&existingUser)
	if result.Error == nil {
		return nil, errors.New("email already registered")
	}
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("database error: %w", result.Error)
	}

	// Hash the password using bcrypt with cost factor 10
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create the user record
	user := models.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate JWT token for the new user
	token, err := s.generateToken(&user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &models.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

// Login authenticates a user with email and password.
// Returns the user and a JWT token, or an error if authentication fails.
func (s *Service) Login(req models.LoginRequest) (*models.AuthResponse, error) {
	// Find user by email
	var user models.User
	result := s.db.Where("email = ?", req.Email).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("invalid email or password")
	}
	if result.Error != nil {
		return nil, fmt.Errorf("database error: %w", result.Error)
	}

	// Compare the provided password with the stored hash
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Generate JWT token for the authenticated user
	token, err := s.generateToken(&user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &models.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

// generateToken creates a signed JWT token for the given user.
// The token includes the user's ID, email, and expiration time.
func (s *Service) generateToken(user *models.User) (string, error) {
	// Set token expiration based on configuration
	expirationTime := time.Now().Add(time.Duration(s.cfg.JWTExpiryHours) * time.Hour)

	// Create JWT claims
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"name":    user.Name,
		"exp":     expirationTime.Unix(),
		"iat":     time.Now().Unix(),
		"iss":     "research-paper-analyzer",
	}

	// Create and sign the token with HMAC-SHA256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
