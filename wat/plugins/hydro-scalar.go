package plugins

import (
	"encoding/json"
	"time"
)

type HydroScalarEvent struct {
	Plugin   string        `json:"plugin" yaml:"plugin"`
	Flows    []float64     `json:"flows"`
	TimeStep time.Duration `json:"timestep"`
	// FlowFrequency    statistics.BootstrappableDistribution `json:"flow_frequency"`
}

func (hse HydroScalarEvent) Payload() ([]byte, error) {
	return json.Marshal(hse)
}

func (hse *HydroScalarEvent) SetPlugin(pluginName string) {
	hse.Plugin = pluginName
}
