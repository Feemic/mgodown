# A Markdown Parser and HTML or Text Renderer for Go

It's fast and supports common extensions.

## Installation

    go get -u github.com/feemic/mgodown

API Docs:

- /ast : defines abstract syntax tree of parsed markdown document
- /parser : parser
- /html : html renderer


## History

mgodown is a fork of https://github.com/gomarkdown/markdown which is a fork of v2 of https://github.com/russross/blackfriday


Blackfriday itself was based on C implementation [sundown](https://github.com/vmg/sundown) which in turn was based on [libsoldout](http://fossil.instinctive.eu/libsoldout/home).

## License

[Simplified BSD License](LICENSE.txt)
