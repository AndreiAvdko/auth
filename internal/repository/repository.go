package repository

import (
	"context"

	"github.com/AndreiAvdko/auth/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, info *model.User) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
}
