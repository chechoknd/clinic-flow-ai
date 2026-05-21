package clinics

import "encoding/json"

type ClinicResponse struct {
	ID                string          `json:"id"`
	Name              string          `json:"name"`
	ClinicType        string          `json:"clinic_type"`
	City              string          `json:"city"`
	Phone             *string         `json:"phone,omitempty"`
	WhatsApp          string          `json:"whatsapp"`
	Address           *string         `json:"address,omitempty"`
	OpeningHours      json.RawMessage `json:"opening_hours"`
	GeneralFAQ        json.RawMessage `json:"general_faq"`
	CommunicationTone string          `json:"communication_tone"`
}

type UpdateClinicRequest struct {
	Name              string          `json:"name"`
	City              string          `json:"city"`
	Phone             *string         `json:"phone,omitempty"`
	WhatsApp          string          `json:"whatsapp"`
	Address           *string         `json:"address,omitempty"`
	OpeningHours      json.RawMessage `json:"opening_hours"`
	GeneralFAQ        json.RawMessage `json:"general_faq"`
	CommunicationTone string          `json:"communication_tone"`
}

type ErrorResponse struct {
	Error APIError `json:"error"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
