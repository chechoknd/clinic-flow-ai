package ai

type GenerationRecord struct {
	ID              string
	ClinicID        string
	UserID          string
	Feature         string
	Provider        string
	Model           string
	Status          string
	SafetyStatus    string
	ErrorCode       string
	InputCharCount  int
	OutputCharCount int
}
