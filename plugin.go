package booklit

type PluginFactory interface {
	NewPlugin(*Section) Plugin
}

var plugins = map[string]PluginFactory{}

func RegisterPlugin(name string, factory PluginFactory) {
	plugins[name] = factory
}

func LookupPlugin(name string) (PluginFactory, bool) {
	plugin, found := plugins[name]
	return plugin, found
}

type PluginFactoryFunc func(*Section) Plugin

func (f PluginFactoryFunc) NewPlugin(s *Section) Plugin { return f(s) }

type Plugin interface {
	// methods are dynamically invoked
}
