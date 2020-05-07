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
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			fmt.Println(err)
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	err = cmd.Execute(args)
	if err != nil {
		if prettyErr, ok := err.(booklit.PrettyError); ok {
			prettyErr.PrettyPrint(os.Stderr)
		} else {
			fmt.Fprintln(os.Stderr, err)
		}

		os.Exit(1)
	}
}
