package booklitcmd

import (
	"fmt"
	"os"

	flags "github.com/jessevdk/go-flags"
	"github.com/vito/booklit"
)

func Main() {
	cmd := &Command{}
	cmd.Version = func() {
		fmt.Println(booklit.Version)
		os.Exit(0)
	}

	parser := flags.NewParser(cmd, flags.Default)
	parser.NamespaceDelimiter = "-"

	args, err := parser.Parse()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = cmd.Execute(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
