package leads

import (
	"time"
)

type LeadResponse struct {
	ID           string     `json:"id"`
	FullName     string     `json:"full_name"`
	Phone        string     `json:"phone"`
	ServiceID    *string    `json:"service_id,omitempty"`
	ServiceName  *string    `json:"service_name,omitempty"`
	Status       string     `json:"status"`
	Source       string     `json:"source"`
	NextActionAt *time.Time `json:"next_action_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type LeadDetailResponse struct {
	ID           string           `json:"id"`
	FullName     string           `json:"full_name"`
	Phone        string           `json:"phone"`
	Service      *LeadServiceDTO  `json:"service,omitempty"`
	Status       string           `json:"status"`
	Source       string           `json:"source"`
	Notes        []LeadNoteResponse `json:"notes"`
	NextActionAt *time.Time       `json:"next_action_at,omitempty"`
	CreatedAt    time.Time        `json:"created_at"`
}

type LeadServiceDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type LeadNoteResponse struct {
	ID        string    `json:"id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateLeadRequest struct {
	FullName     string     `json:"full_name"`
	Phone        string     `json:"phone"`
	ServiceID    *string    `json:"service_id,omitempty"`
	Status       string     `json:"status"`
	Source       string     `json:"source"`
	Notes        string     `json:"notes,omitempty"`
	NextActionAt *time.Time `json:"next_action_at,omitempty"`
}

type UpdateLeadRequest struct {
	Status       string     `json:"status"`
	Note         string     `json:"note,omitempty"`
	NextActionAt *time.Time `json:"next_action_at,omitempty"`
}

type PaginatedLeadsResponse struct {
	Data       []LeadResponse     `json:"data"`
	Pagination PaginationResponse `json:"pagination"`
}

type PaginationResponse struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type ErrorResponse struct {
	Error APIError `json:"error"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
