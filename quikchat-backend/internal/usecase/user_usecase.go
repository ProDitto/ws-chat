package usecase

import (
	"context"
	"quikchat/internal/domain"
)

type UserUseCase interface {
	Register(ctx context.Context, username, password string) (*domain.User, error)
	Login(ctx context.Context, username, password string) (accessToken, refreshToken string, err error)
}

