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

type GetPackageDI struct {
	Usecase domain.GetPackageUsecase
	Logger  logger.Logger
}

const GetPackagePath = "/api/v1/package/{name}/{version}"

type getPackageHTTPX struct {
	usecase domain.GetPackageUsecase
	logger  logger.Logger
}

func (g *getPackageHTTPX) Handler() http.HandlerFunc {
	return g.handleFunc
}

func NewGetPackageHTTPX(di GetPackageDI) httpx.Service {
	return &getPackageHTTPX{
		usecase: di.Usecase,
		logger:  di.Logger,
	}
}

func (g *getPackageHTTPX) Path() string {
	return GetDependenciesPath
}

func (g *getPackageHTTPX) Method() httpx.Method {
	return httpx.MethodGet
}

func (g *getPackageHTTPX) handleFunc(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	packageName, _ := url.QueryUnescape(params["name"])
	packageVersion, _ := url.QueryUnescape(params["version"])
	pkg, err := g.usecase.GetPackage(packageName, packageVersion)

	var response cli_dependencies.PackageDto
	if err != nil {
		g.logger.Error(err.Error())
		response = cli_dependencies.PackageDto{
			Error: err.Error(),
		}
	} else {
		response = cli_dependencies.PackageDto{
			Data: pkg,
		}
	}
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		g.logger.Error("Failed to encode")
		http.Error(w, "Error retrieving the response: "+err.Error(), http.StatusInternalServerError)
	}
}
