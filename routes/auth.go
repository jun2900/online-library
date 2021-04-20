package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jun2900/online-library/controllers"
	"github.com/jun2900/online-library/middleware"
)

func AuthRoutes(app *fiber.App) {
	app.Post("/signup", middleware.VerifyToken, controllers.Signup)
	app.Post("/login", controllers.Login)
}
