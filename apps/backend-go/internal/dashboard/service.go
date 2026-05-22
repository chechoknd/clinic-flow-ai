package dashboard

import (
	"context"
	"errors"
	"strings"
)

var ErrMissingClinicID = errors.New("clinic id is required")

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Summary(ctx context.Context, clinicID string) (SummaryResponse, error) {
	clinicID = strings.TrimSpace(clinicID)
	if clinicID == "" {
		return SummaryResponse{}, ErrMissingClinicID
	}

	summary, err := s.repository.Summary(ctx, clinicID)
	if err != nil {
		return SummaryResponse{}, err
	}

	topServices := make([]TopServiceDTO, len(summary.TopServices))
	for i, service := range summary.TopServices {
		topServices[i] = TopServiceDTO{
			ServiceID:   service.ServiceID,
			ServiceName: service.ServiceName,
			LeadCount:   service.LeadCount,
		}
	}

	return SummaryResponse{
		LeadsTotal:            summary.LeadsTotal,
		LeadsByStatus:         summary.LeadsByStatus,
		TopServices:           topServices,
		PendingFollowUpsToday: summary.PendingFollowUpsToday,
		OverdueFollowUps:      summary.OverdueFollowUps,
		UpcomingFollowUps:     summary.UpcomingFollowUps,
		ConversionRate:        summary.ConversionRate,
	}, nil
}
