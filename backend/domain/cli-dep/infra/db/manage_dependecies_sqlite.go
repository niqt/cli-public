package db

import (
	cli "cli-depend/domain/cli-dep"
	"cli-depend/domain/cli-dep/repository"
	"cli-depend/domain/logger"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strconv"
)

type ManageDependenciesRepositoryDI struct {
	Logger logger.Logger
}

type manageDependenciesRepository struct {
	logger logger.Logger
}

func NewDBRepository(di ManageDependenciesRepositoryDI) repository.ManageDependenciesRepo {
	return &manageDependenciesRepository{logger: di.Logger}
}

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./app.db")
	if err != nil {
		log.Fatal(err) // Log an error and stop the program if the database can't be opened
	}

	createPackageSql := `
 CREATE TABLE IF NOT EXISTS package (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  name TEXT, version TEXT
 );`

	_, err = DB.Exec(createPackageSql)
	if err != nil {
		log.Fatalf("Error creating table: %q: %s\n", err, createPackageSql) // Log an error if table creation fails
	}

	createDependenciesSql := `
 CREATE TABLE IF NOT EXISTS dependency (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  name TEXT, version TEXT, score NUMERIC, packageId INTEGER
 );`

	_, err = DB.Exec(createDependenciesSql)
	if err != nil {
		log.Fatalf("Error creating table: %q: %s\n", err, createDependenciesSql) // Log an error if table creation fails
	}
}

func (m manageDependenciesRepository) SaveDependency(dependency cli.Dependency) (int64, error) {
	dep, err := m.GetDependencyByPackage(dependency)
	if err != nil {
		m.logger.Error(err.Error())
		return 0, fmt.Errorf("errore nel recupero della dipendenza: %w", err)
	}

	if dep == (cli.Dependency{}) {
		return m.InsertDependency(dependency)
	}

	err = m.UpdateDependency(dependency)
	if err != nil {
		m.logger.Error(err.Error())
		return 0, err
	}
	return dep.ID, nil
}

func (m manageDependenciesRepository) GetDependencies(packageId int64, searchName, lower, upper string) ([]cli.Dependency, error) {
	minScore := 0.0
	maxScore := 10.0
	if len(lower) != 0 {
		n, err := strconv.ParseFloat(lower, 32)
		if err != nil {
			m.logger.Error("lower must be a number")
			n = 0
		}
		minScore = n
	}
	if len(lower) != 0 {
		n, err := strconv.ParseFloat(upper, 32)
		if err != nil {
			m.logger.Error("upper must be a number")
			n = 10
		}
		maxScore = n
	}

	searchTextQueryParameter := "%"
	if len(searchName) != 0 {
		searchTextQueryParameter = "%" + searchName + "%"
	}

	rows, err := DB.Query("SELECT * FROM dependency where packageId = ? and name LIKE ? and ((score >= ? and score <= ?) or score < ?)", packageId,
		searchTextQueryParameter, minScore, maxScore, 0)

	if err != nil {
		m.logger.Error(err.Error())
		return []cli.Dependency{}, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			m.logger.Error(err.Error())
		}
	}(rows)

	var dependencies []cli.Dependency
	for rows.Next() {
		var dep cli.Dependency
		if err = rows.Scan(&dep.ID, &dep.Name, &dep.Version, &dep.Score, &dep.PackageID); err != nil {
			m.logger.Error(err.Error())
			return []cli.Dependency{}, err
		}
		dependencies = append(dependencies, dep)
	}
	return dependencies, nil
}

func (m manageDependenciesRepository) GetPackage(name, version string) (cli.Package, error) {

	rows, err := DB.Query("SELECT * FROM package WHERE name = ? and version = ?", name, version)
	pkg := cli.Package{}
	if err != nil {
		fmt.Println("NOT FOUND", err)
		return pkg, nil
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			m.logger.Error(err.Error())
		}
	}(rows)

	for rows.Next() {
		if err := rows.Scan(&pkg.ID, &pkg.Name, &pkg.Version); err != nil {
			m.logger.Error(err.Error())
			return cli.Package{}, err
		}

	}
	return pkg, nil
}

