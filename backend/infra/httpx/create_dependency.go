package httpx

import (
	cli_dependencies "cli-depend/domain/cli-dep"
	"cli-depend/domain/cli-dep/infra/httpx"
	domain "cli-depend/domain/cli-dep/usecase"
	"cli-depend/domain/logger"
	"encoding/json"
	"net/http"
)

type CreateDependencyDI struct {
	Usecase domain.CreateDependencyUsecase
	Logger  logger.Logger
}

const CreateDependencyPath = "/api/v1/package/{name}/{version}/dependency"

type createDependencyHTTPX struct {
	usecase domain.CreateDependencyUsecase
	logger  logger.Logger
}

func (c *createDependencyHTTPX) Handler() http.HandlerFunc {
	return c.handleFunc
}

func NewCreateDependencyHTTPX(di CreateDependencyDI) httpx.Service {
	return &createDependencyHTTPX{
		usecase: di.Usecase,
		logger:  di.Logger,
	}
}

func (c *createDependencyHTTPX) Path() string {
	return CreateDependencyPath
}

func (c *createDependencyHTTPX) Method() httpx.Method {
	return httpx.MethodPost
}

func (c *createDependencyHTTPX) handleFunc(w http.ResponseWriter, req *http.Request) {
	var dep cli_dependencies.Dependency
	err := json.NewDecoder(req.Body).Decode(&dep)
	if err != nil {
		c.logger.Info("Failed to decode")
		http.Error(w, "Error decoding the body: "+err.Error(), http.StatusInternalServerError)
		return
	}
	var insertId int64
	insertId, err = c.usecase.CreateDependency(dep)
	var response cli_dependencies.DependenciesForAppDto
	if err != nil {
		response = cli_dependencies.DependenciesForAppDto{
			Error: err.Error(),
		}
	} else {
		dep.ID = insertId
		response = cli_dependencies.DependenciesForAppDto{
			Data: append([]cli_dependencies.Dependency{}, dep),
		}
	}
	json.NewEncoder(w).Encode(response)
}
