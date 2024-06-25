package routes

import (
	"rimor/pkg/web/api"

	"github.com/gofiber/fiber/v2"
)




type Router struct { 
	Port string `env:"port"`
	Router fiber.Router
}


func NewRouter() *Router{
	r := Router{}
	r.initRoutes()


	return &r
}



func (r *Router) initRoutes() {
	router := r.Router
	e_handler := api.GetEngineHandler()

	router.Get("/query", e_handler.Query)
}