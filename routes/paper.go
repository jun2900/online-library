package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jun2900/online-library/controllers"
	"github.com/jun2900/online-library/middleware"
)

func PaperRoutes(app *fiber.App) {
	app.Get("/papers", controllers.ReadAllPaper)
	app.Get("/paper/create", middleware.VerifyToken, controllers.CreatePaperGet)
	app.Post("/paper/create", middleware.VerifyToken, controllers.CreatePaperPost)
	app.Get("/download/paper/:id", middleware.VerifyToken, middleware.VerifyParamIdPaper, controllers.DownloadPaper)
	app.Get("paper/:id", middleware.VerifyParamIdPaper, controllers.ReadSpecificPaper)
	app.Put("/paper/:id", middleware.VerifyToken, middleware.VerifyParamIdPaper, controllers.UpdatePaperPut)
	app.Delete("/paper/:id", middleware.VerifyToken, middleware.VerifyParamIdPaper, controllers.DeletePaper)
}
