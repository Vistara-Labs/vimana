package config

type SpacecoreType string

const (
	Celestia  SpacecoreType = "celestia"
	Berachain SpacecoreType = "berachain"
	Avail     SpacecoreType = "avail"
	Gmworld   SpacecoreType = "gmworld"
	Eigen     SpacecoreType = "eigen"
	RepoBase  string        = "https://github.com/vistara-labs/spacecore-template.git"
)
