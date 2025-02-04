package controllers

import (
	"blog/database"
	"blog/middlewares"
	"blog/models"
	"blog/utils"
	"bytes"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func SelectAuthedUser(c *fiber.Ctx) error {
	db := database.DB
	user, err := middlewares.SelectAuthenticatedUser(c, db)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON("You are not logged in")
	}
	return c.Status(fiber.StatusOK).JSON(user)
}

func Register(c *fiber.Ctx) error {
	db := database.DB
	userName := c.FormValue("name")
	if len(userName) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON("Name cannot be empty")
	}
	password := []byte(c.FormValue("password"))
	if len(password) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON("Password cannot be empty")
	}
	passwordVer := []byte(c.FormValue("password_verification"))
	if len(passwordVer) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON("Password verification cannot be empty")
	}
	if !bytes.Equal(password, passwordVer) {
		return c.Status(fiber.StatusBadRequest).JSON("Password verification failed")
	}

	var user models.User
	db.NewSelect().Model(&user).Where("name = ?", userName).Relation("Blogs").Scan(ctx)
	if len(user.Name) != 0 {
		return c.Status(fiber.StatusBadRequest).JSON("This name already taken")
	}

	var role models.Role
	db.NewSelect().Model(&role).Where("name = ?", "basic").Scan(ctx)

	user = models.User{Name: userName, RoleID: role.ID}
	user.ChangePassword(password)
	db.NewInsert().Model(&user).Exec(ctx)
	return c.SendStatus(fiber.StatusOK)
}

func Login(c *fiber.Ctx) error {
	db := database.DB
	user := models.User{}
	userName := c.FormValue("name")
	if len(userName) == 0 {
		fmt.Println("Login problem: name cannot be empty")
		return c.Status(fiber.StatusBadRequest).JSON("Name cannot be empty")
	}
	password := []byte(c.FormValue("password"))
	if len(password) == 0 {
		fmt.Println("Login page problem: password cannot be empty")
		return c.Status(fiber.StatusBadRequest).JSON("Password cannot be empty")
	}

	err := db.NewSelect().Model(&user).Where("name = ?", userName).Scan(ctx)
	if !user.CheckPassword(password) {
		fmt.Println("Login err:", err)
		return c.Status(fiber.StatusBadRequest).JSON("Wrong Name or Password")
	}

	c.Cookie(&fiber.Cookie{
		Name:  "token",
		Value: utils.GenerateToken(user.ID),
	})

	return c.SendStatus(fiber.StatusOK)

}

func Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:  "token",
		Value: "",
	})
	c.Method("GET")
	return c.SendStatus(fiber.StatusOK)
}

func SelectUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil && id <= 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	db := database.DB
	var user models.User

	db.NewSelect().Model(&user).Where("\"user\".\"id\" = ?", id).ExcludeColumn("password").Relation("Role").Scan(ctx)

	return c.JSON(models.ViewData{Data: user})
}

func SelectUsers(c *fiber.Ctx) error {
	limit, err := strconv.Atoi(c.Query("limit", "5"))
	if err != nil || limit <= 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page <= 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	db := database.DB

	offset := (page - 1) * limit

	users := []models.User{}
	userLen, _ := db.NewSelect().Model(&users).ExcludeColumn("password").Relation("Role").Limit(limit).Offset(offset).ScanAndCount(ctx)

	paging := models.Pagination{}.Paginate(page, userLen, limit)

	return c.JSON(models.ViewData{
		Data: users,
		Meta: models.Meta{
			Pagination: paging,
		},
	})
}

func AddUser(c *fiber.Ctx) error {
	db := database.DB
	var user models.User
	user.Name = c.FormValue("name")
	password := []byte(c.FormValue("password"))
	roleID, _ := strconv.Atoi(c.FormValue("role_id"))
	user.RoleID = uint(roleID)
	var dummUser models.User
	db.NewSelect().Model(&dummUser).Where("name = ?", user.Name).Scan(ctx)
	if len(dummUser.Name) != 0 {
		return c.Status(fiber.StatusBadRequest).JSON("This name already taken")
	}
	if user.Name == "" || len(password) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON("Name or Password cannot be empty")
	}
	user.ChangePassword(password)
	if user.RoleID <= 0 {
		user.RoleID = 2
	}

	authedUser, err := middlewares.SelectAuthenticatedUser(c, db)
	if authedUser.ID <= 0 || err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	db.NewInsert().Model(&user).Exec(ctx)

	return c.SendStatus(200)
}

func UpdateUser(c *fiber.Ctx) error {
	db := database.DB
	var user, oldUser models.User
	userID, _ := c.ParamsInt("id")
	user.ID = uint(userID)

	user.Name = c.FormValue("name")
	roleID, _ := strconv.Atoi(c.FormValue("role_id"))
	user.RoleID = uint(roleID)

	var authedUser models.User
	authedUser, err := middlewares.SelectAuthenticatedUser(c, db)
	if authedUser.ID <= 0 || err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
		fmt.Println(err)
	}

	if user.ID <= 0 {
		user.ID = authedUser.ID
		oldUser = authedUser
	} else {
		db.NewSelect().Model(&oldUser).Where("id = ?", user.ID).Scan(ctx)
		fmt.Println("old user: ", oldUser)
	}
	if user.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON("name cannot be empty")
	}

	if user.RoleID <= 0 {
		user.RoleID = oldUser.RoleID
	}

	var dummUser models.User
	db.NewSelect().Model(&dummUser).Where("name = ?", user.Name).Scan(ctx)
	fmt.Println("user", user.Name, "old user", oldUser.Name, "dummy user", dummUser.Name)
	if len(dummUser.Name) != 0 && user.Name != oldUser.Name {
		return c.Status(fiber.StatusBadRequest).JSON("This name already taken")
	}

	db.NewUpdate().Model(&user).Column("name", "role_id").Where("id = ?", user.ID).Exec(ctx)

	return c.Status(200).JSON(fiber.Map{"user": user})
}

func DeleteUser(c *fiber.Ctx) error {
	db := database.DB
	id, _ := c.ParamsInt("id")
	if id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON("unacceptable id")
	}

	db.NewDelete().Model(&models.User{}).Where("id = ?", id).Exec(ctx)
	return c.SendStatus(fiber.StatusOK)
}

func UpdatePassword(c *fiber.Ctx) error {
	oldPassword := []byte(c.FormValue("password"))
	if len(oldPassword) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON("Password Cannot Be Empty")
	}
	newPassword := []byte(c.FormValue("new_password"))
	if len(newPassword) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON("New Password Cannot Be Empty")
	}
	passVerification := []byte(c.FormValue("new_password_verification"))
	if len(passVerification) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON("New Password verification Cannot Be Empty")
	}
	if !bytes.Equal(newPassword, passVerification) {
		return c.Status(fiber.StatusBadRequest).JSON("Passwords doesnt match")
	}
	db := database.DB
	var user models.User
	user, err := middlewares.SelectAuthenticatedUser(c, db)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON("Unauthorized")
	}
	if !user.CheckPassword(oldPassword) {
		return c.Status(fiber.StatusBadRequest).JSON("Wrong password")
	}

	user.ChangePassword(newPassword)

	db.NewUpdate().Model(&user).Where("id = ?", user.ID).Exec(ctx)

	return c.SendStatus(fiber.StatusOK)
}
