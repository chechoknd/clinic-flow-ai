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

func (s *Service) PriorityActions(ctx context.Context, clinicID string, limit int) (ActionsResponse, error) {
	clinicID = strings.TrimSpace(clinicID)
	if clinicID == "" {
		return ActionsResponse{}, ErrMissingClinicID
	}
	if limit <= 0 {
		limit = 4
	}
	if limit > 20 {
		limit = 20
	}

	actions, err := s.repository.PriorityActions(ctx, clinicID, limit)
	if err != nil {
		return ActionsResponse{}, err
	}

	res := ActionsResponse{Data: make([]ActionItemDTO, len(actions))}
	for i, action := range actions {
		var nextActionAt *string
		if action.NextActionAt != nil {
			formatted := action.NextActionAt.Format("2006-01-02T15:04:05Z07:00")
			nextActionAt = &formatted
		}

		res.Data[i] = ActionItemDTO{
			Type:         action.Type,
			Tone:         action.Tone,
			Priority:     action.Priority,
			LeadID:       action.LeadID,
			FullName:     action.FullName,
			Phone:        action.Phone,
			ServiceID:    action.ServiceID,
			ServiceName:  action.ServiceName,
			Status:       action.Status,
			Source:       action.Source,
			Reason:       action.Reason,
			NextActionAt: nextActionAt,
			CreatedAt:    action.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return res, nil
}
