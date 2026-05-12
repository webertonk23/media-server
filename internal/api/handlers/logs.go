package handlers
import (
	"media-server/internal/utils"
	"github.com/gofiber/fiber/v2"
)
func GetLogs(c *fiber.Ctx) error {
	if utils.GlobalLogger == nil {
		return c.Status(500).JSON(fiber.Map{"error": "Logger não inicializado"})
	}
	return c.JSON(utils.GlobalLogger.GetLogs())
}
