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
	users.Post("/", controllers.AddUser)
	users.Get("/", controllers.GetUsers)
	users.Get("/:id", controllers.GetUser)
	users.Get("/:id/comment", controllers.GetCommentByUser)
	users.Get("/:id/blogs", controllers.GetBlogsByUser)
	users.Put("/:id", controllers.UpdateUser)
	users.Delete("/:id", controllers.DeleteUser)

	blogs := api.Group("/blogs")
	blogs.Post("/", controllers.AddBlog)
	blogs.Get("/", controllers.GetBlogs)
	blogs.Get("/:id", controllers.GetBlog)
	blogs.Get("/:id/comment", controllers.GetCommentByBlog)
	blogs.Put("/:id", controllers.UpdateBlog)
	blogs.Delete("/:id", controllers.DeleteBlog)

	permissions := api.Group("/permissions")
	permissions.Post("/", controllers.AddPermission)
	permissions.Get("/", controllers.GetPermissions)
	permissions.Get("/:id", controllers.GetPermission)
	permissions.Put("/:id", controllers.UpdatePermission)
	permissions.Delete("/:id", controllers.DeletePermission)

	roles := api.Group("/roles")
	roles.Post("/", controllers.AddRole)
	roles.Get("/", controllers.GetRoles)
	roles.Get("/:id", controllers.GetRole)
	roles.Put("/:id", controllers.UpdateRole)
	roles.Delete("/:id", controllers.DeleteRole)

	comments := api.Group("/comments")
	comments.Post("/", controllers.AddComment)
	comments.Get("/", controllers.GetAllComments)
	comments.Get("/:id", controllers.GetComment)
	comments.Put("/:id", controllers.UpdateComment)
	comments.Delete("/:id", controllers.DeleteComment)

}
