##  breaking changes

* the plugin system has been reverted back to the old reexec approach from pre-v0.8.0.

  the Go `plugin` package doesn't look like it will be getting Windows support any time soon, and the cross-compiled Darwin binaries stopped working, so really this just wasn't worth it.

* a few methods on `*booklit.Section` now take an `ast.Location` argument, which is required for the "better error messages" feature below. plugins should typically adapt by passing it the value from the newly added `.InvokeLocation` field on the section object given to the plugin.

* `booklit.Reference`, `booklit.Target`, and `booklit.Section` now have `Location` fields which should similarly be populated by plugins which construct them, again in service of better error messages.

## new features

errors are actually kinda helpful now!

the work here is multi-faceted:

* error messages have been reworded to be a bit less cryptic. (I think. Let me know.)

* the parser now preserves position information (i.e. line numbers), allowing it to be threaded through Booklit and returned in the errors.

* the CLI will render a helpful message, showing a snippet of the file where the error occurred, and even highlighting the exact region.

  ![booklit-err-cli](https://user-images.githubusercontent.com/1880/81368807-bbf80100-90be-11ea-9275-8e77923320c6.png)

* the web UI will also render a helpful message, but in HTML (*whoaaaaaaa*).

  ![booklit-err-web](https://user-images.githubusercontent.com/1880/81368808-bbf80100-90be-11ea-9088-656654d69dfe.png)

* a `PrettyError` interface has been introduced, allowing plugins to leverage these pretty-error mechanisms.

## bug fixes

* previously, `\include-section` in an inline `\section` would resolve its path relative to wherever `booklit` was run from. now it resolves paths relative to its outer section's file path, consistent with its use in a top-level document.

* `booklit --help` now exits successfully.

* the parser now accepts Windows-style `CRLF` line endings.

* the parser no longer requires a trailing linebreak at the end of the file.

* the parser no longer chokes on arguments containing a single linebreak:

  ```
  \foo{
  }
  ```

## misc

* anchor tags now use the `id` attribute rather than `name` (thanks @jtarchie!)
