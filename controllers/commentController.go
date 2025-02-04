package controllers

import (
	"blog/database"
	"blog/middlewares"
	"blog/models"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/bun"
	"strconv"
)

func SelectComment(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil && id <= 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	db := database.DB
	var comment models.Comment
	err = db.NewSelect().Model(&comment).Where("\"comment\".\"id\" = ?", id).
		Relation("Blog").
		Relation("User", bunPasswordExclude).
		Relation("ParentComment", func(query *bun.SelectQuery) *bun.SelectQuery {
			return query.Relation("User", bunPasswordExclude)
		}).
		Relation("ChildComments", func(query *bun.SelectQuery) *bun.SelectQuery {
			return query.Relation("User", bunPasswordExclude)
		}).Scan(ctx)
	if err != nil {
		fmt.Println("COMEMNT ERROR: ", err)
	}

	return c.JSON(models.ViewData{Data: comment})
}

func SelectAllComments(c *fiber.Ctx) error {
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
	var comments []models.Comment
	commentLen, err := db.NewSelect().Model(&comments).Relation("Blog").Relation("User", bunPasswordExclude).Relation("ParentComment").Relation("ChildComments").Order("id DESC").Limit(limit).Offset(offset).ScanAndCount(ctx)
	if err != nil {
		fmt.Println(err)
	}

	paging := models.Pagination{}.Paginate(page, commentLen, limit)

	return c.JSON(models.ViewData{Data: comments, Meta: models.Meta{paging}})
}

func SelectCommentByBlog(c *fiber.Ctx) error {
	blogId, err := c.ParamsInt("id")
	if err != nil && blogId <= 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	page := c.QueryInt("page", 1)
	if page <= 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	limit := c.QueryInt("limit", 5)
	if limit <= 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	offset := limit * (page - 1)

	db := database.DB
	var comments []models.Comment
	commentLen, err := db.NewSelect().
		Model(&comments).
		Relation("User", bunPasswordExclude).Relation("ChildComments", func(query *bun.SelectQuery) *bun.SelectQuery {
		return query.Relation("User", bunPasswordExclude)
	}).
		WhereGroup("AND", func(query *bun.SelectQuery) *bun.SelectQuery {
			return query.Where("blog_id = ?", blogId).Where("parent_comment_id = ?", 0)
		}).
		Order("id DESC").
		Offset(offset).Limit(limit).
		ScanAndCount(ctx)

	if err != nil {
		fmt.Println("Select Comment By Blog error: ", err)
	}

	paging := models.Pagination{}.Paginate(page, commentLen, limit)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	fmt.Println("BLOG ID:", blogId, "\n", "OFFSET:", offset, "\n", "ERR:", err, "\n", "COMMENT LEN:", commentLen, "\n")
	return c.JSON(models.ViewData{
		Data: comments,
		Meta: models.Meta{
			Pagination: paging,
		},
	})
}

func bunPasswordExclude(query *bun.SelectQuery) *bun.SelectQuery {
	return query.ExcludeColumn("password")
}

func SelectCommentByUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil && id <= 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	page := c.QueryInt("page", 1)
	if page <= 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	limit := c.QueryInt("limit", 5)
	if limit <= 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	offset := limit * (page - 1)

	db := database.DB
	var comment models.Comment
	commentLen, _ := db.NewSelect().Model(&comment).Where("user_id = ?", id).Limit(limit).Offset(offset).ScanAndCount(ctx)

	paginate := models.Pagination{}.Paginate(page, commentLen, limit)
	return c.JSON(models.ViewData{
		Data: comment,
		Meta: models.Meta{paginate},
	})
}

func AddComment(c *fiber.Ctx) error {
	content := c.FormValue("content")
	blogId, err := strconv.Atoi(c.FormValue("blog_id"))

	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if content == "" || blogId <= 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	parentCommentId, _ := strconv.Atoi(c.FormValue("parent_comment_id"))

	db := database.DB
	var comment models.Comment
	user, err := middlewares.SelectAuthenticatedUser(c, db)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	comment.UserID = user.ID
	comment.BlogID = uint(blogId)
	comment.ParentCommentID = uint(parentCommentId)
	comment.Content = content
	comment.IsParent = false
	_, err = db.NewInsert().Model(&comment).Exec(ctx)
	comment.User = &user
	comment.User.Password = nil
	comment.User.Role = &models.Role{}

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if parentCommentId > 0 {
		_, err = db.NewUpdate().Model(&models.Comment{IsParent: true}).Column("is_parent").Where("id = ?", parentCommentId).Exec(ctx)
		if err != nil {
			fmt.Println(err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

	}

	return c.JSON(models.ViewData{Data: comment})
}

func UpdateComment(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil && id <= 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	content := c.FormValue("content")
	if content == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	db := database.DB
	var comment models.Comment
	db.NewSelect().Model(&comment).Where("id = ?", id).Scan(ctx)
	if comment.ID == 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	user, err := middlewares.SelectAuthenticatedUser(c, db)
	if user.ID != comment.UserID {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	_, err = db.NewUpdate().Model(&comment).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusOK).JSON(models.ViewData{Data: comment})
}

func DeleteComment(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil && id <= 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	db := database.DB
	_, err = db.NewDelete().Model(&models.Comment{}).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusOK)
}
