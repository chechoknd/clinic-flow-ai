package services

type ServiceEntity struct {
	ID               string
	ClinicID         string
	Name             string
	Description      *string
	DurationMinutes  *int
	PriceFrom        *float64
	CurrencyCode     string
	Benefits         []byte
	FAQ              []byte
	CommonObjections []byte
	IsActive         bool
}
