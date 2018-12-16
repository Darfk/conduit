package conduit

import (
	"fmt"
)

type Config struct {
	Stages []Stage
}

type Stage struct {
	Route  int
	Grow   int
	Shrink int
	Size   int
}

func CheckConfig(cfg Config) error {
	usedRoutes := make(map[int]struct{})
	for _, stage := range cfg.Stages {
		if _, exists := usedRoutes[stage.Route]; exists {
			return fmt.Errorf("stage already exists at route %d", stage.Route)
		}
		usedRoutes[stage.Route] = struct{}{}

		if stage.Grow < 1 {
			return fmt.Errorf("stage.Grow < 1")
		}
		if stage.Shrink < 1 {
			return fmt.Errorf("stage.Shrink < 1")
		}
		if stage.Size < 1 {
			return fmt.Errorf("stage.Size < 1")
		}
	}
	return nil
}
