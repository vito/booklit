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

docs/outputs/index.html: docs/lit/*.md docs/html/*.tmpl docs/html/*.md treehighlight/*.go treehighlight/internal/tree_sitter_booklit/*.go treehighlight/internal/tree_sitter_booklit/src/parser.c dagger/booklitdoc/main.go dagger/booklitdoc/cmd/lit-syntax/main.go
	go run ./cmd/booklit -i docs/lit/index.md -o docs/outputs --html-templates docs/html

clean:
	rm -f $(targets)
