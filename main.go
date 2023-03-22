package main

import (
	"log"
	"niromash-api/services"
	"niromash-api/utils/environment"
)

func main() {
	if !environment.CheckEnvs() {
		log.Fatalln("Missing environment variables")
	}
	service := services.NewMainService()
	if err := service.Init(); err != nil {
		log.Fatalln(err)
	}
	errCh := make(chan error)
	service.Start(errCh)
	defer service.Close()

	log.Fatalln(<-errCh)
}
