targets=ast/booklit.peg.go docs/css/booklit.css errhtml/errors.css docs/outputs/index.html

all: $(targets)

ast/booklit.peg.go: ast/booklit.peg
	pigeon ast/booklit.peg | goimports > ast/booklit.peg.go

errhtml/errors.css: less/errors.less less/*.less
	yarn run lessc $< $@

docs/css/booklit.css: less/docs.less less/*.less
	yarn run lessc $< $@

less/logo-url.less: docs/css/images/booklit.svg
	yarn run build-logo-url-less

docs/outputs/index.html: docs/hello/*.lit docs/hello/html/*.tmpl docs/hello/go/*
	cd docs/hello && booklit --plugin github.com/vito/booklit/docs/hello/go --html-templates html -i index.lit -o ../outputs

clean:
	rm -f $(targets)
