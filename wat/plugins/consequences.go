package plugins

import "encoding/json"

type ConsequencesEvent struct {
	Plugin string `json:"plugin" yaml:"plugin"`
	Path   string `json:"path" yaml:"path"`
	Source string `json:"source" yaml:"source"`
}

func (ce ConsequencesEvent) Payload() ([]byte, error) {
	return json.Marshal(ce)
}

func (ce *ConsequencesEvent) SetPlugin(pluginName string) {
	ce.Plugin = pluginName
}
