package service

import (
	"context"
	"sber/pkg/repository"
	"sber/types"
)

type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) GetById(ctx context.Context, userId string) (*types.User, error) {
	return s.repo.GetById(ctx, userId)
}