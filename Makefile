all: ast/booklit.peg.go errhtml/normalize.css errhtml/logo.svg.base64 errhtml/bindata.go render/html/bindata.go render/text/bindata.go

ast/booklit.peg.go: ast/booklit.peg
	pigeon ast/booklit.peg | goimports > ast/booklit.peg.go

errhtml/normalize.css: docs/css/normalize.css
	cp $< $@

errhtml/logo.svg.base64: docs/css/images/booklit.svg
	base64 -w0 $< > $@

errhtml/bindata.go: errhtml errhtml/*.tmpl errhtml/*.css errhtml/*.base64
	go-bindata -o errhtml/bindata.go -pkg errhtml errhtml/*.tmpl errhtml/*.css errhtml/*.base64

render/html/bindata.go: render/html render/html/*.tmpl
	go-bindata -o render/html/bindata.go -pkg html render/html/*.tmpl

render/text/bindata.go: render/text render/text/*.tmpl
	go-bindata -o render/text/bindata.go -pkg text render/text/*.tmpl

clean:
	find . -name bindata.go -delete
	rm ast/booklit.peg.go
