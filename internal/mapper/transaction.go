package mapper

import (
	"basic-gin/internal/dto"
	"basic-gin/internal/model"
)

func TransactionToResponse(t *model.Transaction) *dto.TransactionResponse {
	return &dto.TransactionResponse{
		ID:            t.ID,
		FromAccountID: t.FromAccountID,
		ToAccountID:   t.ToAccountID,
		Amount:        t.Amount,
		CreatedAt:     t.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
