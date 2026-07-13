package main

import (
	yamlparser "github.com/tahaSadeghi2308/TapeRunner/internal/adapters/yaml_parser"
	"log"
	"net/http"

	myhttp "github.com/tahaSadeghi2308/TapeRunner/internal/adapters/http"
	"github.com/tahaSadeghi2308/TapeRunner/internal/application"
)

func main() {
	machineSetup := yamlparser.NewYamlParser()
	simService := application.NewSimulationService(&machineSetup)

	handler := myhttp.NewHandler(simService) // myhttp alias depending on imports
	router := myhttp.SetupRouter(handler)

	log.Println("Server starting on http://localhost:3225...")
	err := http.ListenAndServe(":3225", router)
	if err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
