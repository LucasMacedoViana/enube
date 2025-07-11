package handler

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"enube/internal/infra/import"
	"github.com/gofiber/fiber/v2"
)

// ImportFile godoc
// @Summary      Importa dados a partir de planilha .xlsx
// @Tags         Importação
// @Accept       multipart/form-data
// @Produce      json
// @Param        file formData file true "Arquivo Excel (.xlsx)"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Router       /import [post]
func ImportFile(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Arquivo não enviado"})
	}

	// Criar pasta temporária
	savePath := filepath.Join(os.TempDir(), file.Filename)
	if err := c.SaveFile(file, savePath); err != nil {
		log.Println("Erro ao salvar arquivo:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao salvar arquivo"})
	}

	// Chamar função de importação
	if err := importador.ImportFromExcel(savePath); err != nil {
		log.Println("Erro ao importar:", err)
		return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("Erro na importação: %v", err)})
	}

	return c.JSON(fiber.Map{"message": "Importação concluída com sucesso"})
}
