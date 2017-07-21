package booklit

type Plugin interface {
	// methods are dynamically invoked
}

type PluginFactory func(*Section) Plugin

var plugins = map[string]PluginFactory{}

func RegisterPlugin(name string, factory PluginFactory) {
	plugins[name] = factory
}

func LookupPlugin(name string) (PluginFactory, bool) {
	plugin, found := plugins[name]
	return plugin, found
}
