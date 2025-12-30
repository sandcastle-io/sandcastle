package server

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/mbobrovskyi/kube-sandcastle/internal/executor"
	"github.com/mbobrovskyi/kube-sandcastle/pkg/api"
)

const (
	executeTimeout = 15 * time.Second
)

type Executor interface {
	Execute(ctx context.Context, code string) (*executor.Result, error)
}

type WorkerHandler struct {
	executor Executor
}

func NewWorkerHandler(executor Executor) *WorkerHandler {
	return &WorkerHandler{
		executor: executor,
	}
}

func (h *WorkerHandler) HealthzHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("ok")); err != nil {
		log.Printf("Failed to write health check response: %v", err)
	}
}

func (h *WorkerHandler) ReadyzHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("ready")); err != nil {
		log.Printf("Failed to write health check response: %v", err)
	}
}

func (h *WorkerHandler) ExecuteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}

	// Read and decode the request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	req := &api.ExecuteRequest{}
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	log.Printf("Executing code chunk (%d bytes)", len(req.Code))

	ctx, cancel := context.WithTimeout(r.Context(), executeTimeout)
	defer cancel()

	result, err := h.executor.Execute(ctx, req.Code)
	if err != nil {
		log.Printf("Internal server error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ExecutorResultToExecuteResponse(result)); err != nil {
		log.Printf("Internal server error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ExecutorResultToExecuteResponse(result *executor.Result) *api.ExecuteResponse {
	if result == nil {
		return nil
	}
	return &api.ExecuteResponse{
		Stdout:   result.Stdout,
		Stderr:   result.Stderr,
		ExitCode: result.ExitCode,
		CpuTime:  result.CpuTime,
		Memory:   result.Memory,
	}
}
