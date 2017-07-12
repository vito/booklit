all: ast/booklit.peg.go render/bindata.go docs/index.html

docs/index.html: docs/lit/*.lit docs/lit/html/*.tmpl
	go run cmd/booklit/*.go -i docs/lit/index.lit -o docs --html-templates docs/lit/html --plugin github.com/vito/booklit/booklitdoc

ast/booklit.peg.go: ast/booklit.peg
	pigeon -o ast/booklit.peg.go ast/booklit.peg

render/bindata.go: render/html/*.tmpl
	go-bindata -o render/bindata.go -pkg render render/html/*.tmpl
