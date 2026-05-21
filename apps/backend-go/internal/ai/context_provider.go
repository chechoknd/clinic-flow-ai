package ai

import (
	"context"
)

type ClinicContext struct {
	Name              string
	ClinicType        string
	City              string
	CommunicationTone string
}

type ServiceContext struct {
	Name             string
	Benefits         string
	FAQ              string
	CommonObjections string
}

type LeadContext struct {
	FullName string
	Notes    string
}

type ContextProvider interface {
	GetContext(ctx context.Context, clinicID, serviceID, leadID string) (Context, error)
}
