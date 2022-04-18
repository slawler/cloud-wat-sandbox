package plugins

type Manifest interface {
	Payload() ([]byte, error)
	SetPlugin(string)
}
