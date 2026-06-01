module.exports = grammar({
  name: 'booklit',

  extras: $ => [],

  rules: {
    document: $ => repeat(choice(
      $.comment,
      $.command,
      $.tag,
      $.heading,
      $.code_span,
      $.delimiter,
      $.text,
    )),

    command: $ => seq(
      field('marker', $.backslash),
      field('name', $.identifier),
    ),

    backslash: _ => '\\',

    identifier: _ => /[A-Za-z][A-Za-z0-9-]*/,

    tag: $ => seq(
      '<',
      optional('/'),
      field('name', $.identifier),
      repeat(choice($.tag_string, $.tag_expr, $.tag_text)),
      optional('/'),
      '>',
    ),

    tag_string: _ => token(seq('"', repeat(choice(/[^"\\]+/, /\\./)), '"')),

    tag_expr: $ => seq(
      '{',
      repeat(choice($.command, $.tag_string, $.tag_expr, $.tag_text, $.text)),
      '}',
    ),

    tag_text: _ => token(/[^<>"{}]+/),

    heading: _ => token(seq(repeat(choice(' ', '\t')), repeat1('#'), /[^\n]*/)),

    code_span: _ => token(seq('`', repeat(choice(/[^`\\]+/, /\\./)), '`')),

    delimiter: _ => choice('{{{', '}}}', '{{', '}}', '{', '}'),

    comment: _ => token(seq('{-', repeat(choice(/[^-]+/, /-+[^}]/)), '-}')),

    text: _ => token(choice(/[^\\{}"`<#\n]+/, /[\n\\{}"`<#]/)),
  },
});
