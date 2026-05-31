// Command booklit-docs is the booklit binary used to build Booklit's own
// documentation site. It bundles the docs-specific built-ins
// (booklitdoc) on top of the regular booklit command.
package main

import (
	"github.com/vito/booklit/booklitcmd"

	_ "github.com/vito/booklit/docs/booklitdoc"
)

func main() {
	booklitcmd.Main()
}
