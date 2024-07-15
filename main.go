package main

import (
	"github.com/Courtcircuits/HackTheCrous.api/api"
	"github.com/Courtcircuits/HackTheCrous.api/domains/restaurants"
	"github.com/Courtcircuits/HackTheCrous.api/util"
)

func main() {
	port := util.Get("PORT")
	server := api.NewServer(port)
	restaurantRouter := restaurants.RestaurantRouter{
		Controller: restaurants.RestaurantController,
	}
	server.Ignite(restaurantRouter)
}
