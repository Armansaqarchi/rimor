package main

import (
	"log"
	"rimor/pkg/web/routes"
)

func main() {

	router := routes.NewRouter()
	if err := router.Router.Listen(router.Port); err != nil {
		log.Fatal(err.Error())
	}


}



