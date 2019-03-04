## new features

* added a few more helper methods to `*booklit.Section` which can be useful
  when templating:

  * `NextSibling`: returns the next sibling section, if any

  * `Next`: returns the first sub-section, if any, or the `NextSibling`

  * `Prev`: returns the previous sibling section, if any, or the parent section

  these are mainly in service of generating 'next/prev' links at the end of a
  section in order to guide the reader along to the rest of the content, in
  situations where the nav is separate and would be hard to notice
