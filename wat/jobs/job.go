package jobs

// Makes life reeaaallly simple:
// https://zhwt.github.io/yaml-to-go/

type Job struct {
	Job struct {
		Config struct {
			Bucket    string `yaml:"bucket"`
			PrefixIn  string `yaml:"prefix_in"`
			PrefixOut string `yaml:"prefix_out"`
			Details   []struct {
				EventType       string `yaml:"event_type,omitempty"`
				Realizations    int    `yaml:"realizations,omitempty"`
				Lifecycles      int    `yaml:"lifecycles,omitempty"`
				RealizationSeed int    `yaml:"realization_seed,omitempty"`
				EventSeed       int    `yaml:"event_seed,omitempty"`
			} `yaml:"details"`
		} `yaml:"config"`
		Plugins []struct {
			Plugin struct {
				Name      string   `yaml:"name"`
				Payload   string   `yaml:"payload"`
				DependsOn []string `yaml:"depends on,omitempty"`
			} `yaml:"plugin"`
		} `yaml:"plugins"`
	} `yaml:"job"`
}
