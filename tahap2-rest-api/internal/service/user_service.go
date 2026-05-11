package service

import (
	"context"
	"fmt"

	"github.com/user/tahap2-rest-api/internal/model"
	"github.com/user/tahap2-rest-api/internal/repository"
)

type UpdateProfileRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Address   string `json:"address"`
}

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) UpdateProfile(ctx context.Context, userID string, req *UpdateProfileRequest) (*model.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}

	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Address = req.Address

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("update profile: %w", err)
	}

	updated, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("fetch updated user: %w", err)
	}

	return updated, nil
}