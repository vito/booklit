## new features

* booklit can now be run with `--save-search-index`, which will generate a
  `.json` file alongside the generated output containing titles and content for
  each tag. this can be fed in to a client-side search engine to implement
  search functionality on the generated docs.

* added a few helper methods to `*booklit.Section` which can be useful when
  templating:

  * `IsOrHasChild(*Section)`: useful when rendering a navigation tree, to
    determine whether to show a section as 'expanded'

  * `Depth`: analogous to `PageDepth`, this is the absolute depth of the
    section in the tree. this is exposed in the search index, as it can be
    useful to 'weigh' toplevel content higher when ranking the results.

  * `SetTagAnchored`: registers an anchored tag with additional content that
    can be shown in search results

* arbitrary content can now be specified to `\target`, like so:

  ```
  \target{some-tag}{This is a title}{This is arbitrary content.}
  ```

  this is mainly useful when building up search indexes, as it allows you to
  provide content to show for the given tag.


## bug fixes

* the `chroma` plugin will no longer leave a `<p>` wrapping around code blocks
