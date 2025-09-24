package service

import (
	"basic-gin/internal/cache"
	"basic-gin/internal/dto"
	"basic-gin/internal/mapper"
	"basic-gin/internal/model"
	"basic-gin/internal/repository"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"time"
)

type AccountService struct {
	accountRepository repository.AccountRepository
	clientService     ClientService
	cache             cache.Cache
}

func NewAccountService(accountRepository *repository.AccountRepository, clientService *ClientService, cache cache.Cache) *AccountService {
	return &AccountService{
		accountRepository: *accountRepository,
		clientService:     *clientService,
		cache:             cache,
	}
}

func (s *AccountService) GetById(ctx context.Context, id int) (*dto.AccountResponse, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid id")
	}

	if s.cache != nil {
		if bytes, ok, err := s.cache.Get(ctx, s.keyAccount(id)); err == nil && ok {
			var account dto.AccountResponse
			if err := json.Unmarshal(bytes, &account); err == nil {
				return &account, nil
			}
		}
	}

	account, err := s.accountRepository.GetById(ctx, id)

	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	response := mapper.AccountToResponse(account)

	if s.cache != nil {
		if bytes, err := json.Marshal(response); err == nil {
			_ = s.cache.Set(ctx, s.keyAccount(response.ID), bytes, 5*time.Minute)
		}
	}

	return response, nil
}

func (s *AccountService) GetByClientId(ctx context.Context, id int) ([]*dto.AccountResponse, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid id")
	}

	if _, err := s.clientService.GetById(ctx, int64(id)); err != nil {
		return nil, fmt.Errorf("client not found: %w", err)
	}

	key := s.keyAccountsByClient(id)

	if s.cache != nil {
		if b, ok, err := s.cache.Get(ctx, key); err == nil && ok {
			var cached []*dto.AccountResponse
			if err := json.Unmarshal(b, &cached); err == nil {
				return cached, nil
			}
		}
	}

	accounts, err := s.accountRepository.GetByClientId(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get accounts by client id: %w", err)
	}

	resp := mapper.AccountsToResponseSlice(accounts)

	if s.cache != nil {
		if b, err := json.Marshal(resp); err == nil {
			_ = s.cache.Set(ctx, key, b, 60*time.Second)
		}
	}

	return resp, nil
}

func (s *AccountService) Save(ctx context.Context, clientId int) (dto.AccountResponse, error) {
	if clientId <= 0 {
		return dto.AccountResponse{}, fmt.Errorf("invalid client id")
	}

	if _, err := s.clientService.GetById(ctx, int64(clientId)); err != nil {
		return dto.AccountResponse{}, fmt.Errorf("client not found: %w", err)
	}

	const maxAttempts = 3
	var saved *model.Account
	var err error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		acc := model.Account{
			ClientId:      clientId,
			AccountNumber: generateAccountNumber(16),
			Balance:       0,
		}

		saved, err = s.accountRepository.CreateAccount(ctx, &acc)
		if err == nil {
			break
		}
		if strings.Contains(strings.ToLower(err.Error()), "account_number already exists") {
			if attempt < maxAttempts {
				continue
			}
		}
		return dto.AccountResponse{}, fmt.Errorf("create account: %w", err)
	}

	if saved == nil {
		return dto.AccountResponse{}, fmt.Errorf("could not create account after retries")
	}

	respPtr := mapper.AccountToResponse(saved)
	resp := *respPtr

	if s.cache != nil {
		_ = s.cache.Del(ctx, s.keyAccountsByClient(clientId))

		if b, mErr := json.Marshal(resp); mErr == nil {
			_ = s.cache.Set(ctx, s.keyAccount(resp.ID), b, 5*time.Minute)
		}
	}

	return resp, nil
}

func (s *AccountService) Deposit(ctx context.Context, id int, amount float64) (*dto.AccountResponse, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid account id")
	}
	if amount <= 0 {
		return nil, fmt.Errorf("amount must be positive")
	}

	tx, err := s.accountRepository.Pool().Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	updated, err := s.accountRepository.UpdateBalanceDelta(ctx, tx, id, amount)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	if s.cache != nil {
		_ = s.cache.Del(ctx, s.keyAccount(updated.ID))
		_ = s.cache.Del(ctx, s.keyAccountsByClient(updated.ClientId))
		if b, mErr := json.Marshal(mapper.AccountToResponse(updated)); mErr == nil {
			_ = s.cache.Set(ctx, s.keyAccount(updated.ID), b, 5*time.Minute)
		}
	}

	return mapper.AccountToResponse(updated), nil
}

func (s *AccountService) Withdraw(ctx context.Context, id int, amount float64) (*dto.AccountResponse, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid account id")
	}
	if amount <= 0 {
		return nil, fmt.Errorf("amount must be positive")
	}

	tx, err := s.accountRepository.Pool().Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	updated, err := s.accountRepository.UpdateBalanceDelta(ctx, tx, id, -amount)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	if s.cache != nil {
		_ = s.cache.Del(ctx, s.keyAccount(updated.ID))
		_ = s.cache.Del(ctx, s.keyAccountsByClient(updated.ClientId))
		if b, mErr := json.Marshal(mapper.AccountToResponse(updated)); mErr == nil {
			_ = s.cache.Set(ctx, s.keyAccount(updated.ID), b, 5*time.Minute)
		}
	}

	return mapper.AccountToResponse(updated), nil
}

func generateAccountNumber(n int) string {
	const digits = "0123456789"
	if n <= 0 {
		n = 16
	}
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		// crypto/rand
		idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		b[i] = digits[idx.Int64()]
	}
	return string(b)
}

func (s *AccountService) keyAccount(id int) string { return fmt.Sprintf("account:%d", id) }
func (s *AccountService) keyAccountsByClient(id int) string {
	return fmt.Sprintf("accounts:client:%d", id)
}
