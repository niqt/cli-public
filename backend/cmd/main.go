package main

import (
	repo "cli-depend/domain/cli-dep/infra/db"
	"cli-depend/domain/logger"
	"cli-depend/infra/httpx"
	"cli-depend/usecase"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
)

func main() {
	logger.GetLogger().Info("Starting")
	repo.InitDB()
	defer repo.DB.Close()
	router := mux.NewRouter()
	router.UseEncodedPath()
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"accept", "authorization", "content-type"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	})

	dbRepository := repo.NewDBRepository(repo.ManageDependenciesRepositoryDI{
		Logger: *logger.GetLogger(),
	})
	createDependencyUsecase := usecase.NewCreateDependency(
		usecase.DependenciesUsecaseDI{
			Repository: dbRepository,
			Logger:     *logger.GetLogger(),
		})
	createDependencyService := httpx.NewCreateDependencyHTTPX(httpx.CreateDependencyDI{
		Usecase: createDependencyUsecase,
		Logger:  *logger.GetLogger(),
	})

	updateUsecase := usecase.NewUpdateDependency(usecase.DependenciesUsecaseDI{
		Repository: dbRepository,
		Logger:     *logger.GetLogger(),
	})

	updateService := httpx.NewUpdateDependencyHTTPX(httpx.UpdateDependencyDI{
		Usecase: updateUsecase,
		Logger:  *logger.GetLogger(),
	})

	getDependenciesUsecase := usecase.NewGetDependencies(usecase.DependenciesUsecaseDI{
		Repository: dbRepository,
		Logger:     *logger.GetLogger(),
	})

	getDependenciesService := httpx.NewGetDependenciesHTTPX(httpx.GetDependenciesDI{
		Usecase: getDependenciesUsecase,
		Logger:  *logger.GetLogger(),
	})

	getPackageusecase := usecase.NewGetPackage(usecase.DependenciesUsecaseDI{
		Repository: dbRepository,
		Logger:     *logger.GetLogger(),
	})

	getPackageService := httpx.NewGetPackageHTTPX(httpx.GetPackageDI{
		Usecase: getPackageusecase,
		Logger:  *logger.GetLogger(),
	})

	downloadUsecase := usecase.NewDownloadUsecase(usecase.DependenciesUsecaseDI{
		Repository: dbRepository,
		Logger:     *logger.GetLogger(),
	})

	downloadService := httpx.NewDownloadHTTPX(httpx.DownloadDI{
		Usecase: downloadUsecase,
		Logger:  *logger.GetLogger(),
	})

	deleteUsecase := usecase.NewDeleteDependency(usecase.DependenciesUsecaseDI{
		Repository: dbRepository,
		Logger:     *logger.GetLogger(),
	})

	deleteService := httpx.NewDeleteDependencyHTTPX(httpx.DeleteDependencyDI{
		Usecase: deleteUsecase,
		Logger:  *logger.GetLogger(),
	})

	healthService := httpx.NewHealthHTTPX(httpx.HealthDI{
		Logger: *logger.GetLogger(),
	})

	router.HandleFunc(createDependencyService.Path(), createDependencyService.Handler()).Methods(createDependencyService.Method().String())
	router.HandleFunc(updateService.Path(), updateService.Handler()).Methods(updateService.Method().String())
	router.HandleFunc(getDependenciesService.Path(), getDependenciesService.Handler()).Methods(getDependenciesService.Method().String())
	router.HandleFunc(getPackageService.Path(), getPackageService.Handler()).Methods(getPackageService.Method().String())
	router.HandleFunc(downloadService.Path(), downloadService.Handler()).Methods(downloadService.Method().String())
	router.HandleFunc(deleteService.Path(), deleteService.Handler()).Methods(deleteService.Method().String())

	router.HandleFunc(healthService.Path(), healthService.Handler()).Methods(healthService.Method().String())
	err := http.ListenAndServe("0.0.0.0:8080", c.Handler(router))
	if err != nil {
		logger.GetLogger().Error(err.Error())
	}
	logger.GetLogger().Info("Died")
}