func (m manageDependenciesRepository) SavePackage(name, version string) (int64, error) {
	pkg, err := m.GetPackage(name, version)
	if err != nil {
		m.logger.Error(err.Error())
		return 0, fmt.Errorf("error retrieving the package: %w", err)
	}

	if pkg == (cli.Package{}) {
		return m.insertPackage(name, version)
	}

	err = m.updatePackage(name, version, pkg.ID)
	if err != nil {
		m.logger.Error(err.Error())
		return 0, err
	}

	return int64(pkg.ID), nil
}

func (m manageDependenciesRepository) GetDependencyByPackage(dependency cli.Dependency) (cli.Dependency, error) {
	rows, err := DB.Query("SELECT * FROM dependency WHERE name = ? and version = ? and packageId = ?", dependency.Name, dependency.Version, dependency.PackageID)
	dep := cli.Dependency{}
	if err != nil {
		m.logger.Error(err.Error())
		return dep, nil
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			m.logger.Error(err.Error())
		}
	}(rows)
	for rows.Next() {
		if err := rows.Scan(&dep.ID, &dep.Name, &dep.Version, &dep.Score, &dep.PackageID); err != nil {
			m.logger.Error(err.Error())
			return cli.Dependency{}, nil
		}
	}
	return dep, nil
}

func (m manageDependenciesRepository) InsertDependency(dependency cli.Dependency) (int64, error) {
	stmt, err := DB.Prepare("INSERT INTO dependency (name, version, score, packageId) VALUES (?,?,?, ?)")
	if err != nil {
		m.logger.Error(err.Error())
		return 0, fmt.Errorf("insert prepare: %w", err)
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			m.logger.Error(err.Error())
		}
	}(stmt)

	var result sql.Result
	result, err = stmt.Exec(dependency.Name, dependency.Version, dependency.Score, dependency.PackageID)
	if err != nil {
		m.logger.Error(err.Error())
		return 0, fmt.Errorf("error inserting: %w", err)
	}

	return result.LastInsertId()
}

func (m manageDependenciesRepository) UpdateDependency(dependency cli.Dependency) error {
	stmt, err := DB.Prepare("UPDATE dependency SET version = ?, score = ? WHERE id = ?")
	if err != nil {
		m.logger.Error(err.Error())
		return fmt.Errorf("error prepare update: %w", err)
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			m.logger.Error(err.Error())
		}
	}(stmt)

	_, err = stmt.Exec(dependency.Version, dependency.Score, dependency.ID)
	if err != nil {
		m.logger.Error(err.Error())
		return fmt.Errorf("error upgrading: %w", err)
	}

	return nil
}

func (m manageDependenciesRepository) insertPackage(name, version string) (int64, error) {
	stmt, err := DB.Prepare("INSERT INTO package (name, version) VALUES (?,?)")
	if err != nil {
		m.logger.Error(err.Error())
		return 0, fmt.Errorf("errore  preparing insert: %w", err)
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			m.logger.Error(err.Error())
		}
	}(stmt)
	var result sql.Result
	result, err = stmt.Exec(name, version)
	if err != nil {
		m.logger.Error(err.Error())
		return 0, fmt.Errorf("error inserting the package: %w", err)
	}

	return result.LastInsertId()
}

func (m manageDependenciesRepository) updatePackage(name, version string, id int) error {
	stmt, err := DB.Prepare("UPDATE package SET version = ?, name = ? WHERE id = ?")
	if err != nil {
		m.logger.Error(err.Error())
		return fmt.Errorf("error preparing update: %w", err)
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			m.logger.Error(err.Error())
		}
	}(stmt)

	_, err = stmt.Exec(version, name, id)
	if err != nil {
		m.logger.Error(err.Error())
		return fmt.Errorf("error upgrading the package: %w", err)
	}

	return nil
}

func (m manageDependenciesRepository) DeleteDependency(id int64) (int64, error) {
	stmt, err := DB.Prepare("DELETE FROM DEPENDENCY WHERE id = ?")
	if err != nil {
		m.logger.Error(err.Error())
		return 0, fmt.Errorf("error preparing update: %w", err)
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			m.logger.Error(err.Error())
		}
	}(stmt)

	_, err = stmt.Exec(id)
	if err != nil {
		m.logger.Error(err.Error())
		return 0, fmt.Errorf("error upgrading the package: %w", err)
	}

	return id, nil
}
