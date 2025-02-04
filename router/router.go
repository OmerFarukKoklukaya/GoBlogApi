package router

import (
	"blog/controllers"
	"blog/middlewares"
	"github.com/gofiber/fiber/v2"
)

func Router(app *fiber.App) {
	api := app.Group("/api")

	api.Post("/register", controllers.Register)
	api.Post("/login", controllers.Login)

	users := api.Group("/users")
	blogs := api.Group("/blogs")
	permissions := api.Group("/permissions")
	roles := api.Group("/roles")
	comments := api.Group("/comments")

	users.Get("/", controllers.SelectUsers)
	users.Get("/:id", controllers.SelectUser)
	users.Get("/:id/comments", controllers.SelectCommentByUser)

	blogs.Get("/", controllers.SelectBlogs)
	blogs.Get("/:id", controllers.SelectBlog)
	blogs.Get("/:id/image", controllers.SelectImage)
	blogs.Get("/:id/comments", controllers.SelectCommentByBlog)

	api.Use(middlewares.AuthenticationMiddleware)

	api.Post("/logout", controllers.Logout)
	api.Get("/profile", controllers.SelectAuthedUser)

	api.Use(middlewares.AuthorizationMiddleware)

	users.Put("/password", controllers.UpdatePassword)
	users.Get("/:id/blogs", controllers.SelectBlogsByUser)
	users.Post("/", controllers.AddUser)
	users.Put("/:id", controllers.UpdateUser)
	users.Delete("/:id", controllers.DeleteUser)

	blogs.Post("/:id/image", controllers.AddImage)
	blogs.Post("/", controllers.AddBlog)
	blogs.Delete("/:id/image", controllers.DeleteImage)
	blogs.Put("/:id", controllers.UpdateBlog)
	blogs.Delete("/:id", controllers.DeleteBlog)

	permissions.Get("/", controllers.SelectPermissions)
	permissions.Get("/:id", controllers.SelectPermission)
	permissions.Post("/", controllers.AddPermission)
	permissions.Put("/:id", controllers.UpdatePermission)
	permissions.Delete("/:id", controllers.DeletePermission)

	roles.Get("/", controllers.SelectRoles)
	roles.Get("/:id", controllers.SelectRole)
	roles.Post("/", controllers.AddRole)
	roles.Put("/:id", controllers.UpdateRole)
	roles.Delete("/:id", controllers.DeleteRole)

	comments.Get("/", controllers.SelectAllComments)
	comments.Get("/:id", controllers.SelectComment)
	comments.Post("/", controllers.AddComment)
	comments.Put("/:id", controllers.UpdateComment)
	comments.Delete("/:id", controllers.DeleteComment)
}
