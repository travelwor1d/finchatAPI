package service

import (
	contact "github.com/finchatapp/finchat-api/internal/entities/contact"
)

type usecase struct {
	repo contact.Repository
}

func New(repo contact.Repository) contact.Usecase {
	return &usecase{
		repo: repo,
	}
}
