package controllers

import (
	"blog/database"
	"blog/models"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func SelectPermission(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	db := database.DB
	var permission models.Permission

	db.NewSelect().Model(&permission).Where("id = ?", id).Scan(ctx)

	return c.JSON(models.ViewData{Data: permission})
}

func SelectPermissions(c *fiber.Ctx) error {
	db := database.DB
	var permissions []models.Permission
	db.NewSelect().Model(&permissions).Scan(ctx)
	return c.JSON(models.ViewData{Data: permissions})
}

func AddPermission(c *fiber.Ctx) error {
	var permission models.Permission
	err := json.Unmarshal(c.Body(), &permission)
	if err != nil || permission.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON("Permission name cannot be empty")
	}

	db := database.DB
	var dummPermission models.Permission
	db.NewSelect().Model(&dummPermission).Where("name = ?", permission.Name).Scan(ctx)
	if dummPermission.Name != "" {
		return c.Status(fiber.StatusBadRequest).JSON("Already have this permission")
	}
	_, err = db.NewInsert().Model(&permission).Exec(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}
	return c.JSON("permissions")
}

func UpdatePermission(c *fiber.Ctx) error {
	var permission models.Permission
	err := json.Unmarshal(c.Body(), &permission)
	if err != nil {
		fmt.Println("err")
	}
	if len(permission.Name) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON("Permission name cannot be empty")
	}
	db := database.DB
	var dummPermission models.Permission
	db.NewSelect().Model(&dummPermission).Where("name = ?", permission.Name).Scan(ctx)
	if dummPermission.Name != "" {
		return c.Status(fiber.StatusBadRequest).JSON("Already have this permission")
	}
	var oldPermission models.Permission
	oldPermissionID, err := c.ParamsInt("id")
	oldPermission.ID = uint(oldPermissionID)
	if oldPermission.ID == 0 || err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("bad request")
	}

	_, err = db.NewUpdate().Model(&permission).ExcludeColumn("id").Where("id = ?", oldPermission.ID).Exec(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON("update problem")
	}
	return c.Status(fiber.StatusOK).JSON("")
}

func DeletePermission(c *fiber.Ctx) error {
	permissionID, err := c.ParamsInt("id")
	if permissionID == 0 || err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("bad request")
	}
	db := database.DB
	_, err = db.NewDelete().Model(&models.Permission{}).Where("id = ?", permissionID).Exec(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON("delete problem")
	}
	return c.Status(fiber.StatusOK).JSON("")
}
