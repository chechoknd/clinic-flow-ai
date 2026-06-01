package dashboard

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
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

type dayInterval struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type hoursConfig struct {
	Monday    []dayInterval `json:"monday"`
	Tuesday   []dayInterval `json:"tuesday"`
	Wednesday []dayInterval `json:"wednesday"`
	Thursday  []dayInterval `json:"thursday"`
	Friday    []dayInterval `json:"friday"`
	Saturday  []dayInterval `json:"saturday"`
	Sunday    []dayInterval `json:"sunday"`
}

func (s *Service) ScheduleSummary(ctx context.Context, clinicID string, dateStr string) (ScheduleSummaryResponse, error) {
	clinicID = strings.TrimSpace(clinicID)
	if clinicID == "" {
		return ScheduleSummaryResponse{}, ErrMissingClinicID
	}
	dateStr = strings.TrimSpace(dateStr)
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	}

	summary, err := s.repository.ScheduleSummary(ctx, clinicID, dateStr)
	if err != nil {
		return ScheduleSummaryResponse{}, err
	}

	// Calculate available slots
	availableSlots := 0
	profs, err := s.repository.GetActiveProfessionalsWithHours(ctx, clinicID)
	if err == nil {
		appts, err := s.repository.GetAppointmentsForDate(ctx, clinicID, dateStr)
		if err == nil {
			t, err := time.Parse("2006-01-02", dateStr)
			if err == nil {
				dayOfWeek := strings.ToLower(t.Weekday().String())
				for _, prof := range profs {
					if len(prof.WorkingHours) == 0 || string(prof.WorkingHours) == "{}" {
						continue
					}
					var hc hoursConfig
					if err := json.Unmarshal(prof.WorkingHours, &hc); err != nil {
						continue
					}
					intervals := getIntervals(hc, dayOfWeek)
					for _, interval := range intervals {
						startT, err := parseTime(t, interval.Start)
						if err != nil {
							continue
						}
						endT, err := parseTime(t, interval.End)
						if err != nil {
							continue
						}
						// candidate slots of 60 minutes
						candidate := startT
						for {
							candidateEnd := candidate.Add(60 * time.Minute)
							if candidateEnd.After(endT) {
								break
							}
							// check overlap
							overlap := false
							for _, appt := range appts {
								if appt.ProfessionalID == prof.ID {
									if appt.StartsAt.Before(candidateEnd) && appt.EndsAt.After(candidate) {
										overlap = true
										break
									}
								}
							}
							if !overlap {
								availableSlots++
							}
							candidate = candidate.Add(60 * time.Minute)
						}
					}
				}
			}
		}
	}

	profCounts := make([]ProfessionalCount, len(summary.AppointmentsByProfessional))
	for i, pc := range summary.AppointmentsByProfessional {
		profCounts[i] = ProfessionalCount{
			ProfessionalID:   pc.ProfessionalID,
			ProfessionalName: pc.ProfessionalName,
			AppointmentCount: pc.AppointmentCount,
		}
	}

	serviceCounts := make([]ServiceCount, len(summary.TopServicesByScheduleDemand))
	for i, sc := range summary.TopServicesByScheduleDemand {
		serviceCounts[i] = ServiceCount{
			ServiceID:        sc.ServiceID,
			ServiceName:      sc.ServiceName,
			AppointmentCount: sc.AppointmentCount,
		}
	}

	return ScheduleSummaryResponse{
		Date:                            summary.Date,
		TodaysAppointments:              summary.TodaysAppointments,
		AppointmentsPendingConfirmation: summary.AppointmentsPendingConfirmation,
		AvailableSlots:                  availableSlots,
		HotLeadsWithoutAppointment:      summary.HotLeadsWithoutAppointment,
		OverdueFollowUps:                summary.OverdueFollowUps,
		AppointmentsByProfessional:      profCounts,
		LeadsConvertedToAppointments:    summary.LeadsConvertedToAppointments,
		TopServicesByScheduleDemand:     serviceCounts,
	}, nil
}

func getIntervals(hc hoursConfig, day string) []dayInterval {
	switch day {
	case "monday":
		return hc.Monday
	case "tuesday":
		return hc.Tuesday
	case "wednesday":
		return hc.Wednesday
	case "thursday":
		return hc.Thursday
	case "friday":
		return hc.Friday
	case "saturday":
		return hc.Saturday
	case "sunday":
		return hc.Sunday
	default:
		return nil
	}
}

func parseTime(date time.Time, timeStr string) (time.Time, error) {
	var hour, min int
	_, err := fmt.Sscanf(timeStr, "%d:%d", &hour, &min)
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(date.Year(), date.Month(), date.Day(), hour, min, 0, 0, time.UTC), nil
}
