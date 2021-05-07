package middleware

import "github.com/gofiber/fiber/v2"

func VerifyParamIdPaper(c *fiber.Ctx) error {
	paperId, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Error processing param id", "err": err})
	}
	c.Locals("paperId", paperId)
	return c.Next()
}
