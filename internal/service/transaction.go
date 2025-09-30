package service

import (
	"basic-gin/internal/dto"
	"basic-gin/internal/mapper"
	"basic-gin/internal/model"
	"basic-gin/internal/repository"
	"context"
	"errors"
	"time"

	"golang.org/x/sync/errgroup"
)

type TransactionService struct {
	transactionRepository repository.TransactionRepository
	accountRepository     repository.AccountRepository
	accountService        AccountService
}

func NewTransactionService(
	transactionRepository *repository.TransactionRepository,
	accountRepository *repository.AccountRepository,
	accountService AccountService,
) *TransactionService {
	return &TransactionService{
		transactionRepository: *transactionRepository,
		accountRepository:     *accountRepository,
		accountService:        accountService,
	}
}

func (s *TransactionService) CreateTransfer(ctx context.Context, in dto.TransactionCreate) (*dto.TransactionResponse, error) {
	if in.FromAccountID == in.ToAccountID {
		return nil, errors.New("from and to accounts must differ")
	}
	if in.Amount <= 0 {
		return nil, errors.New("amount must be positive")
	}

	g, ctx := errgroup.WithContext(ctx)

	var preFrom, preTo *dto.AccountResponse

	g.Go(func() error {
		acc, err := s.accountService.GetById(ctx, in.ToAccountID)
		if err != nil {
			return err
		}
		preFrom = acc
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	g.Go(func() error {
		acc, err := s.accountService.GetById(ctx, in.FromAccountID)
		if err != nil {
			return err
		}
		preTo = acc
		return nil
	})

	if preFrom.Balance < in.Amount {
		return nil, errors.New("insufficient funds")
	}

	//da izbegnem unused za preTo
	_ = preTo

	tx, err := s.accountRepository.Pool().Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	first, second := in.FromAccountID, in.ToAccountID
	if first > second {
		first, second = second, first
	}

	if _, err := s.accountRepository.GetByIdTx(ctx, tx, first, true); err != nil {
		return nil, err
	}
	if _, err := s.accountRepository.GetByIdTx(ctx, tx, second, true); err != nil {
		return nil, err
	}

	fromAcc, err := s.accountRepository.GetByIdTx(ctx, tx, in.FromAccountID, true)
	if err != nil {
		return nil, err
	}
	if fromAcc.Balance < in.Amount {
		return nil, errors.New("insufficient funds")
	}

	if _, err := s.accountRepository.UpdateBalanceDeltaTx(ctx, tx, in.FromAccountID, -in.Amount); err != nil {
		return nil, err
	}
	if _, err := s.accountRepository.UpdateBalanceDeltaTx(ctx, tx, in.ToAccountID, +in.Amount); err != nil {
		return nil, err
	}

	t := &model.Transaction{
		FromAccountID: in.FromAccountID,
		ToAccountID:   in.ToAccountID,
		Amount:        in.Amount,
	}
	if err := s.transactionRepository.SaveTx(ctx, tx, t); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	// 7) KYT simulacija – klijent čeka 10s
	time.Sleep(10 * time.Second)

	return mapper.TransactionToResponse(t), nil
}

func (s *TransactionService) ListByAccountID(ctx context.Context, accountID, limit, offset int) ([]*dto.TransactionResponse, error) {
	if limit <= 0 || limit > 1000 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	items, err := s.transactionRepository.ListByAccountID(ctx, accountID, limit, offset)
	if err != nil {
		return nil, err
	}
	out := make([]*dto.TransactionResponse, 0, len(items))
	for _, t := range items {
		out = append(out, mapper.TransactionToResponse(t))
	}
	return out, nil
}
