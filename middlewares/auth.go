package middlewares

import (
	"blog/database"
	"blog/models"
	"blog/utils"
	"context"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/bun"
	"strconv"
	"strings"
)

func AuthenticationMiddleware(c *fiber.Ctx) error {
	token := c.Cookies("token")
	tokenClaims := utils.ControlToken(token)
	if tokenClaims == nil {
		return fiber.ErrUnauthorized
	}
	return c.Next()
}

func AuthorizationMiddleware(c *fiber.Ctx) error {
	db := database.DB
	var authedUser, err = SelectAuthenticatedUser(c, db)
	if err != nil {
		fmt.Println(err)
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	var authedRole models.Role
	db.NewSelect().Model(&authedRole).Where("\"role\".\"id\" = ?", authedUser.RoleID).Relation("Permissions").Scan(context.Background())

	groupName, _ := strings.CutPrefix(c.Path(), "/api")

	if groupName == "/" {
		return c.Next()
	}
	groupName, _, _ = strings.Cut(groupName[1:], "/")
	fmt.Println(groupName)
	if groupName[len(groupName)-1:] != "s" {
		groupName = groupName + "s"
	}
	fmt.Println(groupName)

	var modelMap = make(map[string]any)
	parameter := strings.Split(c.Path(), "/")
	targetID, err := strconv.Atoi(parameter[len(parameter)-1])
	if err == nil {
		err = db.NewSelect().Model(&modelMap).Table(groupName).Where("id = ?", targetID).Scan(context.Background())
	}

	if (groupName == "blogs" && c.Route().Method == "POST") || (groupName == "blogs" && modelMap["user_id"] == authedUser.ID) {
		return c.Next()
	} else if groupName == "users" && modelMap["id"] == authedUser.ID {
		return c.Next()
	}

	for _, permission := range authedRole.Permissions {
		if c.Route().Method == "POST" || c.Route().Method == "PUT" || c.Route().Method == "DELETE" {
			if permission.Name == "edit_"+groupName {
				return c.Next()
			}
		} else if c.Route().Method == "GET" {
			if permission.Name == "read_"+groupName {
				return c.Next()
			}
		}
	}

	return c.SendStatus(fiber.StatusUnauthorized)

}

func SelectAuthenticatedUser(c *fiber.Ctx, db *bun.DB) (models.User, error) {
	rawToken := c.Cookies("token")
	if len(rawToken) == 0 {
		return models.User{}, errors.New("not logged in")
	}
	tokenClaims := utils.ControlToken(rawToken)
	if tokenClaims == nil {
		return models.User{}, errors.New("invalid token")
	}
	var user models.User

	err := db.NewSelect().Model(&user).Where("\"user\".\"id\" = ?", tokenClaims.ID).Relation("Role").Scan(context.Background())
	if err != nil {
		fmt.Println("Authenticated", err)
		return models.User{}, err
	}
	var role models.Role
	db.NewSelect().Model(&role).Where("id = ?", user.RoleID).Relation("Permissions").Scan(context.Background())
	user.Role = &role
	return user, nil
}
