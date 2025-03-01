package httpx

import (
	cli_dependencies "cli-depend/domain/cli-dep"
	"cli-depend/domain/cli-dep/infra/httpx"
	domain "cli-depend/domain/cli-dep/usecase"
	"cli-depend/domain/logger"
	"encoding/json"
	"net/http"
)

type UpdateDependencyDI struct {
	Usecase domain.UpdateDependencyUsecase
	Logger  logger.Logger
}

const UpdateDependencyPath = "/api/v1/package/{name}/{version}/dependency/{id}"

type updateDependencyHTTPX struct {
	usecase domain.UpdateDependencyUsecase
	logger  logger.Logger
}

func (u *updateDependencyHTTPX) Handler() http.HandlerFunc {
	return u.handleFunc
}

func NewUpdateDependencyHTTPX(di UpdateDependencyDI) httpx.Service {
	return &updateDependencyHTTPX{
		usecase: di.Usecase,
		logger:  di.Logger,
	}
}

func (u *updateDependencyHTTPX) Path() string {
	return UpdateDependencyPath
}

func (u *updateDependencyHTTPX) Method() httpx.Method {
	return httpx.MethodPut
}

func (u *updateDependencyHTTPX) handleFunc(w http.ResponseWriter, req *http.Request) {
	var dep cli_dependencies.Dependency
	err := json.NewDecoder(req.Body).Decode(&dep)
	if err != nil {
		u.logger.Error("Failed to decode")
		http.Error(w, "Errore durante la decodifica della dipendenza: "+err.Error(), http.StatusBadRequest)
		return
	}
	_, err = u.usecase.UpdateDependency(dep)
	var response cli_dependencies.DependenciesForAppDto
	if err != nil {
		u.logger.Error(err.Error())
		response = cli_dependencies.DependenciesForAppDto{
			Error: err.Error(),
		}
	} else {
		response = cli_dependencies.DependenciesForAppDto{
			Data: append([]cli_dependencies.Dependency{}, dep),
		}
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		u.logger.Error("Failed to encode")
		http.Error(w, "Error retrieving the response: "+err.Error(), http.StatusInternalServerError)
	}
}
