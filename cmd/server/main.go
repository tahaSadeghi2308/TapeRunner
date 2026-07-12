package main

import (
	"log"
	"net/http"

	"github.com/tahaSadeghi2308/TapeRunner/internal/application"
	myhttp "github.com/tahaSadeghi2308/TapeRunner/internal/adapters/http"
)

func main() {
	simService := application.NewSimulationService()
	
	handler := myhttp.NewHandler(simService) // myhttp alias depending on imports
	router := myhttp.SetupRouter(handler)

	log.Println("Server starting on http://localhost:3225...")
	err := http.ListenAndServe(":3225", router)
	if err != nil {
		log.Fatal("Server failed to start:", err)
	}
}