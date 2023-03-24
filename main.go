package main

import (
	"github.com/Niromash/niromash-api/services"
	"github.com/Niromash/niromash-api/utils/environment"
	"log"
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
