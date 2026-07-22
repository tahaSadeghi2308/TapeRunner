package http

import "net/http"

func SetupRouter(handler *Handler) *http.ServeMux {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", handler.ServeUI)
	mux.HandleFunc("/api/machines", handler.HandleListMachines)
	mux.HandleFunc("/api/run", handler.HandleRun)

	return mux
}