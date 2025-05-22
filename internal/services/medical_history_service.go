package services

import (
	"be/internal/models"
	"be/internal/repositories"
	"context"
	"errors"
)

// MedicalHistoryService provides business logic methods for MedicalHistory
type MedicalHistoryService struct {
	medicalHistoryRepo *repositories.MedicalHistoryRepository
}

// NewMedicalHistoryService creates a new MedicalHistoryService
func NewMedicalHistoryService(medicalHistoryRepo *repositories.MedicalHistoryRepository) *MedicalHistoryService {
	return &MedicalHistoryService{medicalHistoryRepo}
}

// Create creates a new medical history record
func (s *MedicalHistoryService) Create(ctx context.Context, medicalHistory *models.MedicalHistory) (*models.MedicalHistory, error) {
	if medicalHistory.HealthID.IsZero() {
		return nil, errors.New("health ID is required")
	}
	return s.medicalHistoryRepo.Create(ctx, medicalHistory)
}

// GetByID retrieves a medical history by ID
func (s *MedicalHistoryService) GetByID(ctx context.Context, id string) (*models.MedicalHistory, error) {
	return s.medicalHistoryRepo.GetByID(ctx, id)
}

// GetByHealthID retrieves a medical history by health ID
func (s *MedicalHistoryService) GetByHealthID(ctx context.Context, healthID string) (*models.MedicalHistory, error) {
	return s.medicalHistoryRepo.GetByHealthID(ctx, healthID)
}

// GetAll retrieves all medical history records with pagination
func (s *MedicalHistoryService) GetAll(ctx context.Context, page, limit int64) ([]models.MedicalHistory, error) {
	return s.medicalHistoryRepo.GetAll(ctx, page, limit)
}

// Update updates medical history information
func (s *MedicalHistoryService) Update(ctx context.Context, medicalHistory *models.MedicalHistory) (*models.MedicalHistory, error) {
	if medicalHistory.ID.IsZero() {
		return nil, errors.New("invalid medical history ID")
	}
	if medicalHistory.HealthID.IsZero() {
		return nil, errors.New("health ID is required")
	}
	return s.medicalHistoryRepo.Update(ctx, medicalHistory)
}

// Delete deletes a medical history by ID
func (s *MedicalHistoryService) Delete(ctx context.Context, id string) error {
	return s.medicalHistoryRepo.Delete(ctx, id)
}