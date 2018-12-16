package conduit

import (
	"testing"
)

func TestConfigShouldPass(t *testing.T) {
	cfg := Config{
		Stages: []Stage{
			Stage{
				Route: 1, Grow: 1, Shrink: 1, Size: 1,
			},
			Stage{
				Route: 2, Grow: 32, Shrink: 16, Size: 4,
			},
		},
	}

	if err := CheckConfig(cfg); err != nil {
		t.Error(err)
	}

}

func TestConfigShouldFail(t *testing.T) {
	cfgs := []Config{
		Config{Stages: []Stage{Stage{Route: 1, Grow: 0, Shrink: 1, Size: 1}}},
		Config{Stages: []Stage{Stage{Route: 1, Grow: 1, Shrink: 0, Size: 1}}},
		Config{Stages: []Stage{Stage{Route: 1, Grow: 1, Shrink: 1, Size: 0}}},
		Config{
			Stages: []Stage{
				Stage{Route: 1, Grow: 1, Shrink: 1, Size: 1},
				Stage{Route: 1, Grow: 1, Shrink: 1, Size: 1},
			},
		},
	}

	for _, cfg := range cfgs {
		if err := CheckConfig(cfg); err == nil {
			t.Error("bad config did not return error")
		}
	}

}
