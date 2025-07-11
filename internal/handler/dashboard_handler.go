package handler

import (
	"enube/internal/domain/dto"
	"enube/internal/domain/model"
	"enube/internal/infra/db"
	"github.com/gofiber/fiber/v2"
)

// GetSummary godoc
// @Summary      Totais gerais de faturamento e quantidade
// @Description  Retorna valores agregados do sistema como total faturado e quantidade de registros
// @Tags         Dashboard
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]interface{}
// @Failure      500 {object} map[string]string
// @Router       /api/dashboard/summary [get]
func GetSummary(c *fiber.Ctx) error {
	var total float64
	var count int64

	err := db.DB.Model(&model.BillingItem{}).
		Select("ROUND(sum(billing_pre_tax_total))").
		Scan(&total).Error
	if err != nil {
		return fiber.ErrInternalServerError
	}

	db.DB.Model(&model.BillingItem{}).Count(&count)

	return c.JSON(fiber.Map{
		"total_faturado": total,
		"quantidade":     count,
	})
}

// GetByMonth godoc
// @Summary      Faturamento mensal
// @Description  Retorna o total de faturamento agrupado por mês de cobrança
// @Tags         Dashboard
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} dto.DashboardEntryDTO
// @Failure      500 {object} map[string]string
// @Router       /api/dashboard/monthly [get]
func GetByMonth(c *fiber.Ctx) error {
	var results []dto.DashboardEntryDTO
	err := db.DB.Model(&model.BillingItem{}).
		Select("to_char(charge_start_date, 'YYYY-MM') AS label, ROUND(sum(billing_pre_tax_total)) AS total").
		Group("label").
		Order("label").
		Scan(&results).Error
	if err != nil {
		return fiber.ErrInternalServerError
	}
	return c.JSON(results)
}

// GetByClient godoc
// @Summary      Faturamento por cliente
// @Description  Lista o total de faturamento agrupado por cliente
// @Tags         Dashboard
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} dto.DashboardEntryDTO
// @Failure      500 {object} map[string]string
// @Router       /api/dashboard/by-client [get]
func GetByClient(c *fiber.Ctx) error {
	var results []dto.DashboardEntryDTO
	err := db.DB.Model(&model.BillingItem{}).
		Select("customers.name AS label, ROUND(sum(billing_items.billing_pre_tax_total)) AS total").
		Joins("JOIN customers ON customers.id = billing_items.customer_id").
		Group("customers.name").
		Order("total DESC").
		Scan(&results).Error
	if err != nil {
		return fiber.ErrInternalServerError
	}
	return c.JSON(results)
}

// GetByResource godoc
// @Summary      Faturamento por recurso
// @Description  Lista o faturamento agrupado por recurso (campo consumed_service)
// @Tags         Dashboard
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} dto.DashboardEntryDTO
// @Failure      500 {object} map[string]string
// @Router       /api/dashboard/by-resource [get]
func GetByResource(c *fiber.Ctx) error {
	var results []dto.DashboardEntryDTO
	err := db.DB.Model(&model.BillingItem{}).
		Select("consumed_service AS label, ROUND(sum(billing_pre_tax_total)) AS total").
		Group("consumed_service").
		Order("total DESC").
		Scan(&results).Error
	if err != nil {
		return fiber.ErrInternalServerError
	}
	return c.JSON(results)
}

// GetByCategory godoc
// @Summary      Faturamento por categoria de medição
// @Description  Retorna o total faturado agrupado por categoria do medidor (meter.category)
// @Tags         Dashboard
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} dto.DashboardEntryDTO
// @Failure      500 {object} map[string]string
// @Router       /api/dashboard/by-category [get]
func GetByCategory(c *fiber.Ctx) error {
	var results []dto.DashboardEntryDTO
	err := db.DB.Model(&model.BillingItem{}).
		Joins("JOIN meters ON meters.id = billing_items.meter_id").
		Select("meters.category AS label, ROUND(sum(billing_pre_tax_total)) AS total").
		Group("meters.category").
		Order("total DESC").
		Scan(&results).Error
	if err != nil {
		return fiber.ErrInternalServerError
	}
	return c.JSON(results)
}
