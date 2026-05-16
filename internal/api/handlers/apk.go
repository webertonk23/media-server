package handlers

import (
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

func DownloadAPK(c *fiber.Ctx) error {
	apkDir := "./apk"
	apkFileName := "media-server.apk"
	apkPath := filepath.Join(apkDir, apkFileName)

	if _, err := os.Stat(apkPath); os.IsNotExist(err) {
		return c.Status(404).JSON(fiber.Map{
			"error": "APK não encontrado",
		})
	}

	c.Set("Content-Type", "application/vnd.android.package-archive")
	c.Set("Content-Disposition", "attachment; filename="+apkFileName)

	return c.SendFile(apkPath)
}

func GetAPKInfo(c *fiber.Ctx) error {
	apkDir := "./apk"
	apkFileName := "media-server.apk"
	apkPath := filepath.Join(apkDir, apkFileName)

	fileInfo, err := os.Stat(apkPath)
	if os.IsNotExist(err) {
		return c.Status(404).JSON(fiber.Map{
			"error":     "APK não encontrado",
			"available": false,
		})
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Erro ao obter informações do APK",
		})
	}

	return c.JSON(fiber.Map{
		"available": true,
		"filename":  apkFileName,
		"size":      fileInfo.Size(),
		"modified":  fileInfo.ModTime(),
	})
}
