package controllers

import (
	"blog/database"
	"blog/models"
	"blog/utils"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"os"
	"strconv"
	"strings"
)

func SelectImage(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if id == 0 || err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.SendFile("./images/blog-"+strconv.Itoa(id)+".jpeg", true)
}

func AddImage(c *fiber.Ctx) error {
	db := database.DB

	form, err := c.MultipartForm()
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	var blog models.Blog

	db.NewSelect().Model(&blog).Where("id = ?", id).Scan(ctx)
	if blog.ID == 0 {
		return c.SendStatus(fiber.StatusNotFound)
	}

	if err != nil {
		fmt.Println("POST IMAGE: Img read problem:", err)
		return c.JSON("Img read problem")
	}

	fmt.Println("form", form)
	fmt.Println("form file", form.File)
	fmt.Println("form value", form.Value)

	if form.File["image"] == nil {
		fmt.Println("POST IMAGE: Incoming image problem: file is empty")
		return c.Status(fiber.StatusBadRequest).JSON("Incoming image file is empty")
	}
	file := form.File["image"][0]
	Header := strings.Split(file.Header["Content-Type"][0], "/")
	fileType := Header[0]
	if fileType != "image" {
		return c.SendString("unacceptable file type")
	}
	file.Filename = "blog-" + strconv.Itoa(int(blog.ID))
	err = utils.ImgScaleAndSave(file)
	if err != nil {
		fmt.Println("Img problem:", err)
		return c.Status(fiber.StatusInternalServerError).JSON("Img problem")
	}

	blog.Image = "/blogs/" + strconv.Itoa(int(blog.ID)) + "/image"
	_, err = db.NewUpdate().Model(&blog).Column("image").Where("id = ?", blog.ID).Exec(ctx)
	if err != nil {
		fmt.Println("POST IMAGE: update blog err:", err)
	}

	return c.SendStatus(200)
}

func DeleteImage(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	var blog models.Blog
	db := database.DB

	db.NewSelect().Model(&blog).Where("id = ?", id).Scan(ctx)
	if blog.ID == 0 {
		return c.SendStatus(fiber.StatusNotFound)
	}

	if id <= 0 || err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	blog.Image = ""

	_, err = db.NewUpdate().Model(&blog).Column("image").Where("id = ?", id).Exec(ctx)
	os.Remove("./images/blog-" + strconv.Itoa(int(blog.ID)) + ".jpeg")
	return c.JSON("Image Deleted")
}
