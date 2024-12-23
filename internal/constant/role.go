package constant

const (
	// Role IDs
	RoleIDAdmin   int64 = 1
	RoleIDStaff   int64 = 2
	RoleIDCoach   int64 = 3
	RoleIDCoachee int64 = 4
)

var roleToID = map[string]int64{
	"admin":   RoleIDAdmin,
	"staff":   RoleIDStaff,
	"coach":   RoleIDCoach,
	"coachee": RoleIDCoachee,
}

func GetRoleID(role string) int64 {
	if id, exists := roleToID[role]; exists {
		return id
	}
	return 0
}
