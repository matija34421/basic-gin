package handler

import (
	"basic-gin/internal/service"
)

type Dependencies struct {
	AccountHandler     *AccountHandler
	ClientHandler      *ClientHandler
	TransactionHandler *TransactionHandler
}

func NewDependencies(cs *service.ClientService, as *service.AccountService, ts *service.TransactionService) *Dependencies {
	var ch *ClientHandler
	if cs != nil {
		ch = NewClientHandler(cs)
	}
	var ah *AccountHandler
	if as != nil {
		ah = NewAccountHandler(as)
	}
	var th *TransactionHandler
	if ts != nil {
		th = NewTransactionHandler(ts)
	}
	return &Dependencies{
		ClientHandler:      ch,
		AccountHandler:     ah,
		TransactionHandler: th,
	}
}
