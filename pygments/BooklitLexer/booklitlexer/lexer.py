from pygments.lexer import RegexLexer
from pygments.token import *

class BooklitLexer(RegexLexer):
    name = 'Booklit'
    aliases = ['booklit']
    filenames = ['*.lit']

    tokens = {
        'root': [
            (r'[^\\{}]+', Text),
            (r'\{\{\{', String.Double, 'verbatim'),
            (r'\{-', Comment.Multiline, 'comment'),
            (r'[{}]', Name.Builtin),
            (r'\\([a-z-]+)', Keyword),
            (r'\\[\\{}]+', Text),
        ],
        'verbatim': [
            (r'\}\}\}', String.Double, '#pop'),
            (r'[^}]+', String.Double),
            (r'}[^\}]', String.Double),
        ],
        'comment': [
            (r'[^-{}]+', Comment.Multiline),
            (r'\{-', Comment.Multiline, '#push'),
            (r'-\}', Comment.Multiline, '#pop'),
            (r'[-{}]', Comment.Multiline),
        ],
    }
