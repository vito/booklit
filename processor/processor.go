package processor

type Processor struct {
	Plugins []Plugin
}

type Plugin interface {
	// methods are dynamically invoked
}

func (processor Processor) Load(path string) *backlit.Section {
}
