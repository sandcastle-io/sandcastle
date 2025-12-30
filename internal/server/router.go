package server

import "net/http"

func NewRouter(workerHandler *WorkerHandler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", workerHandler.HealthzHandler)
	mux.HandleFunc("/readyz", workerHandler.ReadyzHandler)
	mux.HandleFunc("/execute", workerHandler.ExecuteHandler)
	return setupMiddleware(mux)
}

func setupMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "Sandcastle-Worker")
		next.ServeHTTP(w, r)
	})
}
