package service

import (
	user "github.com/finchatapp/finchat-api/internal/entities/user"
)

type usecase struct {
	repo user.Repository
}

func New(repo user.Repository) user.Usecase {
	return &usecase{
		repo: repo,
	}
}
