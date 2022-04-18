package server

import (
	"fmt"
	"wat/plugins"
)

func registeredPlugins(pluginName string) (plugins.Manifest, error) {
	var event plugins.Manifest
	switch pluginName {
	case "hec-ras":
		event = new(plugins.RasEvent)

	case "hydro-scalar":
		event = new(plugins.HydroScalarEvent)

	case "consequences":
		event = new(plugins.ConsequencesEvent)

	default:
		return event, fmt.Errorf("plugin not supported: %v", pluginName)
	}
	return event, nil
}
