from setuptools import setup, find_packages

setup (
  name='booklitlexer',
  packages=find_packages(),
  entry_points =
  """
  [pygments.lexers]
  booklitlexer = booklitlexer.lexer:BooklitLexer
  """,
)
