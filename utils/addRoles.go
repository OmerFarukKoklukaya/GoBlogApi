package utils

import (
	"blog/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/uptrace/bun"
	"os"
)

func addAdmin(db *bun.DB, allPermissions []models.Permission) error {
	var roleAdmin models.Role
	err := db.NewSelect().Model(&roleAdmin).Relation("Permissions").Where("name = ?", "admin").Scan(context.Background())
	if err != nil {
		err = nil
		roleAdmin.Name = "admin"
		db.NewInsert().Model(&roleAdmin).Exec(context.Background())
	}

	for _, permission := range allPermissions {
		b := true
		for _, rolePermission := range roleAdmin.Permissions {
			if permission.Name == rolePermission.Name {
				b = false
				break
			}
		}
		if b {
			_, err = db.NewInsert().Model(&models.RoleToPermission{
				RoleID:       roleAdmin.ID,
				PermissionID: permission.ID,
			}).Exec(context.Background())
		}
	}
	if err != nil {
		return err
	}
	return nil

}

func addBasic(db *bun.DB, allPermissions []models.Permission) error {
	var readPermissions []models.Permission
	for _, permission := range allPermissions {
		if permission.Name[0:4] == "read" {
			readPermissions = append(readPermissions, permission)
		}
	}
	if len(readPermissions) == 0 {
		return errors.New("no read permission")
	}

	var roleBasic models.Role

	err := db.NewSelect().Model(&roleBasic).Relation("Permissions").Where("name = ?", "basic").Scan(context.Background())
	if err != nil {
		err = nil
		roleBasic.Name = "basic"
		db.NewInsert().Model(&roleBasic).Exec(context.Background())
	}

	for _, readPermission := range readPermissions {
		b := true
		for _, rolePermission := range roleBasic.Permissions {
			if readPermission.Name == rolePermission.Name {
				fmt.Println("HI")
				b = false
				break
			}
		}
		if b {
			_, err = db.NewInsert().Model(&models.RoleToPermission{
				RoleID:       roleBasic.ID,
				PermissionID: readPermission.ID,
			}).Exec(context.Background())
		}
	}
	if err != nil {
		return err
	}
	return nil
}

func addMininumRequirementForSeeAdminPanel(db *bun.DB, allPermissions []models.Permission) error {

	return nil
}

func AddRoles(db *bun.DB) error {
	allPermissions, err := writePermissions(db)
	if err != nil {
		return err
	}

	err = addAdmin(db, allPermissions)
	if err != nil {
		fmt.Println(err)
		return err
	}

	addBasic(db, allPermissions)

	//addMininumRequirementForSeeAdminPanel()
	return nil
}

func writePermissions(db *bun.DB) ([]models.Permission, error) {
	var permission models.Permission
	var ret []models.Permission
	permissionsJSON, err := os.ReadFile("./permissions.json")
	if err != nil {
		return nil, err
	}
	var permissionMap map[string]string

	err = json.Unmarshal(permissionsJSON, &permissionMap)
	if err != nil {
		return nil, err
	}

	for _, value := range permissionMap {
		err = db.NewSelect().Model(&permission).Where("name = ?", value).Scan(context.Background())
		if err != nil {
			permission.Name = value
			db.NewInsert().Model(&permission).Exec(context.Background())
		}
		err = nil
		ret = append(ret, permission)
		permission.Name = ""
		permission.ID = 0
	}
	fmt.Println("RET: ", ret)

	return ret, nil
}
