package mapper

import (
	"basic-gin/internal/dto"
	"basic-gin/internal/model"
	"fmt"
	"strings"
	"time"
)

const dateLayout = "2006-01-02"

func ClientToResponse(c *model.Client) dto.ClientResponse {
	return dto.ClientResponse{
		ID:               c.ID,
		FirstName:        c.FirstName,
		LastName:         c.LastName,
		Email:            c.Email,
		ResidenceAddress: c.ResidenceAddress,
		BirthDate:        c.BirthDate.Format("2006-01-02"),
		CreatedAt:        c.CreatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
	}
}

func ClientsToResponseSlice(items []*model.Client) []dto.ClientResponse {
	out := make([]dto.ClientResponse, 0, len(items))
	for _, c := range items {
		out = append(out, ClientToResponse(c))
	}
	return out
}

func ToClientFromCreate(in dto.ClientCreate) (model.Client, error) {
	bd, err := time.Parse(dateLayout, strings.TrimSpace(in.BirthDate))
	if err != nil {
		return model.Client{}, fmt.Errorf("invalid birth_date (use YYYY-MM-DD): %w", err)
	}
	return model.Client{
		FirstName:        strings.TrimSpace(in.FirstName),
		LastName:         strings.TrimSpace(in.LastName),
		Email:            strings.TrimSpace(in.Email),
		ResidenceAddress: strings.TrimSpace(in.ResidenceAddress),
		BirthDate:        bd,
	}, nil
}

func ToClientFromUpdate(in dto.ClientUpdate) (model.Client, error) {
	if in.ID <= 0 {
		return model.Client{}, fmt.Errorf("invalid id")
	}
	bd, err := time.Parse(dateLayout, strings.TrimSpace(in.BirthDate))
	if err != nil {
		return model.Client{}, fmt.Errorf("invalid birth_date (use YYYY-MM-DD): %w", err)
	}
	return model.Client{
		ID:               in.ID,
		FirstName:        strings.TrimSpace(in.FirstName),
		LastName:         strings.TrimSpace(in.LastName),
		Email:            strings.TrimSpace(in.Email),
		ResidenceAddress: strings.TrimSpace(in.ResidenceAddress),
		BirthDate:        bd,
	}, nil
}
