package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jun2900/online-library/controllers"
	"github.com/jun2900/online-library/middleware"
)

func AuthorRoutes(app *fiber.App) {
	app.Post("/author/create", middleware.VerifyToken, controllers.CreateAuthorPost)
}
