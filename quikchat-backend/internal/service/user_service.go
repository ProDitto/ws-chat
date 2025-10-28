package service

import (
	"context"
	"errors"
	"quikchat/internal/domain"
	"quikchat/internal/repository"
	"quikchat/internal/usecase"
	"quikchat/pkg/util"
	"time"
)

type userService struct {
	userRepo           repository.UserRepository
	blockRepo          repository.BlockRepository
	jwtSecret          string
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
}

func NewUserService(userRepo repository.UserRepository, blockRepo repository.BlockRepository, jwtSecret string, accessExp, refreshExp time.Duration) usecase.UserUseCase {
	return &userService{
		userRepo:           userRepo,
		blockRepo:          blockRepo,
		jwtSecret:          jwtSecret,
		accessTokenExpiry:  accessExp,
		refreshTokenExpiry: refreshExp,
	}
}

func (s *userService) Register(ctx context.Context, username, password string) (*domain.User, error) {
	existingUser, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Username:     username,
		PasswordHash: hashedPassword,
		DisplayName:  username,
	}

	if err := s.userRepo.Save(ctx, user); err != nil {
		return nil, err
	}

	user.PasswordHash = ""
	return user, nil
}

func (s *userService) Login(ctx context.Context, username, password string) (accessToken, refreshToken string, err error) {
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return "", "", err
	}
	if user == nil || !util.CheckPasswordHash(password, user.PasswordHash) {
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

func (s *userService) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	user, err := s.userRepo.FindByUsername(ctx, username)
	if user != nil {
		user.PasswordHash = ""
	}
	return user, err
}

func (s *userService) UpdateProfile(ctx context.Context, userID int64, displayName *string, profileImageURL *string) (*domain.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	if displayName != nil {
		user.DisplayName = *displayName
	}
	if profileImageURL != nil {
		user.ProfileImageURL = profileImageURL
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}
	user.PasswordHash = ""
	return user, nil
}

func (s *userService) BlockUser(ctx context.Context, blockerID int64, blockedUsername string) error {
	blockedUser, err := s.userRepo.FindByUsername(ctx, blockedUsername)
	if err != nil {
		return err
	}
	if blockedUser == nil {
		return errors.New("user to block not found")
	}
	if blockerID == blockedUser.ID {
		return errors.New("cannot block yourself")
	}
	return s.blockRepo.Save(ctx, blockerID, blockedUser.ID)
}

func (s *userService) UnblockUser(ctx context.Context, blockerID int64, blockedUsername string) error {
	blockedUser, err := s.userRepo.FindByUsername(ctx, blockedUsername)
	if err != nil {
		return err
	}
	if blockedUser == nil {
		return errors.New("user to unblock not found")
	}
	return s.blockRepo.Delete(ctx, blockerID, blockedUser.ID)
}

