package usecase

import (
	"cli-depend/api/client"
	cli "cli-depend/domain/cli-dep"
	repo "cli-depend/domain/cli-dep/repository"
	domain "cli-depend/domain/cli-dep/usecase"
	"cli-depend/domain/logger"
	"fmt"
	"sync"
)

type DependenciesUsecaseDI struct {
	Repository repo.ManageDependenciesRepo
	Logger     logger.Logger
}

type dependenciesUsecase struct {
	repository repo.ManageDependenciesRepo
	logger     logger.Logger
}

func (d *dependenciesUsecase) Download(packageName, packageVersion string) ([]cli.Dependency, error) {
	pkg, err := d.repository.GetPackage(packageName, packageVersion)
	if err != nil {
		fmt.Println("Loading error", err)
		return []cli.Dependency{}, err
	}

	var packageId int64
	if pkg == (cli.Package{}) {
		packageId, err = d.repository.SavePackage(packageName, packageVersion)
	} else {
		packageId = int64(pkg.ID)
	}

	packageInfo := cli.Package{
		ID:      int(packageId),
		Name:    packageName,
		Version: packageVersion,
	}
	ch := make(chan cli.Dependency, 5)
	done := make(chan bool)

	var wg sync.WaitGroup
	wg.Add(2)
	go d.readDependency(packageInfo, ch, done, &wg)
	go d.writeDependency(ch, done, &wg)
	wg.Wait()
	return d.GetDependencies(packageName, packageVersion, "", "", "")
}

func (d *dependenciesUsecase) readDependency(packageInfo cli.Package, ch chan<- cli.Dependency, done chan<- bool, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(ch)

	dependencies, err := client.GetDependencies(packageInfo.Name, packageInfo.Version)
	if err != nil {
		d.logger.Error("Failed to get dependencies")
		return
	}
	for _, dep := range dependencies.Nodes {
		score, scoreError := client.GetScore(dep.VersionKey.Name)
		if scoreError != nil {
			score = -1
		}
		var dependency = cli.Dependency{
			Name:      dep.VersionKey.Name,
			Version:   dep.VersionKey.Version,
			Score:     score,
			PackageID: int64(packageInfo.ID),
		}
		ch <- dependency
	}
	done <- true
}

func (d *dependenciesUsecase) writeDependency(ch <-chan cli.Dependency, done <-chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case dep, ok := <-ch:
			if !ok {
				return
			}
			_, err := d.repository.SaveDependency(dep)
			if err != nil {
				d.logger.Error("Failed to save dependencies")
				return
			}

		case <-done:
			for dep := range ch {
				_, err := d.repository.SaveDependency(dep)
				if err != nil {
					d.logger.Error("Failed to save dependencies")
					return
				}
			}
			return
		}
	}
}

func (d *dependenciesUsecase) GetDependencies(packageName, packageVersion, searchName, lower, upper string) ([]cli.Dependency, error) {
	pkg, err := d.repository.GetPackage(packageName, packageVersion)
	if err != nil {
		d.logger.Error(err.Error())
		return nil, fmt.Errorf("package not found: %w", err)
	}
	if pkg != (cli.Package{}) {
		return d.repository.GetDependencies(int64(pkg.ID), searchName, lower, upper)
	}
	return []cli.Dependency{}, nil
}

func (d *dependenciesUsecase) UpdateDependency(dep cli.Dependency) (int64, error) {
	return d.repository.SaveDependency(dep)
}

func (d *dependenciesUsecase) GetPackage(packageName, packageVersion string) (cli.Package, error) {
	pkg, err := d.repository.GetPackage(packageName, packageVersion)
	if err != nil {
		d.logger.Error(err.Error())
		return cli.Package{}, fmt.Errorf("package not found: %w", err)
	}
	return pkg, nil
}

func (d *dependenciesUsecase) CreateDependency(dep cli.Dependency) (int64, error) {
	insertedId, err := d.repository.InsertDependency(dep)
	if err != nil {
		d.logger.Error(err.Error())
		return 0, fmt.Errorf("package not found: %w", err)
	}
	return insertedId, nil
}

func (d *dependenciesUsecase) DeleteDependency(id int64) (int64, error) {
	deletedId, err := d.repository.DeleteDependency(id)
	if err != nil {
		d.logger.Error(err.Error())
		return 0, fmt.Errorf("package not found: %w", err)
	}
	return deletedId, nil
}

func NewDownloadUsecase(di DependenciesUsecaseDI) domain.DownloadUsecase {
	return &dependenciesUsecase{
		repository: di.Repository,
		logger:     di.Logger,
	}
}

func NewGetDependencies(di DependenciesUsecaseDI) domain.GetDependenciesUsecase {
	return &dependenciesUsecase{
		repository: di.Repository,
		logger:     di.Logger,
	}
}

func NewUpdateDependency(di DependenciesUsecaseDI) domain.UpdateDependencyUsecase {
	return &dependenciesUsecase{
		repository: di.Repository,
		logger:     di.Logger,
	}
}

func NewGetPackage(di DependenciesUsecaseDI) domain.GetPackageUsecase {
	return &dependenciesUsecase{
		repository: di.Repository,
		logger:     di.Logger,
	}
}

func NewCreateDependency(di DependenciesUsecaseDI) domain.CreateDependencyUsecase {
	return &dependenciesUsecase{
		repository: di.Repository,
		logger:     di.Logger,
	}
}

func NewDeleteDependency(di DependenciesUsecaseDI) domain.DeleteDependencyUsecase {
	return &dependenciesUsecase{
		repository: di.Repository,
		logger:     di.Logger,
	}
}
