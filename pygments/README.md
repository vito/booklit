A plugin providing syntax highlighting via `\syntax{language}{code...}`. Shells
out to `pygmentize` to perform the actual highlighting.

A highlighter for Booklit code is available under `./BooklitLexer`. To install
it, run:

```bash
cd BooklitLexer
python setup.py install
```

...which may require `sudo`.
