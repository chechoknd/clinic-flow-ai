package ai

import (
	"fmt"
)

type Context struct {
	ClinicName        string
	ClinicType        string
	City              string
	CommunicationTone string
	ServiceName       string
	ServicePriceFrom  string
	ServiceBenefits   string
	ServiceFAQ        string
	CommonObjections  string
	PatientName       string
	LeadNotes         string
}

func BuildSystemPrompt(ctx Context) string {
	return fmt.Sprintf(`Eres un asistente de comunicación comercial para la clínica "%s" (%s) ubicada en %s. 
Tu tono de voz es %s. 

REGLAS DE SEGURIDAD (CRÍTICO):
- NO diagnostiques.
- NO interpretes síntomas.
- NO recomiendes medicamentos ni dosis.
- NO garantices resultados médicos.
- Siempre anima a agendar una valoración profesional.
- Compórtate como un asistente comercial, NO como un médico.

CONTEXTO DEL SERVICIO:
Servicio: %s
Precio desde: %s
Beneficios: %s
Preguntas Frecuentes: %s
Objeciones Comunes: %s

DATOS DEL PROSPECTO:
Nombre: %s
Notas Previas: %s

Tu objetivo es responder de forma persuasiva pero segura, enfocándote en que el paciente asista a una consulta de valoración.`,
		ctx.ClinicName, ctx.ClinicType, ctx.City, ctx.CommunicationTone,
		ctx.ServiceName, ctx.ServicePriceFrom, ctx.ServiceBenefits, ctx.ServiceFAQ, ctx.CommonObjections,
		ctx.PatientName, ctx.LeadNotes,
	)
}

func BuildReplySuggestionUserPrompt(patientMessage, desiredTone string) string {
	tonePrompt := ""
	if desiredTone != "" {
		tonePrompt = fmt.Sprintf(" Ajusta el tono a: %s.", desiredTone)
	}

	return fmt.Sprintf(`El paciente dice: "%s"

Genera 4 variantes de respuesta en formato JSON con la siguiente estructura:
{
  "short": "Respuesta muy breve y directa",
  "persuasive": "Respuesta enfocada en beneficios y agendar",
  "technical": "Respuesta que menciona aspectos del proceso sin diagnosticar",
  "closing_question": "Una pregunta de cierre para incentivar la respuesta"
}

Asegúrate de que el JSON sea válido y no incluyas texto fuera del objeto.%s`, patientMessage, tonePrompt)
}

func BuildObjectionHandlerUserPrompt(objection string) string {
	return fmt.Sprintf(`El paciente tiene la siguiente objeción: "%s"

Genera una estrategia de manejo de objeción y una respuesta sugerida en formato JSON:
{
  "objection_type": "Categoría de la objeción (precio, miedo, tiempo, etc.)",
  "recommended_strategy": "Breve descripción de la estrategia comercial a usar",
  "suggested_message": "Respuesta persuasiva para el paciente",
  "closing_question": "Pregunta de cierre para mantener la conversación"
}

Asegúrate de que el JSON sea válido y no incluyas texto fuera del objeto.`, objection)
}

func BuildFollowUpMessageUserPrompt(lastContactNote string) string {
	if lastContactNote == "" {
		lastContactNote = "No hay una nota adicional del último contacto."
	}

	return fmt.Sprintf(`El prospecto no ha avanzado y necesita seguimiento manual por WhatsApp.
Última nota de contacto: "%s"

Genera una respuesta de reenganche en formato JSON con la siguiente estructura:
{
  "suggested_message": "Mensaje breve, amable y listo para copiar en WhatsApp",
  "recommended_timing": "Cuándo conviene enviarlo",
  "next_step": "Siguiente acción comercial recomendada"
}

El mensaje debe sonar humano, no insistente, y debe invitar a agendar una valoración profesional. No incluyas texto fuera del objeto JSON.`, lastContactNote)
}
