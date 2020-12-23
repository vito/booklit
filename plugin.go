package booklit

import "github.com/sirupsen/logrus"

// Plugin is an arbitrary object which is initialized with the Section that is
// using the plugin.
//
// Methods on the plugin object will be invoked during the "evaluation" stage
// by function calling syntax in Booklit documents.
//
// See https://booklit.page/plugins.html for more information.
type Plugin interface {
	// methods are dynamically invoked
}

// PluginFactory constructs a Plugin for a given Section.
type PluginFactory func(*Section) Plugin

var plugins = map[string]PluginFactory{}

// RegisterPlugin registers a PluginFactory under a name. Booklit sections can
// then use the plugin by calling \use-plugin with the same name.
//
// This is typically called by a plugin package's init() function.
func RegisterPlugin(name string, factory PluginFactory) {
	plugins[name] = factory
	logrus.WithField("plugin", name).Info("plugin registered")
}

// LookupPlugin looks up the given plugin factory.
func LookupPlugin(name string) (PluginFactory, bool) {
	plugin, found := plugins[name]
	return plugin, found
}
