package service

import (
	"context"
	"errors"
	"quikchat/internal/domain"
	"quikchat/internal/repository"
	"quikchat/pkg/util"
	"time"
)

type userService struct {
	repo               repository.UserRepository
	jwtSecret          string
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
}

func NewUserService(repo repository.UserRepository, jwtSecret string, accessExp, refreshExp time.Duration) *userService {
	return &userService{
		repo:               repo,
		jwtSecret:          jwtSecret,
		accessTokenExpiry:  accessExp,
		refreshTokenExpiry: refreshExp,
	}
}

func (s *userService) Register(ctx context.Context, username, password string) (*domain.User, error) {
	// Check if user already exists
	existingUser, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("username already taken")
	}

	// Hash password
	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &domain.User{
		Username:     username,
		PasswordHash: hashedPassword,
		DisplayName:  username, // Default display name to username
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.Save(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Login(ctx context.Context, username, password string) (accessToken, refreshToken string, err error) {
	user, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		return "", "", err
	}
	if user == nil {
		return "", "", errors.New("invalid credentials")
	}

	if !util.CheckPasswordHash(password, user.PasswordHash) {
		return "", "", errors.New("invalid credentials")
	}

	accessToken, err = util.GenerateToken(user.ID, s.accessTokenExpiry, s.jwtSecret)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = util.GenerateToken(user.ID, s.refreshTokenExpiry, s.jwtSecret)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

