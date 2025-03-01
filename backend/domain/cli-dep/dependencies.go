package cli_dep

type Dependency struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Version   string  `json:"version"`
	Score     float32 `json:"score"`
	PackageID int64   `json:"packageId"`
}
