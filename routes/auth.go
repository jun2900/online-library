package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jun2900/online-library/controllers"
)

func AuthRoutes(app *fiber.App) {
	app.Post("/signup", controllers.Signup)
}
