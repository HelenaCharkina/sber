package service

import (
	"context"
	"sber/pkg/repository"
	"sber/types"
)


type User interface {
	GetById(ctx context.Context, userId string) (*types.User, error)
}

type Service struct {
	User
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		User: NewUserService(repo.User),
	}
}