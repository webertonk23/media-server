package handlers

import (
	"os"
	"path/filepath"
	"strings"

	"media-server/internal/config"

	"github.com/gofiber/fiber/v2"
)

type FileInfo struct {
	Name  string `json:"name"`
	Path  string `json:"path"`
	IsDir bool   `json:"is_dir"`
}

func ListDirectory(c *fiber.Ctx) error {
	requestedPath := c.Query("path", "/")

	basePath := config.AppConfig.MediaPath
	if basePath == "" {
		basePath = "./media"
	}

	cleanPath := filepath.Clean(requestedPath)
	fullPath := filepath.Join(basePath, cleanPath)

	absBasePath, err := filepath.Abs(basePath)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Erro interno do servidor",
		})
	}

	absFullPath, err := filepath.Abs(fullPath)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Caminho inválido",
		})
	}

	if !strings.HasPrefix(absFullPath, absBasePath) {
		return c.Status(403).JSON(fiber.Map{
			"error": "Acesso negado: caminho fora do diretório permitido",
		})
	}

	stat, err := os.Stat(absFullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return c.Status(404).JSON(fiber.Map{
				"error": "Diretório não encontrado",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"error": "Erro ao acessar diretório",
		})
	}

	if !stat.IsDir() {
		return c.Status(400).JSON(fiber.Map{
			"error": "Caminho não é um diretório",
		})
	}

	files, err := os.ReadDir(absFullPath)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Não foi possível ler o diretório",
		})
	}

	var result []FileInfo
	for _, file := range files {
		if file.IsDir() {
			relPath, _ := filepath.Rel(absBasePath, filepath.Join(absFullPath, file.Name()))
			result = append(result, FileInfo{
				Name:  file.Name(),
				Path:  relPath,
				IsDir: true,
			})
		}
	}

	return c.JSON(result)
}
