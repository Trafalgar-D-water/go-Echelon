package services

import (
	"github.com/go-Echelon/go-Echelon/pkg/core/database"
)

type UserService struct {
	db *database.Database
}

func NewUserService(db *database.Database) *UserService {
	return &UserService{db: db}
}
