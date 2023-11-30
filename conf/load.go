package conf

import "fmt"

func Load(verbose bool, loaders ...Loader) Map {
	cfg := Map{}

	for _, loader := range loaders {
		cfg = loader.Load(cfg)
	}

	if verbose {
		fmt.Println("Configuration verbose: ON")
	}
	return cfg.Interpolate(verbose)
}
