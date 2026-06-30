package http

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/tahaSadeghi2308/TapeRunner/internal/adapters/parser"
	"github.com/tahaSadeghi2308/TapeRunner/internal/application"
)

type Handler struct {
	simService *application.SimulationService
}

func NewHandler(simService *application.SimulationService) *Handler {
	return &Handler{simService: simService}
}

func (h *Handler) ServeUI(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

type RunRequest struct {
	MachineJSON string `json:"machine_json"`
	InitialTape string `json:"initial_tape"`
}

func (h *Handler) HandleRun(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RunRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	machine, err := parser.ParseMachine([]byte(req.MachineJSON))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	results := h.simService.Run(machine, req.InitialTape)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}