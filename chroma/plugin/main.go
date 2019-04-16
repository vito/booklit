package main

import (
	"github.com/vito/booklit"
	"github.com/vito/booklit/chroma"
)

func init() {
	booklit.RegisterPlugin("chroma", chroma.NewPlugin)
}
