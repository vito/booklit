all: ast/booklit.peg.go pkged.go

ast/booklit.peg.go: ast/booklit.peg
	pigeon ast/booklit.peg | goimports > ast/booklit.peg.go

pkged.go: render/html/*.tmpl render/text/*.tmpl
	pkger
