package services

import (
	"context"
	"fmt"

	"github.com/puremike/online_auction_api/internal/models"
	"github.com/puremike/online_auction_api/internal/store"
)

type CSService struct {
	repo store.CSRepository
}

func NewCSService(repo store.CSRepository) *CSService {
	return &CSService{
		repo: repo,
	}
}

func (c *CSService) ContactSupport(ctx context.Context, req *models.ContactSupport) (*models.SupportRes, error) {

	if req.Message == "" || req.Subject == "" || req.UserID == "" {
		return nil, fmt.Errorf("subject, message, and user_id are required")
	}

	supportreq := &models.ContactSupport{
		UserID:  req.UserID,
		Subject: req.Subject,
		Message: req.Message,
	}

	support, err := c.repo.ContactSupport(ctx, supportreq)
	if err != nil {
		return nil, fmt.Errorf("failed to contact support: %w", err)
	}

	return &models.SupportRes{
		ID:      support.ID,
		UserID:  support.UserID,
		Subject: support.Subject,
		Message: support.Message,
	}, nil
}
