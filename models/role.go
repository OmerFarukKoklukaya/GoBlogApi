package models

type Role struct {
	ID          int          `bun:"id,pk,autoincrement" json:"id"`
	Name        string       `bun:"name" json:"name"`
	Permissions []Permission `bun:"m2m:role_to_permissions,join:Role=Permission"`
}

type Permission struct {
	ID   int    `bun:"id,pk,autoincrement" json:"id"`
	Name string `bun:"name" json:"name"`
}

type RoleToPermission struct {
	Role         *Role       `bun:"rel:belongs-to,join:role_id=id"`
	RoleID       int         `bun:"role_id" json:"role_id"`
	Permission   *Permission `bun:"rel:belongs-to,join:permission_id=id"`
	PermissionID int         `bun:"permission_id" json:"permission_id"`
}

func (role Role) IsHaveEditPermission() bool {
	for _, permission := range role.Permissions {
		if permission.Name[:4] == "edit" {
			return true
		}
	}
	return false
}

func (role Role) IsHavePermission(permissionName string) bool {
	for _, permission := range role.Permissions {
		if permission.Name == permissionName {
			return true
		}
	}
	return false
}
