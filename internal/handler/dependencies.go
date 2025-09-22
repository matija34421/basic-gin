package handler

import (
	"basic-gin/internal/service"
)

type Dependencies struct {
	AccountHandler *AccountHandler
	ClientHandler  *ClientHandler
}

func NewDependencies(cs *service.ClientService, as *service.AccountService) *Dependencies {
	var ch *ClientHandler
	if cs != nil {
		ch = NewClientHandler(cs)
	}
	var ah *AccountHandler
	if as != nil {
		ah = NewAccountHandler(as)
	}
	return &Dependencies{
		ClientHandler:  ch,
		AccountHandler: ah,
	}
}
