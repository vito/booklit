all: ast/booklit.peg.go render/bindata.go

ast/booklit.peg.go: ast/booklit.peg
	pigeon -o ast/booklit.peg.go ast/booklit.peg

render/bindata.go: render/html render/html/*.tmpl
	go-bindata -o render/bindata.go -pkg render render/html/*.tmpl
