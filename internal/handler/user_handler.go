package handler

import (
	"enube/internal/domain/dto"
	"enube/internal/domain/model"
	"enube/internal/infra/db"
	"enube/internal/utils"
	"github.com/go-playground/validator/v10"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// CreateUser godoc
// @Summary      Cadastra um novo usuário
// @Description  Cria um novo usuário com nome e senha
// @Tags         Usuários
// @Accept       json
// @Produce      json
// @Param        user body dto.UserInputDTO true "Dados do usuário"
// @Success      201 {object} dto.UserOutputDTO
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /api/users [post]
// @Security     BearerAuth
func CreateUser(c *fiber.Ctx) error {
	var input dto.UserInputDTO
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "JSON inválido")
	}

	validate := validator.New()

	if err := validate.Struct(input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Campos obrigatórios inválidos")
	}

	// Verifica duplicidade
	var exists model.User
	if err := db.DB.Where("name = ?", input.Name).First(&exists).Error; err == nil {
		return fiber.NewError(fiber.StatusBadRequest, "Usuário já existe")
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Erro ao criptografar senha")
	}

	user := model.User{
		Name:     input.Name,
		Password: hashedPassword,
	}

	if err := db.DB.Create(&user).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Erro ao salvar usuário")
	}

	output := dto.UserOutputDTO{
		Name: user.Name,
	}

	return c.Status(fiber.StatusCreated).JSON(output)
}

// GetAllUsers godoc
// @Summary      Lista todos os usuários
// @Description  Retorna todos os usuários cadastrados
// @Tags         Usuários
// @Produce      json
// @Success      200 {array} dto.UserOutputDTO
// @Failure      500 {object} map[string]string
// @Router       /api/users [get]
// @Security     BearerAuth
func GetAllUsers(c *fiber.Ctx) error {
	var users []model.User
	if err := db.DB.Find(&users).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Erro ao buscar usuários")
	}

	var output []dto.UserOutputDTO
	for _, u := range users {
		output = append(output, dto.UserOutputDTO{
			Name: u.Name,
		})
	}

	return c.JSON(output)
}

// GetUserByID godoc
// @Summary      Busca usuário por ID
// @Description  Retorna um único usuário pelo ID
// @Tags         Usuários
// @Produce      json
// @Param        id path int true "ID do usuário"
// @Success      200 {object} dto.UserOutputDTO
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /api/users/{id} [get]
// @Security     BearerAuth
func GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")

	var user model.User
	if err := db.DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fiber.NewError(fiber.StatusNotFound, "Usuário não encontrado")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Erro ao buscar usuário")
	}

	output := dto.UserOutputDTO{
		Name: user.Name,
	}

	return c.JSON(output)
}
