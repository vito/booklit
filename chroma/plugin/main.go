package main

import (
	"github.com/vito/booklit"
)

func init() {
	booklit.RegisterPlugin("chroma", chroma.NewPlugin)
}
