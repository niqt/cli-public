package cli_dep

type DependenciesForAppDto struct {
	Data  []Dependency `json:"data,omitempty"`
	Error string       `json:"error,omitempty"`
}
