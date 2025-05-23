package api

import (
	"PerpusGo/domain"
	"PerpusGo/dto"
	"PerpusGo/internal/util"
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

type CustomerAPI struct {
	customerService domain.CustomerService
}

func NewCustomer(app *fiber.App, cutomerService domain.CustomerService, auzMidd fiber.Handler) {
	ca := CustomerAPI{
		customerService: cutomerService,
	}

	app.Get("/customers", auzMidd, ca.Index)
	app.Post("/customers", auzMidd, ca.Create)
	app.Put("/customers/:id", auzMidd, ca.Update)
	app.Delete("/customers/:id", auzMidd, ca.Delete)
	app.Get("/customers/:id", auzMidd, ca.Show)
}

func (ca *CustomerAPI) Index(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	res, err := ca.customerService.Index(c)

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).
			JSON(dto.CreateResponseError(err.Error()))
	}

	return ctx.JSON(dto.CreateResponseSuccess(res))
}

func (ca *CustomerAPI) Create(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var req dto.CreateCustomerRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.SendStatus(http.StatusUnprocessableEntity)
	}

	fails := util.Validate(req)

	if len(fails) > 0 {
		return ctx.Status(http.StatusBadRequest).
			JSON(dto.CreateResponseErrorData("validation error", fails))
	}

	err := ca.customerService.Create(c, req)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).
			JSON(dto.CreateResponseError(err.Error()))
	}

	return ctx.Status(http.StatusCreated).
		JSON(dto.CreateResponseSuccess("customer created"))
}

func (ca *CustomerAPI) Update(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var req dto.UpdateCustomerRequest

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.SendStatus(http.StatusUnprocessableEntity)
	}

	fails := util.Validate(req)
	if len(fails) > 0 {
		return ctx.Status(http.StatusBadRequest).
			JSON(dto.CreateResponseErrorData("validation error", fails))
	}

	req.ID = ctx.Params("id")
	err := ca.customerService.Update(c, req)

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).
			JSON(dto.CreateResponseError(err.Error()))
	}

	return ctx.Status(http.StatusOK).
		JSON(dto.CreateResponseSuccess(""))

}

func (ca *CustomerAPI) Delete(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	id := ctx.Params("id")
	err := ca.customerService.Delete(c, id)

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).
			JSON(dto.CreateResponseError(err.Error()))
	}

	return ctx.SendStatus(http.StatusNoContent)
}

func (ca *CustomerAPI) Show(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	id := ctx.Params("id")
	data, err := ca.customerService.Show(c, id)

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).
			JSON(dto.CreateResponseError(err.Error()))
	}

	return ctx.Status(http.StatusOK).
		JSON(dto.CreateResponseSuccess(data))
}
