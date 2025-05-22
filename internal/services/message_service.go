package services

import (
	"be/internal/models"
	"be/internal/repositories"
	"context"
	"errors"
)

// MessageService provides business logic methods for Message
type MessageService struct {
	messageRepo *repositories.MessageRepository
}

// NewMessageService creates a new MessageService
func NewMessageService(messageRepo *repositories.MessageRepository) *MessageService {
	return &MessageService{messageRepo}
}

// Create creates a new message
func (s *MessageService) Create(ctx context.Context, message *models.Message) (*models.Message, error) {
	if message.GroupID.IsZero() {
		return nil, errors.New("group ID is required")
	}
	if message.SenderID.IsZero() {
		return nil, errors.New("sender ID is required")
	}
	if message.Content == "" {
		return nil, errors.New("message content is required")
	}
	return s.messageRepo.Create(ctx, message)
}

// GetByID retrieves a message by ID
func (s *MessageService) GetByID(ctx context.Context, id string) (*models.Message, error) {
	return s.messageRepo.GetByID(ctx, id)
}

// GetByGroupID retrieves messages by group ID
func (s *MessageService) GetByGroupID(ctx context.Context, groupID string) ([]models.Message, error) {
	return s.messageRepo.GetByGroupID(ctx, groupID)
}

// GetAll retrieves all messages with pagination
func (s *MessageService) GetAll(ctx context.Context, page, limit int64) ([]models.Message, error) {
	return s.messageRepo.GetAll(ctx, page, limit)
}

// Update updates message information
func (s *MessageService) Update(ctx context.Context, message *models.Message) (*models.Message, error) {
	if message.ID.IsZero() {
		return nil, errors.New("invalid message ID")
	}
	if message.GroupID.IsZero() {
		return nil, errors.New("group ID is required")
	}
	if message.SenderID.IsZero() {
		return nil, errors.New("sender ID is required")
	}
	if message.Content == "" {
		return nil, errors.New("message content is required")
	}
	return s.messageRepo.Update(ctx, message)
}

// Delete deletes a message by ID
func (s *MessageService) Delete(ctx context.Context, id string) error {
	return s.messageRepo.Delete(ctx, id)
}