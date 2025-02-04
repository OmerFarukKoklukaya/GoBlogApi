package router

import (
	"blog/controllers"
	"github.com/gofiber/fiber/v2"
)

func Router(app *fiber.App) {
	api := app.Group("/api")

	api.Post("/register", controllers.Register)
	api.Post("/login", controllers.Login)
	api.Post("/logout", controllers.Logout)

	users := api.Group("/users")
	users.Get("/", controllers.SelectUsers)
	users.Get("/:id", controllers.SelectUser)
	users.Post("/", controllers.AddUser)
	users.Get("/:id/comments", controllers.SelectCommentByUser)
	users.Get("/:id/blogs", controllers.SelectBlogsByUser)
	users.Put("/password", controllers.UpdatePassword)
	users.Put("/:id", controllers.UpdateUser)
	users.Delete("/:id", controllers.DeleteUser)

	blogs := api.Group("/blogs")
	blogs.Get("/", controllers.SelectBlogs)
	blogs.Get("/:id", controllers.SelectBlog)
	blogs.Get("/:id/image", controllers.SelectImage)
	blogs.Get("/:id/comments", controllers.SelectCommentByBlog)
	blogs.Post("/:id/image", controllers.AddImage)
	blogs.Post("/", controllers.AddBlog)
	blogs.Delete("/:id/image", controllers.DeleteImage)
	blogs.Put("/:id", controllers.UpdateBlog)
	blogs.Delete("/:id", controllers.DeleteBlog)

	permissions := api.Group("/permissions")
	permissions.Get("/", controllers.SelectPermissions)
	permissions.Get("/:id", controllers.SelectPermission)
	permissions.Post("/", controllers.AddPermission)
	permissions.Put("/:id", controllers.UpdatePermission)
	permissions.Delete("/:id", controllers.DeletePermission)

	roles := api.Group("/roles")
	roles.Get("/", controllers.SelectRoles)
	roles.Get("/:id", controllers.SelectRole)
	roles.Post("/", controllers.AddRole)
	roles.Put("/:id", controllers.UpdateRole)
	roles.Delete("/:id", controllers.DeleteRole)

	comments := api.Group("/comments")
	comments.Get("/", controllers.SelectAllComments)
	comments.Get("/:id", controllers.SelectComment)
	comments.Post("/", controllers.AddComment)
	comments.Put("/:id", controllers.UpdateComment)
	comments.Delete("/:id", controllers.DeleteComment)

	api.Get("/profile", controllers.SelectAuthedUser)
}
