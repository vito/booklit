package booklit

// Plugin is an arbitrary object which is initialized with the Section that is
// using the plugin.
//
// Methods on the plugin object will be invoked during the "evaluation" stage
// by function calling syntax in Booklit documents.
type Plugin any

// PluginFactory constructs a Plugin for a given Section.
type PluginFactory func(*Section) Plugin
