package mapper

import (
	"basic-gin/internal/dto"
	"basic-gin/internal/model"
	"fmt"
	"strings"
	"time"
)

type ClientMapper struct {
}

func NewClientMapper() *ClientMapper {
	return &ClientMapper{}
}

const dateLayout = "2006-01-02"

func (m *ClientMapper) ToResponse(c *model.Client) dto.ClientResponse {
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

func (m *ClientMapper) ToResponseSlice(items []*model.Client) []dto.ClientResponse {
	out := make([]dto.ClientResponse, 0, len(items))
	for _, c := range items {
		out = append(out, m.ToResponse(c))
	}
	return out
}

func (m *ClientMapper) ToEntity(in dto.ClientCreate) (model.Client, error) {
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

func (m *ClientMapper) ToEntityFromUpdate(in dto.ClientUpdate) (model.Client, error) {
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
