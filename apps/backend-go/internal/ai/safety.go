package ai

import (
	"strings"
)

var riskyKeywords = []string{
	"paracetamol", "ibuprofeno", "antibiótico", "antibiotico", "amoxicilina", "dosis",
	"diagnóstico", "diagnostico", "tienes", "sufres de", "enfermedad", "receta",
	"curación total", "curacion total", "garantizado", "sin riesgo",
}

func ValidateSafety(text string) (string, bool) {
	lower := strings.ToLower(text)
	for _, kw := range riskyKeywords {
		if strings.Contains(lower, kw) {
			return "passed_with_warnings", false // For MVP, let's just flag it
		}
	}
	return "passed", true
}
