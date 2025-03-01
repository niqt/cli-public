package cli_dep

type VersionKey struct {
	System  string `json:"system"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Node struct {
	VersionKey VersionKey `json:"versionKey"`
	Bundled    bool       `json:"bundled"`
	Relation   string     `json:"relation"`
	Errors     []string   `json:"errors"`
}

type DependenciesDto struct {
	Nodes []Node `json:"nodes"`
}
