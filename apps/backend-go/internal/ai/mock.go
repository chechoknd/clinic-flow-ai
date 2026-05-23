package ai

import (
	"context"
	"strings"
)

type MockProvider struct{}

func NewMockProvider() *MockProvider {
	return &MockProvider{}
}

func (p *MockProvider) Name() string {
	return "mock"
}

func (p *MockProvider) Generate(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	prompt := strings.ToLower(userPrompt)
	switch {
	case strings.Contains(prompt, "objec"):
		return `{"objection_type":"precio","recommended_strategy":"Validar interes y reforzar valor sin presionar.","suggested_message":"Entiendo tu inquietud. Podemos revisar opciones y resolver tus dudas antes de tomar una decision.","closing_question":"Quieres que te compartamos los pasos para una valoracion?"}`, nil
	case strings.Contains(prompt, "seguimiento") || strings.Contains(prompt, "follow"):
		return `{"suggested_message":"Hola, queriamos saber si pudiste revisar la informacion. Estamos atentos para ayudarte con el siguiente paso.","recommended_timing":"Hoy en horario laboral","next_step":"Enviar mensaje manual por WhatsApp"}`, nil
	default:
		return `{"short":"Claro, con gusto te ayudamos con la informacion.","persuasive":"Podemos orientarte y revisar la mejor opcion segun tu necesidad.","technical":"Primero realizamos una valoracion para confirmar el plan adecuado.","closing_question":"Quieres que te ayudemos a coordinar una valoracion?"}`, nil
	}
}
