all: ast/booklit.peg.go render/bindata.go docs/index.html

docs/index.html: docs/index.lit
	go run cmd/booklit/*.go -i docs/index.lit -o docs

ast/booklit.peg.go: ast/booklit.peg
	pigeon -o ast/booklit.peg.go ast/booklit.peg

render/bindata.go: render/html/*.html
	go-bindata -o render/bindata.go -pkg render render/html/*.html
