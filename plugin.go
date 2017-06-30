package booklit

type PluginFactory interface {
	NewPlugin(*Section) Plugin
}

type Plugin interface {
	// methods are dynamically invoked
}
