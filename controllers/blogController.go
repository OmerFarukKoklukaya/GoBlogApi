package controllers

import (
	"blog/database"
	"blog/middlewares"
	"blog/models"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

var ctx = context.Background()

func GetBlog(c *fiber.Ctx) error {
	blogId, err := c.ParamsInt("id")
	if err != nil || blogId <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": "Wrong url path"})
	}
	blog := &models.Blog{}
	db := database.DB

	db.NewSelect().Model(&blog).Where("\"blog\".\"id\" = ?", blogId).Relation("User").Scan(ctx)

	if blog.User == nil {
		blog.User = &models.User{
			Name: "this author deleted",
		}
	}

	return c.JSON(models.ViewData{
		Data: blog,
	})
}

func GetBlogs(c *fiber.Ctx) error {
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
	var blogs []models.Blog
	blogLen, _ := db.NewSelect().Model(&blogs).Relation("User").Order("id DESC").Limit(limit).Offset(offset).ScanAndCount(ctx)

	paging := models.Pagination{}.Paginate(page, blogLen, limit)

	for _, blog := range blogs {
		if blog.User == nil {
			blog.User = &models.User{Name: "this author deleted"}
		}
	}

	return c.JSON(models.ViewData{
		Data: blogs,
		Meta: models.Meta{
			Pagination: paging,
		},
	})
}

func GetBlogsByUser(c *fiber.Ctx) error {
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
	userId, err := c.ParamsInt("id")
	if err != nil || userId <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": "Wrong url path"})
	}

	blogs := []models.Blog{}
	blogLen, _ := db.NewSelect().Model(&blogs).Where("user_id = ?", userId).Order("id DESC").Limit(limit).Offset(offset).ScanAndCount(ctx)
	for i, _ := range blogs {
		if blogs[i].User == nil {
			blogs[i].User = &models.User{Name: "this author deleted"}
		}
	}

	paging := models.Pagination{}.Paginate(page, blogLen, limit)

	return c.JSON(models.ViewData{
		Data: blogs,
		Meta: models.Meta{
			Pagination: paging,
		},
	})
}

func AddBlog(c *fiber.Ctx) error {
	db := database.DB
	newBlog := models.Blog{}
	userID, _ := strconv.Atoi(c.FormValue("user_id"))
	newBlog.Title = c.FormValue("title")
	newBlog.Body = c.FormValue("body")
	newBlog.Summary = c.FormValue("summary")
	if newBlog.Title == "" || newBlog.Body == "" || newBlog.Summary == "" {
		return c.Status(fiber.StatusBadRequest).JSON("Title nor body cannot be empty")
	}

	user := models.User{}
	user, _ = middlewares.SelectAuthenticatedUser(c, db)

	newBlog.UserID = user.ID
	for _, permission := range user.Role.Permissions {
		if permission.Name == "edit_blogs" && userID > 0 {
			newBlog.UserID = uint(userID)
			break
		}
	}

	_, err := db.NewInsert().Model(&newBlog).Exec(ctx)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": newBlog})
}

func UpdateBlog(c *fiber.Ctx) error {
	db := database.DB
	var err error
	var oldBlog models.Blog
	var blog models.Blog
	id, _ := c.ParamsInt("id")
	blog.ID = uint(id)
	blog.Title = c.FormValue("title")
	blog.Body = c.FormValue("body")
	userID, _ := strconv.Atoi(c.FormValue("user_id"))
	blog.Summary = c.FormValue("summary")

	if blog.Title == "" || blog.Body == "" || blog.Summary == "" {
		return c.Status(fiber.StatusBadRequest).JSON("Title nor body cannot be empty")
	}

	db.NewSelect().Model(&oldBlog).Where("id = ?", blog.ID).Scan(ctx)
	if oldBlog.ID != blog.ID {
		return c.Status(fiber.StatusBadRequest).JSON("Bad request")
	}

	authedUser, _ := middlewares.SelectAuthenticatedUser(c, db)
	blog.UserID = authedUser.ID
	for _, permission := range authedUser.Role.Permissions {
		if permission.Name == "edit_blogs" && userID > 0 {
			blog.UserID = uint(userID)
		}
	}

	_, err = db.NewUpdate().Model(&blog).Column("title", "body", "summary", "user_id").Where("id = ?", blog.ID).Exec(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": blog})
}

func DeleteBlog(c *fiber.Ctx) error {
	db := database.DB
	var blogId, _ = c.ParamsInt("id")
	var blog models.Blog
	if blogId <= 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	db.NewSelect().Model(&blog).Where("id = ?", blogId).Scan(ctx)
	db.NewDelete().Model(&models.Blog{}).Where("id = ?", blogId).Exec(ctx)
	return c.Status(fiber.StatusOK).JSON("Blog successfully deleted")
}
