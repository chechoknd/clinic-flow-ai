package schedule

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrMissingClinicID        = errors.New("clinic id is required")
	ErrProfessionalIDRequired = errors.New("professional_id is required")
	ErrInvalidDateRange       = errors.New("date_to must be after or equal to date_from")
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetAvailability(ctx context.Context, clinicID string, req AvailabilityRequest) (AvailabilityResponse, error) {
	clinicID = strings.TrimSpace(clinicID)
	if clinicID == "" {
		return AvailabilityResponse{}, ErrMissingClinicID
	}
	req.ProfessionalID = strings.TrimSpace(req.ProfessionalID)
	if req.ProfessionalID == "" {
		return AvailabilityResponse{}, ErrProfessionalIDRequired
	}
	if req.DateTo.Before(req.DateFrom) {
		return AvailabilityResponse{}, ErrInvalidDateRange
	}

	// 1. Resolve duration
	durationMins := req.DurationMinutes
	if durationMins <= 0 && req.ServiceID != "" {
		dur, err := s.repository.GetServiceDuration(ctx, clinicID, req.ServiceID)
		if err == nil && dur > 0 {
			durationMins = dur
		}
	}
	if durationMins <= 0 {
		durationMins = 60 // fallback default to 60 minutes
	}
	duration := time.Duration(durationMins) * time.Minute

	// 2. Fetch professional and working hours
	prof, err := s.repository.FindProfessional(ctx, clinicID, req.ProfessionalID)
	if err != nil {
		return AvailabilityResponse{}, err
	}

	var wh WorkingHours
	if len(prof.WorkingHours) > 0 && string(prof.WorkingHours) != "{}" {
		if err := json.Unmarshal(prof.WorkingHours, &wh); err != nil {
			wh = WorkingHours{}
		}
	}

	// 3. Fetch active appointments in range
	// We want to fetch from the start of DateFrom to the end of DateTo.
	// Make sure we encompass the whole time range.
	startRange := time.Date(req.DateFrom.Year(), req.DateFrom.Month(), req.DateFrom.Day(), 0, 0, 0, 0, req.DateFrom.Location())
	endRange := time.Date(req.DateTo.Year(), req.DateTo.Month(), req.DateTo.Day(), 23, 59, 59, 999999999, req.DateTo.Location())

	activeAppts, err := s.repository.FindActiveAppointments(ctx, clinicID, req.ProfessionalID, startRange, endRange)
	if err != nil {
		return AvailabilityResponse{}, err
	}

	// 4. Generate slots day by day
	loc := req.DateFrom.Location()
	var slots []TimeSlot

	// Limit calculation to a maximum range of 31 days to prevent abuse/performance issues
	maxDays := 31
	currentDay := time.Date(req.DateFrom.Year(), req.DateFrom.Month(), req.DateFrom.Day(), 0, 0, 0, 0, loc)
	endDayLimit := time.Date(req.DateTo.Year(), req.DateTo.Month(), req.DateTo.Day(), 0, 0, 0, 0, loc)

	for dayIdx := 0; dayIdx <= maxDays; dayIdx++ {
		if currentDay.After(endDayLimit) {
			break
		}

		dayOfWeek := strings.ToLower(currentDay.Weekday().String())
		intervals := getIntervalsForDay(wh, dayOfWeek)

		for _, interval := range intervals {
			intervalStart, err := parseTimeOfDay(currentDay, interval.Start, loc)
			if err != nil {
				continue
			}
			intervalEnd, err := parseTimeOfDay(currentDay, interval.End, loc)
			if err != nil {
				continue
			}

			// Generate candidate slots within the interval, stepping by 30 minutes
			candidateStart := intervalStart
			for {
				candidateEnd := candidateStart.Add(duration)
				if candidateEnd.After(intervalEnd) {
					break
				}

				// Check overlap
				hasOverlap := false
				for _, appt := range activeAppts {
					if appt.StartsAt.Before(candidateEnd) && appt.EndsAt.After(candidateStart) {
						hasOverlap = true
						break
					}
				}

				if !hasOverlap {
					slots = append(slots, TimeSlot{
						StartsAt: candidateStart,
						EndsAt:   candidateEnd,
					})
				}

				candidateStart = candidateStart.Add(30 * time.Minute)
			}
		}

		currentDay = currentDay.AddDate(0, 0, 1)
	}

	return AvailabilityResponse{
		ProfessionalID: prof.ID,
		Slots:          slots,
	}, nil
}

func getIntervalsForDay(wh WorkingHours, day string) []DayWorkingInterval {
	switch day {
	case "monday":
		return wh.Monday
	case "tuesday":
		return wh.Tuesday
	case "wednesday":
		return wh.Wednesday
	case "thursday":
		return wh.Thursday
	case "friday":
		return wh.Friday
	case "saturday":
		return wh.Saturday
	case "sunday":
		return wh.Sunday
	default:
		return nil
	}
}

func parseTimeOfDay(date time.Time, timeStr string, loc *time.Location) (time.Time, error) {
	parts := strings.Split(timeStr, ":")
	if len(parts) < 2 {
		return time.Time{}, fmt.Errorf("invalid time format: %s", timeStr)
	}
	var hour, min int
	_, err := fmt.Sscanf(timeStr, "%d:%d", &hour, &min)
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(date.Year(), date.Month(), date.Day(), hour, min, 0, 0, loc), nil
}
