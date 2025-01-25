package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mestvl-shop-app/backend/internal/domain"
	"github.com/mestvl-shop-app/backend/internal/repository"

	"github.com/google/uuid"
)

type clientService struct {
	clientRepository repository.ClientRepositoryInterface
}

func newClientService(
	clientRepository repository.ClientRepositoryInterface,
) *clientService {
	return &clientService{
		clientRepository: clientRepository,
	}
}

type RegisterClientDTO struct {
	Firstname string
	Surname   string
	Birthday  *time.Time
	Gender    *domain.ClientGenderString
}

func (s *clientService) Register(ctx context.Context, dto *RegisterClientDTO) error {
	clientID, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("generate new uuid v7 failed: %w", err)
	}

	if err := s.clientRepository.Create(ctx, &domain.Client{
		ID:        clientID,
		Firstname: dto.Firstname,
		Surname:   dto.Surname,
		Birthday:  dto.Birthday,
		Gender:    dto.Gender.CodeFromPointer(),
	}); err != nil {
		if errors.Is(err, domain.ErrDuplicateEntry) {
			return ClientAlreadyExists
		}
		return fmt.Errorf("create client failed: %w", err)
	}

	return nil
}
