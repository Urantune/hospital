package service

const (
	RoleSystemAdmin = "system_admin"
	RoleClinicAdmin = "clinic_admin"
	RoleStaff       = "staff"
	RolePatient     = "patient"
)

const (
	PermissionProfileRead   = "profile.read"
	PermissionProfileUpdate = "profile.update"
	PermissionClinicRead    = "clinic.read"
	PermissionClinicManage  = "clinic.manage"
	PermissionConfigRead    = "config.read"
	PermissionConfigManage  = "config.manage"
	PermissionAuditRead     = "audit.read"
	PermissionCMSSync       = "cms.sync"
)

var rolePermissions = map[string]map[string]bool{
	RoleSystemAdmin: {
		PermissionProfileRead:   true,
		PermissionProfileUpdate: true,
		PermissionClinicRead:    true,
		PermissionClinicManage:  true,
		PermissionConfigRead:    true,
		PermissionConfigManage:  true,
		PermissionAuditRead:     true,
		PermissionCMSSync:       true,
	},
	RoleClinicAdmin: {
		PermissionProfileRead:   true,
		PermissionProfileUpdate: true,
		PermissionClinicRead:    true,
		PermissionClinicManage:  true,
		PermissionConfigRead:    true,
		PermissionConfigManage:  true,
	},
	RoleStaff: {
		PermissionProfileRead:   true,
		PermissionProfileUpdate: true,
		PermissionClinicRead:    true,
		PermissionConfigRead:    true,
	},
	RolePatient: {
		PermissionProfileRead:   true,
		PermissionProfileUpdate: true,
	},
}

func HasPermission(role, permission string) bool {
	permissions, ok := rolePermissions[role]
	if !ok {
		return false
	}

	return permissions[permission]
}
