package httpx

import (
	"cli-depend/domain/cli-dep/infra/httpx"
	"cli-depend/domain/logger"
	"net/http"
)

type HealthDI struct {
	Logger logger.Logger
}

const HealthPath = "/health"

type healthHTTPX struct {
	logger logger.Logger
}

func (h *healthHTTPX) Handler() http.HandlerFunc {
	return h.handleFunc
}

func NewHealthHTTPX(di HealthDI) httpx.Service {
	return &healthHTTPX{
		logger: di.Logger,
	}
}

func (h *healthHTTPX) Path() string {
	return HealthPath
}

func (h *healthHTTPX) Method() httpx.Method {
	return httpx.MethodGet
}

func (h *healthHTTPX) handleFunc(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
