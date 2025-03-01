package repository

import (
	cli "cli-depend/domain/cli-dep"
)

type ManageDependenciesRepo interface {
	GetDependencies(packageId int64, searchName, lower, upper string) ([]cli.Dependency, error)
	SaveDependency(dependency cli.Dependency) (int64, error)
	GetPackage(name, version string) (cli.Package, error)
	SavePackage(name, version string) (int64, error)
	GetDependencyByPackage(dependency cli.Dependency) (cli.Dependency, error)
	InsertDependency(dependency cli.Dependency) (int64, error)
	UpdateDependency(dependency cli.Dependency) error
	DeleteDependency(id int64) (int64, error)
}
