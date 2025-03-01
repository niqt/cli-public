package usecase

import (
	cli "cli-depend/domain/cli-dep"
)

type ManageDependenciesUsecase interface {
	DownloadUsecase
	GetDependenciesUsecase
	UpdateDependencyUsecase
	GetPackageUsecase
	CreateDependencyUsecase
	DeleteDependencyUsecase
}

type DownloadUsecase interface {
	Download(packageName, packageVersion string) ([]cli.Dependency, error)
}

type GetDependenciesUsecase interface {
	GetDependencies(packageName, packageVersion, searchName, lower, upper string) ([]cli.Dependency, error)
}

type UpdateDependencyUsecase interface {
	UpdateDependency(dep cli.Dependency) (int64, error)
}

type GetPackageUsecase interface {
	GetPackage(packageName, packageVersion string) (cli.Package, error)
}

type CreateDependencyUsecase interface {
	CreateDependency(dep cli.Dependency) (int64, error)
}

type DeleteDependencyUsecase interface {
	DeleteDependency(id int64) (int64, error)
}
