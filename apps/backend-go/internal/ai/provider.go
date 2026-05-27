package ai

import (
	"context"
)

type Provider interface {
	Generate(ctx context.Context, systemPrompt, userPrompt string) (string, error)
	Name() string
	Model() string
}
