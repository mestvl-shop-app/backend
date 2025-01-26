package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mestvl-shop-app/backend/internal/client/auth"
	"github.com/mestvl-shop-app/backend/internal/config"
	"github.com/mestvl-shop-app/backend/internal/domain"
	"github.com/mestvl-shop-app/backend/internal/repository"
)

type clientService struct {
	clientRepository repository.ClientRepositoryInterface
	authClient       auth.ClientInterface
	cfg              *config.Config
}

func newClientService(
	clientRepository repository.ClientRepositoryInterface,
	authClient auth.ClientInterface,
	cfg *config.Config,
) *clientService {
	return &clientService{
		clientRepository: clientRepository,
		authClient:       authClient,
		cfg:              cfg,
	}
}

type RegisterClientInput struct {
	Email     string
	Password  string
	Firstname string
	Surname   string
	Birthday  *time.Time
	Gender    *domain.ClientGenderString
}

func (s *clientService) Register(ctx context.Context, input *RegisterClientInput) error {
	clientID, err := s.authClient.Register(ctx, &auth.RegisterInput{
		Email:    input.Email,
		Password: input.Password,
	})

	if err != nil {
		if errors.Is(err, auth.ErrClientAlreadyExists) {
			return ClientAlreadyExists
		}
		return fmt.Errorf("register new client failed: %w", err)
	}

	if err := s.clientRepository.Create(ctx, &domain.Client{
		ID:        *clientID,
		Firstname: input.Firstname,
		Surname:   input.Surname,
		Birthday:  input.Birthday,
		Gender:    input.Gender.CodeFromPointer(),
		Email:     input.Email,
	}); err != nil {
		if errors.Is(err, domain.ErrDuplicateEntry) {
			return ClientAlreadyExists
		}
		return fmt.Errorf("create client failed: %w", err)
	}

	return nil
}

func (s *clientService) Login(ctx context.Context, email string, password string) (string, error) {
	token, err := s.authClient.Login(ctx, &auth.LoginInput{
		Email:    email,
		Password: password,
		AppID:    s.cfg.AppID,
	})
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return "", ClientInvalidCredentials
		}
		return "", fmt.Errorf("login failed: %w", err)
	}

	return token, nil
}
