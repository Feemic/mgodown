package mgodown

import (
	"bytes"
	"fmt"
	"io"
	"regexp"

	"github.com/Feemic/mgodown/ast"
	"github.com/Feemic/mgodown/html"
	"github.com/Feemic/mgodown/parser"
)

// Renderer is an interface for implementing custom renderers.
type Renderer interface {
	// RenderNode renders markdown node to w.
	// It's called once for a leaf node.
	// It's called twice for non-leaf nodes:
	// * first with entering=true
	// * then with entering=false
	//
	// Return value is a way to tell the calling walker to adjust its walk
	// pattern: e.g. it can terminate the traversal by returning Terminate. Or it
	// can ask the walker to skip a subtree of this node by returning SkipChildren.
	// The typical behavior is to return GoToNext, which asks for the usual
	// traversal to the next node.
	RenderNode(w io.Writer, node ast.Node, entering bool) ast.WalkStatus

	// RenderHeader is a method that allows the renderer to produce some
	// content preceding the main body of the output document. The header is
	// understood in the broad sense here. For example, the default HTML
	// renderer will write not only the HTML document preamble, but also the
	// table of contents if it was requested.
	//
	// The method will be passed an entire document tree, in case a particular
	// implementation needs to inspect it to produce output.
	//
	// The output should be written to the supplied writer w. If your
	// implementation has no header to write, supply an empty implementation.
	RenderHeader(w io.Writer, ast ast.Node)

	// RenderFooter is a symmetric counterpart of RenderHeader.
	RenderFooter(w io.Writer, ast ast.Node)
}

// Parse parsers a markdown document using provided parser. If parser is nil,
// we use parser configured with parser.CommonExtensions.
//
// It returns AST (abstract syntax tree) that can be converted to another
// format using Render function.
func Parse(markdown []byte, p *parser.Parser) ast.Node {
	if p == nil {
		p = parser.New()
	}
	markdown = Preparse(markdown)
	return p.Parse(markdown)
}

//feemic add:Pre process, hand line break etc.
func Preparse(markdown []byte) []byte{
	markdown = bytes.ReplaceAll(markdown,[]byte("\r\n"),[]byte("\n"))
	markdown = bytes.ReplaceAll(markdown,[]byte("\r"),[]byte("\n"))
	re, _ := regexp.Compile(`\[TOC\]\n?`)
	if ix := re.FindIndex(markdown); len(ix) > 0{
		markdown = markdown[ix[1]:]
	}
	return markdown
}
// Render uses renderer to convert parsed markdown document into a different format.
//
// To convert to HTML, pass html.Renderer
func Render(doc ast.Node, renderer Renderer) []byte {
	var buf bytes.Buffer
	renderer.RenderHeader(&buf, doc)
	ast.WalkFunc(doc, func(node ast.Node, entering bool) ast.WalkStatus {
		return renderer.RenderNode(&buf, node, entering)
	})
	renderer.RenderFooter(&buf, doc)
	return buf.Bytes()
}

func HeadingRender(renderer Renderer, doc ast.Node) []byte{
	buf := bytes.Buffer{}

	inHeading := false
	tocLevel := 0
	headingCount := 0
	minHeading := -1

	ast.WalkFunc(doc, func(node ast.Node, entering bool) ast.WalkStatus {
		if nodeData, ok := node.(*ast.Heading); ok && !nodeData.IsTitleblock {
			//fmt.Println(nodeData.Level, tocLevel,entering)
			inHeading = entering
			if !entering {
				buf.WriteString("</a>") //ready to go next, add content close
				return ast.GoToNext  //go to next, inheading => true
			}
			nodeData.HeadingID = fmt.Sprintf("toc_%d", headingCount)
			if minHeading == -1 || nodeData.Level <= minHeading{
				minHeading = nodeData.Level
			}
			if nodeData.Level == tocLevel {
				buf.WriteString("</li>\n\n<li>") //simple level close and start
			} else if nodeData.Level < tocLevel {
				for nodeData.Level < tocLevel {
					tocLevel--
					buf.WriteString("</li>\n</ul>")
				}
				buf.WriteString("</li>\n\n<li>")  //close and start
			} else {
				for nodeData.Level > tocLevel {
					tocLevel++
					buf.WriteString("\n<ul>\n<li>")  //has son, so start
				}
			}
			//fmt.Println(nodeData.Level, tocLevel)

			fmt.Fprintf(&buf, `<a href="#toc_%d">`, headingCount) //header content
			headingCount++
			return ast.GoToNext
		}

		if inHeading {
			return renderer.RenderNode(&buf, node, entering) //attach text
		}

		return ast.GoToNext
	})
	//final close
	for ; tocLevel > 0; tocLevel-- {
		buf.WriteString("</li>\n</ul>")
	}
	b := buf.Bytes()
	if minHeading >1 {
		leftCut := (minHeading-1)*10
		rightCut := (minHeading-1)*11
		b = b[leftCut:len(b)-rightCut]
		return b
	}
	return b
}

// ToHTML converts markdownDoc to HTML.
//
// You can optionally pass a parser and renderer. This allows to customize
// a parser, use a customized html render or use completely custom renderer.
//
// If you pass nil for both, we use parser configured with parser.CommonExtensions
// and html.Renderer configured with html.CommonFlags.
func ToHTML(markdown []byte, p *parser.Parser, renderer Renderer) ([]byte, []byte) {
	doc := Parse(markdown, p)

	if renderer == nil {
		opts := html.RendererOptions{
			Flags: html.CommonFlags,
		}
		renderer = html.NewRenderer(opts)
	}

	heading := HeaderToHtml(markdown, doc, renderer)
	return heading, Render(doc, renderer)
}

//feemic add:Just render header
func HeaderToHtml(markdown []byte, doc ast.Node, renderer Renderer) ([]byte) {
	parseHeading := false
	re, _ := regexp.Compile(`\[TOC\]\n?`)
	if ix := re.FindIndex(markdown); len(ix) > 0{
		parseHeading = true
	}
	var heading []byte
	if parseHeading{
		heading = HeadingRender(renderer, doc)
	}
	return heading
}