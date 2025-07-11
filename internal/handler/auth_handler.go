package handler

import (
	"enube/internal/domain/dto"
	"enube/internal/domain/model"
	"enube/internal/infra/auth"
	"enube/internal/infra/db"
	"enube/internal/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// Login godoc
// @Summary      Autenticação do usuário
// @Description  Realiza o login e retorna um token JWT
// @Tags         Autenticação
// @Accept       json
// @Produce      json
// @Param        login body dto.LoginRequestDTO true "Credenciais de login"
// @Success      200 {object} dto.LoginResponseDTO
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Router       /api/login [post]
func Login(c *fiber.Ctx) error {
	var req dto.LoginRequestDTO
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Erro no JSON")
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Campos inválidos")
	}

	var user model.User
	if err := db.DB.Where("name = ?", req.Name).First(&user).Error; err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Usuário não encontrado")
	}

	if err := utils.CheckPasswordHash(req.Password, user.Password); err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Senha incorreta")
	}

	token, err := auth.GenerateJWT(user.Name)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Erro ao gerar token")
	}

	return c.JSON(dto.LoginResponseDTO{Token: token})
}
