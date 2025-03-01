package httpx

import (
	cli_dependencies "cli-depend/domain/cli-dep"
	"cli-depend/domain/cli-dep/infra/httpx"
	domain "cli-depend/domain/cli-dep/usecase"
	"cli-depend/domain/logger"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"net/url"
)

type GetDependenciesDI struct {
	Usecase domain.GetDependenciesUsecase
	Logger  logger.Logger
}

const GetDependenciesPath = "/api/v1/package/{name}/{version}/dependency"

type getDependenciesHTTPX struct {
	usecase domain.GetDependenciesUsecase
	logger  logger.Logger
}

func (g *getDependenciesHTTPX) Handler() http.HandlerFunc {
	return g.handleFunc
}

func NewGetDependenciesHTTPX(di GetDependenciesDI) httpx.Service {
	return &getDependenciesHTTPX{
		usecase: di.Usecase,
		logger:  di.Logger,
	}
}

func (g *getDependenciesHTTPX) Path() string {
	return GetDependenciesPath
}

func (g *getDependenciesHTTPX) Method() httpx.Method {
	return httpx.MethodGet
}

func (g *getDependenciesHTTPX) handleFunc(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	packageName := params["name"]
	packageVersion := params["version"]

	v := req.URL.Query()

	searchName := v.Get("searchName")
	lower := v.Get("lower")
	upper := v.Get("upper")

	decodedName, err := url.QueryUnescape(packageName)
	if err != nil {
		g.logger.Error(err.Error())
		http.Error(w, "Error decoding", http.StatusInternalServerError)
		return
	}
	var dependencies []cli_dependencies.Dependency
	var response cli_dependencies.DependenciesForAppDto
	dependencies, err = g.usecase.GetDependencies(decodedName, packageVersion, searchName, lower, upper)

	if err != nil {
		response = cli_dependencies.DependenciesForAppDto{
			Error: err.Error(),
		}
	} else {
		response = cli_dependencies.DependenciesForAppDto{
			Data: dependencies,
		}
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		g.logger.Info("Failed to encode")
		http.Error(w, "Error retrieving the response: "+err.Error(), http.StatusInternalServerError)
	}
}
