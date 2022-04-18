package plugins

import (
	"encoding/json"
)

type RasEvent struct {
	Plugin          string `json:"plugin" yaml:"plugin"`
	BasePath        string `json:"basepath"`
	ProjectFilePath string `json:"projectfile"`
	Pfile           string `json:"planfile"`
	// Ufile            RasFlowFile          `json:"unsteadyfile"`
	// Links            component.ModelLinks `json:"-"`
}

func (re RasEvent) Payload() ([]byte, error) {
	return json.Marshal(re)
}

func (re *RasEvent) SetPlugin(pluginName string) {
	re.Plugin = pluginName
}
