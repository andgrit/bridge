package main

import (
	"fmt"
	"net/http"
	"log"
	"github.com/andgrit/bridge/controllers"
	"github.com/andgrit/bridge/configuration"
	"github.com/andgrit/bridge/storage"
)

func main() {
	// get the application configuration
	appConfiguration, err := configuration.DefaultConfiguration()
	if err != nil {
		panic(err)
	}
	err = storage.SetConfiguration(appConfiguration)
	if err != nil {
		panic(err)
	}
	err = controllers.SetConfiguration(appConfiguration)

	const portString = ":8080"
	const apiPrefix = "/api/v1"

	fmt.Printf("listening on port %s, maybe 'curl http://localhost%s/api/v1/version' will return stuff", portString, portString)
	log.Fatal(http.ListenAndServe(portString, controllers.MuxRouter(apiPrefix)))
}

