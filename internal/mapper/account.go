package mapper

import (
	"basic-gin/internal/dto"
	"basic-gin/internal/model"
)

type AccountMapper struct{}

func NewAccountMapper() *AccountMapper { return &AccountMapper{} }

func (m *AccountMapper) ToResponse(a *model.Account) *dto.AccountResponse {
	return &dto.AccountResponse{
		ID:            a.ID,
		ClientID:      a.ClientId,
		AccountNumber: a.AccountNumber,
		Balance:       a.Balance,
		CreatedAt:     a.CreatedAt,
	}
}

func (m *AccountMapper) ToEntityFromUpdate(existing model.Account, in dto.AccountUpdate) model.Account {
	out := existing

	if in.ClientID != nil {
		out.ClientId = *in.ClientID
	}
	if in.AccountNumber != nil {
		out.AccountNumber = *in.AccountNumber
	}
	if in.Balance != nil {
		out.Balance = *in.Balance
	}
	return out
}

func (m *AccountMapper) ToResponseSlice(items []*model.Account) []*dto.AccountResponse {
	if len(items) == 0 {
		return []*dto.AccountResponse{}
	}
	res := make([]*dto.AccountResponse, 0, len(items))
	for _, a := range items {
		res = append(res, m.ToResponse(a))
	}
	return res
}
