package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"sber/types"
)

type User interface {
	GetById(ctx context.Context, userId string) (*types.User, error)
}

type Repository struct {
	User
}

func NewRepository(client *mongo.Client) *Repository {
	return &Repository{
		User: NewUserMongo(client),
	}
}
