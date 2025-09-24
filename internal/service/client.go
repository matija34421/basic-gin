package service

import (
	"basic-gin/internal/cache"
	"basic-gin/internal/dto"
	"basic-gin/internal/mapper"
	"basic-gin/internal/repository"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

type ClientService struct {
	clientRepository repository.ClientRepository
	cache            cache.Cache
}

func NewClientService(clientRepository repository.ClientRepository, cache cache.Cache) *ClientService {
	return &ClientService{
		clientRepository: clientRepository,
		cache:            cache,
	}
}

func (s *ClientService) GetAll(ctx context.Context) ([]dto.ClientResponse, error) {
	if s.cache != nil {
		if bytes, ok, err := s.cache.Get(ctx, s.keyClientsAll()); err == nil && ok {
			var cached []dto.ClientResponse
			if err := json.Unmarshal(bytes, &cached); err == nil {
				return cached, nil
			}
		}
	}

	clients, err := s.clientRepository.GetAll(ctx)

	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	response := mapper.ClientsToResponseSlice(clients)

	if s.cache != nil {
		if b, err := json.Marshal(response); err == nil {
			_ = s.cache.Set(ctx, s.keyClientsAll(), b, 30*time.Second)
		}
	}

	return response, nil
}

func (s *ClientService) GetById(ctx context.Context, id int64) (*dto.ClientResponse, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid id")
	}

	if s.cache != nil {
		if bytes, ok, err := s.cache.Get(ctx, s.keyClient(id)); err == nil && ok {
			var cacheClient dto.ClientResponse
			if err := json.Unmarshal(bytes, &cacheClient); err == nil {
				return &cacheClient, nil
			}
		}
	}

	client, err := s.clientRepository.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	response := mapper.ClientToResponse(client)

	if s.cache != nil {
		if bytes, err := json.Marshal(response); err == nil {
			_ = s.cache.Set(ctx, s.keyClient(id), bytes, 5*time.Minute)
		}
	}

	return &response, nil
}

func (s *ClientService) Save(ctx context.Context, in dto.ClientCreate) (*dto.ClientResponse, error) {
	validationErr := validateClientCreate(in)

	if validationErr != nil {
		return nil, fmt.Errorf("%w", validationErr)
	}

	client, err := mapper.ToClientFromCreate(in)

	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	saved, createErr := s.clientRepository.CreateClient(ctx, client)

	if createErr != nil {
		return nil, fmt.Errorf("%w", createErr)
	}

	response := mapper.ClientToResponse(saved)

	if s.cache != nil {
		_ = s.cache.Del(ctx, s.keyClientsAll())
		if bytes, err := json.Marshal(response); err == nil {
			_ = s.cache.Set(ctx, s.keyClient(response.ID), bytes, 5*time.Minute)
		}
	}

	return &response, nil
}

func (s *ClientService) Update(ctx context.Context, in dto.ClientUpdate) (*dto.ClientResponse, error) {
	validationErr := validateClientUpdate(in)

	if validationErr != nil {
		return nil, fmt.Errorf("%w", validationErr)
	}

	client, err := mapper.ToClientFromUpdate(in)

	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	saved, createErr := s.clientRepository.UpdateClient(ctx, client)

	if createErr != nil {
		return nil, fmt.Errorf("%w", createErr)
	}

	response := mapper.ClientToResponse(saved)

	if s.cache != nil {
		_ = s.cache.Del(ctx, s.keyClientsAll())
		if bytes, err := json.Marshal(response); err == nil {
			_ = s.cache.Set(ctx, s.keyClient(response.ID), bytes, 5*time.Minute)
		}
	}

	return &response, nil
}

func validateClientCreate(in dto.ClientCreate) error {
	if strings.TrimSpace(in.FirstName) == "" ||
		strings.TrimSpace(in.LastName) == "" ||
		strings.TrimSpace(in.Email) == "" ||
		strings.TrimSpace(in.BirthDate) == "" {
		return errors.New("missing required fields")
	}
	return nil
}

func validateClientUpdate(in dto.ClientUpdate) error {
	if in.ID <= 0 {
		return errors.New("invalid id")
	}
	if strings.TrimSpace(in.FirstName) == "" ||
		strings.TrimSpace(in.LastName) == "" ||
		strings.TrimSpace(in.Email) == "" ||
		strings.TrimSpace(in.BirthDate) == "" {
		return errors.New("missing required fields")
	}
	return nil
}

func (s *ClientService) keyClient(id int64) string { return fmt.Sprintf("client:%d", id) }
func (s *ClientService) keyClientsAll() string     { return "clients:all" }
