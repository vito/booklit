// Package stages implements the three stages for processing Booklit documents:
//
// The first stage is Evaluate, which interprets a node to generate content
// while calling plugin functions.
//
// The second stage is Collect, which traverses all content and sets tags
// defined by Targets.
//
// The third and final stage is Resolve, which traverses all content and
// resolves References to their tags.
package stages
