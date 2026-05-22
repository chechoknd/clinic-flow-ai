package ai

import (
	"errors"
	"fmt"
)

var (
	ErrValidation      = errors.New("validation error")
	ErrAIProvider      = errors.New("ai provider error")
	ErrAIResponse      = errors.New("invalid ai response")
	ErrAISafetyBlocked = errors.New("ai safety blocked")
)

type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	if e.Field == "" {
		return e.Message
	}
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

func (e ValidationError) Is(target error) bool {
	return target == ErrValidation
}

type ProviderError struct {
	Provider string
	Err      error
}

func (e ProviderError) Error() string {
	if e.Provider == "" {
		return "ai provider request failed"
	}
	return fmt.Sprintf("%s provider request failed", e.Provider)
}

func (e ProviderError) Unwrap() error {
	return e.Err
}

func (e ProviderError) Is(target error) bool {
	return target == ErrAIProvider
}

type ResponseError struct {
	Err error
}

func (e ResponseError) Error() string {
	return "ai provider returned an invalid response"
}

func (e ResponseError) Unwrap() error {
	return e.Err
}

func (e ResponseError) Is(target error) bool {
	return target == ErrAIResponse
}

type SafetyBlockedError struct {
	Status string
}

func (e SafetyBlockedError) Error() string {
	return "ai response failed safety validation"
}

func (e SafetyBlockedError) Is(target error) bool {
	return target == ErrAISafetyBlocked
}
