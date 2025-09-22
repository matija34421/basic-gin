package dto

type ClientCreate struct {
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	Email            string `json:"email"`
	ResidenceAddress string `json:"residence_address"`
	BirthDate        string `json:"birth_date"`
}

type ClientUpdate struct {
	ID               int64  `json:"id"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	Email            string `json:"email"`
	ResidenceAddress string `json:"residence_address"`
	BirthDate        string `json:"birth_date"`
}

type ClientResponse struct {
	ID               int64  `json:"id"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	Email            string `json:"email"`
	ResidenceAddress string `json:"residence_address"`
	BirthDate        string `json:"birth_date"`
	CreatedAt        string `json:"created_at"`
}
