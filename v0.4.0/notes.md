## breaking changes

this release removes `booklit.Block` and `booklit.Element`. these two types of
`booklit.Content` were cop-outs that only served to allow adding CSS classes to
content that they wrapped. this leaked presentation/templating concerns into
content. you should use `booklit.Styled` instead.

the one use of `booklit.Block`, which was to force flow content to be block
content, is now localized to `booklit.Styled`, which was the only place it was
ever needed. there is a boolean field, `Block`, which can be set on
`booklit.Styled` to force it to be block content.
