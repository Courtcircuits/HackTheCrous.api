package restaurants

import "github.com/gofiber/fiber/v2"

type RestaurantRouter struct {
	Controller IRestaurantController
}

func (r RestaurantRouter) SubscribeRoutes(app *fiber.Router) {
	(*app).Get("/restaurants/:id", r.Controller.GetRestaurant)
	(*app).Get("/restaurants", r.Controller.GetRestaurants)
	(*app).Get("/search", r.Controller.Search)
	(*app).Get("/restaurants/meals/:id", r.Controller.GetMeals)
}
