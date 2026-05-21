package auth

const (
	RoleSuperadmin  = "superadmin"
	RoleClinicAdmin = "clinic_admin"
	RoleAssistant   = "assistant"
	TokenTypeBearer = "Bearer"

	defaultTokenTTLSeconds = 3600
)

type User struct {
	ID           string
	ClinicID     string
	FullName     string
	Email        string
	PasswordHash string
	Role         string
	IsActive     bool
}
