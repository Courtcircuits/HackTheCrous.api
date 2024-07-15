package restaurants

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type IRestaurantController interface {
	GetRestaurant(ctx *fiber.Ctx) error
	GetRestaurants(ctx *fiber.Ctx) error
	Search(ctx *fiber.Ctx) error
	GetMeals(ctx *fiber.Ctx) error
}

type RestaurantControllerImpl struct {
	restaurantService IRestaurantService
}

func (controller RestaurantControllerImpl) GetRestaurant(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		log.Error(err)
		return ctx.Status(400).SendString(err.Error())
	}

	restaurant, err := controller.restaurantService.GetRestaurant(id)
	if err != nil {
		log.Error(err)
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.JSON(restaurant)
}

func (controller RestaurantControllerImpl) GetRestaurants(ctx *fiber.Ctx) error {
	restaurants, err := controller.restaurantService.GetRestaurants()
	if err != nil {
		log.Error(err)
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.JSON(restaurants)
}

func (controller RestaurantControllerImpl) Search(ctx *fiber.Ctx) error {
	query := ctx.Query("q")
	restaurants, err := controller.restaurantService.Search(query)
	if err != nil {
		log.Error(err)
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.JSON(restaurants)
}

func (controller RestaurantControllerImpl) GetMeals(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	meals, err := controller.restaurantService.GetMeals(id)
	if err != nil {
		log.Error(err)
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.JSON(meals)
}

var RestaurantController = RestaurantControllerImpl{
	restaurantService: RestaurantService,
}
