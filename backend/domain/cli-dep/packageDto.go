package cli_dep

type PackageDto struct {
	Data  Package `json:"data,omitempty"`
	Error string  `json:"error,omitempty"`
}
