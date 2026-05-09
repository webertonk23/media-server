package handlers

import (
	"os"
	"path/filepath"
	"github.com/gofiber/fiber/v2"
)

type FileInfo struct {
	Name  string `json:"name"`
	Path  string `json:"path"`
	IsDir bool   `json:"is_dir"`
}

func ListDirectory(c *fiber.Ctx) error {
	path := c.Query("path", "/")
	
	files, err := os.ReadDir(path)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Não foi possível ler o diretório"})
	}

	var result []FileInfo
	for _, file := range files {
		if file.IsDir() {
			result = append(result, FileInfo{
				Name:  file.Name(),
				Path:  filepath.Join(path, file.Name()),
				IsDir: true,
			})
		}
	}

	return c.JSON(result)
}
