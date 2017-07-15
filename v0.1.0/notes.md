this release packs quite a punch. in using booklit for https://concouse.ci, we
found a lot of room for improvement, and iterated quickly along the way to
where we are now.

there were also just a bunch of silly things missing. like comments.

so, all that seems to deserve at least a minor bump!

here are the changes grouped by section:

## content

* sections can now contain "partials", which are named bits of content that can
  be set by the document, accessed by the template, and stored or accessed by
  plugins.
  
  this introduces a much nicer way of developing websites with a lot of
  non-prose content; portions can be written in the document, taking advantage
  of booklit semantics, and plucked piece-by-piece into the template.

  plugins can also store data in partials, allowing them to communicate data to
  their own templates.

  in general, this will reduce the number of lines of code in plugins just to
  get things looking the right way, and instead allow templates, documents, and
  plugins to do what they individually do best.

* styled content via `booklit.Styled` will now treat the style as an artbirary
  string denoting a template to use. this was basically already how it worked,
  but only with the pre-declared set of styles.

  with this change, plugins can use whatever templates they want to render any
  content, greatly increasing flexibility.

* a section can declare `\single-page`, which will cause all `\split-sections`
  calls to no longer result in sections being split, and instead only reset the
  "page depth" (which is used for e.g. determining `<h1-h6>`. this is useful
  for having a single-page form of your otherwise many-paged documentation.

* added a `rawHTML` template function, which can be used for directly rendering
  trusted HTML content. this is useful for e.g. using a syntax highlighter to
  generate a bunch of styled HTML.

* added `\table` and `\definitions`

* added `\omit-children-from-table-of-contents`, which modifies the section
  to... omit its children from any table of contents.

* added `\aux{...}`, which is used to omit text from a title when referenced by
  a link or by the table of contents.

  note that any custom templates must be changed to call `stripAux` before
  passing the title to `render`.

* added `\aside{...}` for adding side-notes and such. by default, it is
  rendered as a `<blockquote class="aside">`.

* added `\inset{...}` for rendering a block of content indented a bit.

* added `\image{path/to/file.png}`.

* added `booklit.Element`, which is like `booklit.Block` but for inline
  (sentence) elements. it renders as a `<span>` with the given class by
  default.

  this is primarily useful for plugins.


## compiler

* added `--allow-broken-references`, which can be handy for quick local
  iteration where expectations are lower.

* fixed `render` calls in templates not using user-provided templates when
  recursing.

* fixed handling of `\split-sections` in sub-sections.


## syntax

* support for comments:

  ```haskell
  foo {- bar -}baz

  {-
    This is a block comment.

    {- This is a nested block comment. -}

    Bye!
  -}
  ```

* plugin methods with no arguments can now be spliced in to words like so:

  ```latex
  foo{\bar}baz
  ```

* sentences in arguments can now be wrapped onto the next line; a single space
  will be put in place of the linebreaks

  ```latex
  this is a really really really \italic{really really really
  long sentence} here!
  ```

* empty arguments like `\foo{bar}{}{baz}` are now supported, and act like empty
  strings
