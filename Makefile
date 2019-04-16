all: ast/booklit.peg.go render/html/bindata.go render/text/bindata.go

ast/booklit.peg.go: ast/booklit.peg
	pigeon -o ast/booklit.peg.go ast/booklit.peg

render/html/bindata.go: render/html render/html/*.tmpl
	go-bindata -o render/html/bindata.go -pkg html render/html/*.tmpl

render/text/bindata.go: render/text render/text/*.tmpl
	go-bindata -o render/text/bindata.go -pkg text render/text/*.tmpl
