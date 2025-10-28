package usecase

import (
	"context"
	"quikchat/internal/domain"
)

type UserUseCase interface {
	Register(ctx context.Context, username, password string) (*domain.User, error)
	Login(ctx context.Context, username, password string) (accessToken, refreshToken string, err error)
	GetUserByUsername(ctx context.Context, username string) (*domain.User, error)
	UpdateProfile(ctx context.Context, userID int64, displayName *string, profileImageURL *string) (*domain.User, error)
	BlockUser(ctx context.Context, blockerID int64, blockedUsername string) error
	UnblockUser(ctx context.Context, blockerID int64, blockedUsername string) error
}

