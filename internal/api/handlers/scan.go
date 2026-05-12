package handlers
import (
	"github.com/gofiber/fiber/v2"
)
func ScanLibrary(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Scanner is running in background",
	})
}
