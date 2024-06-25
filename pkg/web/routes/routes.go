package routes

import (
	"rimor/pkg/web/api"
	"github.com/gofiber/fiber/v2"
)




type Router struct { 
	Port string `env:"port"`
	Router *fiber.App
}


func NewRouter() *Router{
	r := Router{
		Port: ":8080",
		Router: fiber.New(),
	}
	r.initRoutes()
	return &r
}



func (r *Router) initRoutes() {
	router := r.Router
	e_handler := api.GetEngineHandler()

	router.Get("/query", e_handler.Query)
}