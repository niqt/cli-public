package httpx

import (
	cli_dependencies "cli-depend/domain/cli-dep"
	"cli-depend/domain/cli-dep/infra/httpx"
	domain "cli-depend/domain/cli-dep/usecase"
	"cli-depend/domain/logger"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type DeleteDependencyDI struct {
	Usecase domain.DeleteDependencyUsecase
	Logger  logger.Logger
}

const DeleteDependencyPath = "/api/v1/package/{name}/{version}/dependency/{id}"

type deleteDependencyHTTPX struct {
	usecase domain.DeleteDependencyUsecase
	logger  logger.Logger
}

func (d *deleteDependencyHTTPX) Handler() http.HandlerFunc {
	return d.handleFunc
}

func NewDeleteDependencyHTTPX(di DeleteDependencyDI) httpx.Service {
	return &deleteDependencyHTTPX{
		usecase: di.Usecase,
		logger:  di.Logger,
	}
}

func (d *deleteDependencyHTTPX) Path() string {
	return DeleteDependencyPath
}

func (d *deleteDependencyHTTPX) Method() httpx.Method {
	return httpx.MethodDelete
}

func (d *deleteDependencyHTTPX) handleFunc(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		d.logger.Info("Failed to parse id")
		http.Error(w, "Error parsing the id: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = d.usecase.DeleteDependency(int64(id))
	var response cli_dependencies.DependenciesForAppDto
	if err != nil {
		response = cli_dependencies.DependenciesForAppDto{
			Error: err.Error(),
		}
	} else {
		response = cli_dependencies.DependenciesForAppDto{
			Data: []cli_dependencies.Dependency{},
		}
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		d.logger.Info("Failed to encode")
	}
}
