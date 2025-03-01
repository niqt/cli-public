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

type DownloadDI struct {
	Usecase domain.DownloadUsecase
	Logger  logger.Logger
}

const DownloadPath = "/api/v1/package/{name}/{version}"

type downloadHTTPX struct {
	usecase domain.DownloadUsecase
	logger  logger.Logger
}

func (d *downloadHTTPX) Handler() http.HandlerFunc {
	return d.handleFunc
}

func NewDownloadHTTPX(di DownloadDI) httpx.Service {
	return &downloadHTTPX{
		usecase: di.Usecase,
		logger:  di.Logger,
	}
}

func (d *downloadHTTPX) Path() string {
	return DownloadPath
}

func (d *downloadHTTPX) Method() httpx.Method {
	return httpx.MethodGet
}

func (d *downloadHTTPX) handleFunc(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	packageName, _ := url.QueryUnescape(params["name"])
	packageVersion, _ := url.QueryUnescape(params["version"])
	dependencies, err := d.usecase.Download(packageName, packageVersion)
	var response cli_dependencies.DependenciesForAppDto

	if err != nil {
		d.logger.Error(err.Error())
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
		d.logger.Info("Failed to encode")
		http.Error(w, "Error retrieving the response: "+err.Error(), http.StatusInternalServerError)
	}
}
