package services

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

type DecimalString string

func (d *DecimalString) UnmarshalJSON(data []byte) error {
	value := strings.TrimSpace(string(data))
	if value == "null" {
		return nil
	}

	if strings.HasPrefix(value, "\"") {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		s = strings.TrimSpace(s)
		*d = DecimalString(s)
		return nil
	}

	if value == "" || value == "true" || value == "false" || strings.HasPrefix(value, "{") || strings.HasPrefix(value, "[") {
		return fmt.Errorf("price_from must be a decimal string or number")
	}

	*d = DecimalString(value)
	return nil
}

func (d DecimalString) String() string {
	return string(d)
}

func (d *DecimalString) Scan(value any) error {
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case string:
		*d = DecimalString(v)
	case []byte:
		*d = DecimalString(string(v))
	default:
		return fmt.Errorf("cannot scan decimal value %T", value)
	}
	return nil
}

func (d DecimalString) Value() (driver.Value, error) {
	return string(d), nil
}

type ServiceResponse struct {
	ID               string          `json:"id"`
	ClinicID         string          `json:"clinic_id"`
	Name             string          `json:"name"`
	Description      *string         `json:"description,omitempty"`
	DurationMinutes  *int            `json:"duration_minutes,omitempty"`
	PriceFrom        *DecimalString  `json:"price_from,omitempty"`
	CurrencyCode     string          `json:"currency_code"`
	Benefits         json.RawMessage `json:"benefits"`
	FAQ              json.RawMessage `json:"faq"`
	CommonObjections json.RawMessage `json:"common_objections"`
	IsActive         bool            `json:"is_active"`
}

type CreateServiceRequest struct {
	Name             string          `json:"name"`
	Description      *string         `json:"description,omitempty"`
	DurationMinutes  *int            `json:"duration_minutes,omitempty"`
	PriceFrom        *DecimalString  `json:"price_from,omitempty"`
	Benefits         json.RawMessage `json:"benefits"`
	FAQ              json.RawMessage `json:"faq"`
	CommonObjections json.RawMessage `json:"common_objections"`
}

type UpdateServiceRequest struct {
	Name             string          `json:"name"`
	Description      *string         `json:"description,omitempty"`
	DurationMinutes  *int            `json:"duration_minutes,omitempty"`
	PriceFrom        *DecimalString  `json:"price_from,omitempty"`
	Benefits         json.RawMessage `json:"benefits"`
	FAQ              json.RawMessage `json:"faq"`
	CommonObjections json.RawMessage `json:"common_objections"`
	IsActive         bool            `json:"is_active"`
}

type ErrorResponse struct {
	Error APIError `json:"error"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
