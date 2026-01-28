package command

import (
	"log"
	"money/config"
	"money/routes"
	"net/http"
)

func Serve() {
	configuration := config.GetConfig()

	config.ConnectDatabase()
	routes.SetupRoutes()

	log.Fatal(http.ListenAndServe(":"+configuration.GolangPort, nil))
}
