package constant

const (
	RoleIDAdmin     int64 = 1
	RoleIDStaff     int64 = 2
	RoleIDCoach     int64 = 3
	RoleIDCoachee   int64 = 4
	RoleIDMarketing int64 = 5
)

var roleToID = map[string]int64{
	"Admin":     RoleIDAdmin,
	"Staff":     RoleIDStaff,
	"Coach":     RoleIDCoach,
	"Coachee":   RoleIDCoachee,
	"Marketing": RoleIDMarketing,
}

var roleIdToName = map[int64]string{
	RoleIDAdmin:     "Admin",
	RoleIDStaff:     "Staff",
	RoleIDCoach:     "Coach",
	RoleIDCoachee:   "Coachee",
	RoleIDMarketing: "Marketing",
}

var backofficeRoles = map[int64]struct{}{
	RoleIDAdmin:     {},
	RoleIDStaff:     {},
	RoleIDMarketing: {},
}

var dashboardRoles = map[int64]struct{}{
	RoleIDCoach:   {},
	RoleIDCoachee: {},
}

func IsBackofficeUser(roleID int64) bool {
	_, exists := backofficeRoles[roleID]
	return exists
}

// IsDashboardUser checks if a role belongs to the dashboard category.
func IsDashboardUser(roleID int64) bool {
	_, exists := dashboardRoles[roleID]
	return exists
}

func GetRoleID(role string) int64 {
	if id, exists := roleToID[role]; exists {
		return id
	}
	return 0
}

func GetRoleName(roleId *int64) string {
	if roleId == nil {
		return ""
	}
	if name, exists := roleIdToName[*roleId]; exists {
		return name
	}
	return ""
}
