package controllers

import (
	"blog/database"
	"blog/models"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func SelectRole(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	db := database.DB
	var role models.Role
	db.NewSelect().Model(&role).Relation("Permissions").Where("id = ?", id).Scan(ctx)

	return c.JSON(models.ViewData{Data: role})
}

func SelectRoles(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	if page <= 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	limit := c.QueryInt("limit", 20)
	if limit <= 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	offset := limit * (page - 1)

	db := database.DB
	var roles []models.Role
	roleLen, _ := db.NewSelect().Model(&roles).Relation("Permissions").Limit(limit).Offset(offset).ScanAndCount(ctx)

	paging := models.Pagination{}.Paginate(page, roleLen, limit)

	return c.JSON(models.ViewData{
		Data: roles,
		Meta: models.Meta{
			Pagination: paging,
		},
	})
}

func AddRole(c *fiber.Ctx) error {
	var roleAndPermissions = struct {
		Name        string `json:"name"`
		Permissions []int  `json:"permissions"`
	}{}
	err := json.Unmarshal(c.Body(), &roleAndPermissions)

	if err != nil {
		fmt.Println(err)
	}
	if len(roleAndPermissions.Name) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON("Role name cannot be empty")
	}
	var role = models.Role{Name: roleAndPermissions.Name}
	db := database.DB
	db.NewInsert().Model(&role).Exec(ctx)
	var rtp models.RoleToPermission
	for _, permissionID := range roleAndPermissions.Permissions {
		rtp = models.RoleToPermission{RoleID: role.ID, PermissionID: uint(permissionID)}
		db.NewInsert().Model(&rtp).Exec(ctx)
	}
	return c.Status(fiber.StatusOK).JSON("roles")
}

func UpdateRole(c *fiber.Ctx) error {
	var role models.Role
	var err error

	roleID, _ := c.ParamsInt("id")
	role.ID = uint(roleID)
	if role.ID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON("corrupted id")
	}
	var roleAndPermissions = struct {
		Name        string `json:"name"`
		Permissions []int  `json:"permissions"`
	}{}
	err = json.Unmarshal(c.Body(), &roleAndPermissions)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("some problem?")
		fmt.Println("err")
	}
	if len(roleAndPermissions.Name) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON("role name cannot be empty")
	}
	role.Name = roleAndPermissions.Name
	db := database.DB
	roleID, err = c.ParamsInt("id")
	if roleID <= 0 || err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("bad request")
	}

	var rtp models.RoleToPermission
	db.NewDelete().Model(&rtp).Where("role_id = ?", roleID).Exec(ctx)
	_, err = db.NewUpdate().Model(&role).Column("name").Where("id = ?", roleID).Exec(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON("update problem")
	}

	for _, permissionID := range roleAndPermissions.Permissions {
		rtp = models.RoleToPermission{RoleID: role.ID, PermissionID: uint(permissionID)}
		db.NewInsert().Model(&rtp).Exec(ctx)
	}

	return c.Status(fiber.StatusOK).JSON("roles")
}

func DeleteRole(c *fiber.Ctx) error {
	roleID, err := c.ParamsInt("id")
	if roleID == 0 || err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("bad request")
	}
	db := database.DB
	_, err = db.NewDelete().Model(&models.Role{}).Where("id = ?", roleID).Exec(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON("delete problem")
	}
	return c.Status(fiber.StatusOK).JSON("Successfully deleted.")
}
