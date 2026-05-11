package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/user/tahap2-rest-api/config"
	"github.com/user/tahap2-rest-api/internal/middleware"
	"github.com/user/tahap2-rest-api/internal/model"
	"github.com/user/tahap2-rest-api/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPhoneAlreadyRegistered = errors.New("phone already registered")
	ErrInvalidCredentials      = errors.New("phone number and PIN doesn't match")
	ErrUserNotFound             = errors.New("user not found")
)

type RegisterRequest struct {
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Address     string `json:"address"`
	PIN         string `json:"pin" binding:"required,min=6,max=6"`
}

type LoginRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	PIN         string `json:"pin" binding:"required"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Register(ctx context.Context, req *RegisterRequest) (*model.User, error) {
	// Check if phone already registered
	existing, err := s.userRepo.FindByPhone(ctx, req.PhoneNumber)
	if err != nil && !errors.Is(err, repository.ErrUserNotFound) {
		return nil, fmt.Errorf("check existing user: %w", err)
	}
	if existing != nil {
		return nil, ErrPhoneAlreadyRegistered
	}

	// Hash PIN with bcrypt
	hashedPin, err := bcrypt.GenerateFromPassword([]byte(req.PIN), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash pin: %w", err)
	}

	user := &model.User{
		UserID:      uuid.New().String(),
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		Address:     req.Address,
		PIN:         string(hashedPin),
		Balance:     0,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*TokenPair, error) {
	user, err := s.userRepo.FindByPhone(ctx, req.PhoneNumber)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("find user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PIN), []byte(req.PIN)); err != nil {
		return nil, ErrInvalidCredentials
	}

	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) generateAccessToken(user *model.User) (string, error) {
	claims := &middleware.Claims{
		UserID:      user.UserID,
		PhoneNumber: user.PhoneNumber,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.UserID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.GetConfig().JWTSecret))
}

func (s *AuthService) generateRefreshToken(user *model.User) (string, error) {
	claims := &middleware.Claims{
		UserID:      user.UserID,
		PhoneNumber: user.PhoneNumber,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.UserID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.GetConfig().JWTSecret))
}
