# A Markdown Parser and HTML or Text Renderer for Go

It's fast and supports common extensions.

## Installation

```
go get -u github.com/Feemic/mgodown
```

API Docs:

- /ast : defines abstract syntax tree of parsed markdown document
- /parser : parser
- /html : html renderer

## Usage

To convert markdown text to HTML (split into heading and main content):
```
md := []byte("[TOC]\r\n## markdown document")

extensions := parser.CommonExtensions
mkparser := parser.NewWithExtensions(extensions)
htmlFlags := html.CommonFlags | html.TOC
opts := html.RendererOptions{Flags: htmlFlags}
renderer := html.NewRenderer(opts)

heading, content := mgodown.ToHTML(md, mkparser, renderer)
```

## Extensions & Difference
Default option:no empty line before block are supported  
Tables: tables are supported if there are spaces in the table header  
MathJaX:  
add display math block  

> 
```math
n \geq 2^{\frac {h} {2}} - 1
```

## History

mgodown is a fork of [gomarkdown](https://github.com/gomarkdown/markdown) which is a fork of v2 of [blackfriday](https://github.com/russross/blackfriday)


Blackfriday itself was based on C implementation [sundown](https://github.com/vmg/sundown) which in turn was based on [libsoldout](http://fossil.instinctive.eu/libsoldout/home).

## License

[Simplified BSD License](LICENSE.txt)
