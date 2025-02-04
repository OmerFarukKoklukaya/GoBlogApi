package models

type Role struct {
	ID          uint         `bun:"id,pk,autoincrement" json:"id"`
	Name        string       `bun:"name" json:"name"`
	Permissions []Permission `bun:"m2m:role_to_permissions,join:Role=Permission"`
}

type Permission struct {
	ID   uint   `bun:"id,pk,autoincrement" json:"id"`
	Name string `bun:"name" json:"name"`
}

type RoleToPermission struct {
	Role         *Role       `bun:"rel:belongs-to,join:role_id=id"`
	RoleID       uint        `bun:"role_id" json:"roleID"`
	Permission   *Permission `bun:"rel:belongs-to,join:permission_id=id"`
	PermissionID uint        `bun:"permission_id" json:"permissionID"`
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
