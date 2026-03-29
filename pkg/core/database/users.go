package database

import (
	"context"

	"github.com/go-Echelon/go-Echelon/pkg/core/models"
)

type UserStore interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	CountByEmail(ctx context.Context, email string) (int64, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	VerifyUserOTP(ctx context.Context, email string, otp string) (bool, error)
}
