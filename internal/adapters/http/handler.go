package http

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path/filepath"

	yamlparser "github.com/tahaSadeghi2308/TapeRunner/internal/adapters/yaml_parser"
	"github.com/tahaSadeghi2308/TapeRunner/internal/application"
)

const MACHINES_FOLDER string = "machines"

type Handler struct {
	simService  *application.SimulationService
	machinesDir string
}

func NewHandler(simService *application.SimulationService) *Handler {
	return &Handler{
		simService:  simService,
		machinesDir: MACHINES_FOLDER,
	}
}

func (h *Handler) ServeUI(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func (h *Handler) HandleListMachines(w http.ResponseWriter, r *http.Request) {
	parser := yamlparser.NewYamlParser()
	machines, err := parser.ListMachines(h.machinesDir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(machines)
}

type RunRequest struct {
	MachineName string `json:"machine_name"`
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

	if req.MachineName == "" {
		http.Error(w, "machine_name is required", http.StatusBadRequest)
		return
	}

	machinePath := filepath.Join(h.machinesDir, req.MachineName)
	machine, err := h.simService.MachineSetup.Read(machinePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	results := h.simService.Run(machine, req.InitialTape)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
