package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jun2900/online-library/controllers"
	"github.com/jun2900/online-library/middleware"
)

func PaperRoutes(app *fiber.App) {
	app.Get("/papers", controllers.ReadAllAcceptedPaper)
	app.Get("/papers/auth", middleware.VerifyToken, middleware.IsAdmin, controllers.ReadAllPaper)
	app.Put("/paper/auth/:id/", middleware.VerifyToken, middleware.IsAdmin, middleware.VerifyParamIdPaper, controllers.UpdatePaperStatus)
	app.Get("/paper/create", middleware.VerifyToken, controllers.CreatePaperGet)
	app.Post("/paper/create", middleware.VerifyToken, controllers.CreatePaperPost)
	app.Get("/paper/download/:id", middleware.VerifyToken, middleware.VerifyParamIdPaper, controllers.DownloadPaper)
	app.Get("paper/:id", middleware.VerifyParamIdPaper, controllers.ReadSpecificPaper)
	app.Put("/paper/:id", middleware.VerifyToken, middleware.VerifyParamIdPaper, controllers.UpdatePaperPut)
	app.Delete("/paper/:id", middleware.VerifyToken, middleware.VerifyParamIdPaper, controllers.DeletePaper)
}
