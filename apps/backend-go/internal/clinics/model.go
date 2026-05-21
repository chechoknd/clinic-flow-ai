package clinics

type Clinic struct {
	ID                string
	Name              string
	ClinicType        string
	City              string
	Phone             *string
	WhatsApp          string
	Address           *string
	OpeningHours      []byte
	GeneralFAQ        []byte
	CommunicationTone string
}
