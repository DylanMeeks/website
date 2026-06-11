package main

import (
	"fmt"
	"io"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	mdhtml "github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

func WriteHighlightCSS(w io.Writer) error {
	if highlightStyle == nil {
		return fmt.Errorf("highlight style not initialized")
	}

	return htmlFormatter.WriteCSS(w, highlightStyle)
}

func registerGruberDarkerStyle() {
	styles.Register(chroma.MustNewStyle("gruber-darker", chroma.StyleEntries{
		chroma.Background: "noinherit #e4e4ef bg:#181818",
		chroma.Text:       "#e4e4ef",

		chroma.Keyword:            "bold #ffdd33",
		chroma.KeywordConstant:    "bold #ffdd33",
		chroma.KeywordDeclaration: "bold #ffdd33",
		chroma.KeywordNamespace:   "bold #ffdd33",
		chroma.KeywordPseudo:      "bold #ffdd33",
		chroma.KeywordReserved:    "bold #ffdd33",
		chroma.KeywordType:        "#95a99f",

		chroma.Name:                 "#e4e4ef",
		chroma.NameAttribute:        "#f4f4ff",
		chroma.NameBuiltin:          "#ffdd33",
		chroma.NameBuiltinPseudo:    "#ffdd33",
		chroma.NameClass:            "#95a99f",
		chroma.NameConstant:         "#95a99f",
		chroma.NameDecorator:        "#ffdd33",
		chroma.NameEntity:           "#95a99f",
		chroma.NameException:        "#f43841",
		chroma.NameFunction:         "#96a6c8",
		chroma.NameFunctionMagic:    "#96a6c8",
		chroma.NameLabel:            "#f43841",
		chroma.NameNamespace:        "#95a99f",
		chroma.NameOther:            "#e4e4ef",
		chroma.NameProperty:         "#f4f4ff",
		chroma.NameTag:              "#ffdd33",
		chroma.NameVariable:         "#f4f4ff",
		chroma.NameVariableClass:    "#f4f4ff",
		chroma.NameVariableGlobal:   "#f4f4ff",
		chroma.NameVariableInstance: "#f4f4ff",
		chroma.NameVariableMagic:    "#f4f4ff",

		chroma.Literal:                  "#e4e4ef",
		chroma.LiteralString:            "#73c936",
		chroma.LiteralStringAffix:       "#73c936",
		chroma.LiteralStringBacktick:    "#73c936",
		chroma.LiteralStringChar:        "#73c936",
		chroma.LiteralStringDelimiter:   "#73c936",
		chroma.LiteralStringDoc:         "#73c936",
		chroma.LiteralStringDouble:      "#73c936",
		chroma.LiteralStringEscape:      "#9e95c7",
		chroma.LiteralStringHeredoc:     "#73c936",
		chroma.LiteralStringInterpol:    "#9e95c7",
		chroma.LiteralStringOther:       "#73c936",
		chroma.LiteralStringRegex:       "#73c936",
		chroma.LiteralStringSingle:      "#73c936",
		chroma.LiteralStringSymbol:      "#9e95c7",
		chroma.LiteralNumber:            "#9e95c7",
		chroma.LiteralNumberBin:         "#9e95c7",
		chroma.LiteralNumberFloat:       "#9e95c7",
		chroma.LiteralNumberHex:         "#9e95c7",
		chroma.LiteralNumberInteger:     "#9e95c7",
		chroma.LiteralNumberIntegerLong: "#9e95c7",
		chroma.LiteralNumberOct:         "#9e95c7",

		chroma.Operator:     "#ffdd33",
		chroma.OperatorWord: "bold #ffdd33",
		chroma.Punctuation:  "#e4e4ef",

		chroma.Comment:            "#cc8c3c",
		chroma.CommentHashbang:    "#cc8c3c",
		chroma.CommentMultiline:   "#cc8c3c",
		chroma.CommentSingle:      "#cc8c3c",
		chroma.CommentSpecial:     "#cc8c3c",
		chroma.CommentPreproc:     "#95a99f",
		chroma.CommentPreprocFile: "#95a99f",

		chroma.Generic:           "#e4e4ef",
		chroma.GenericDeleted:    "#ff4f58",
		chroma.GenericEmph:       "italic #e4e4ef",
		chroma.GenericError:      "#f43841",
		chroma.GenericHeading:    "bold #96a6c8",
		chroma.GenericInserted:   "#73c936",
		chroma.GenericOutput:     "#52494e",
		chroma.GenericPrompt:     "#ffdd33",
		chroma.GenericStrong:     "bold #e4e4ef",
		chroma.GenericSubheading: "bold #96a6c8",
		chroma.GenericTraceback:  "#f43841",

		chroma.Error: "#f43841 bg:#181818",
	}))
}

var (
	htmlFormatter  *html.Formatter
	highlightStyle *chroma.Style
)

func InitRenderer() {
	registerGruberDarkerStyle()

	htmlFormatter = html.New(html.WithClasses(true), html.TabWidth(2))
	if htmlFormatter == nil {
		panic("couldn't create html formatter")
	}

	styleName := "gruber-darker"
	highlightStyle = styles.Get(styleName)
	if highlightStyle == nil {
		panic(fmt.Sprintf("didn't find style '%s'", styleName))
	}
}

// based on https://github.com/alecthomas/chroma/blob/master/quick/quick.go
func htmlHighlight(w io.Writer, source, lang, defaultLang string) error {
	if lang == "" {
		lang = defaultLang
	}
	l := lexers.Get(lang)
	if l == nil {
		l = lexers.Analyse(source)
	}
	if l == nil {
		l = lexers.Fallback
	}
	l = chroma.Coalesce(l)

	it, err := l.Tokenise(nil, source)
	if err != nil {
		return err
	}
	return htmlFormatter.Format(w, highlightStyle, it)
}

func renderCode(w io.Writer, codeBlock *ast.CodeBlock, entering bool) {
	defaultLang := ""
	lang := string(codeBlock.Info)
	htmlHighlight(w, string(codeBlock.Literal), lang, defaultLang)
}

func myRenderHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	if code, ok := node.(*ast.CodeBlock); ok {
		renderCode(w, code, entering)
		return ast.GoToNext, true
	}
	return ast.GoToNext, false
}

func newCustomizedRender() *mdhtml.Renderer {
	opts := mdhtml.RendererOptions{
		Flags:          mdhtml.CommonFlags,
		RenderNodeHook: myRenderHook,
	}
	return mdhtml.NewRenderer(opts)
}

func MdToHTML(md []byte) []byte {

	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	renderer := newCustomizedRender()
	html := markdown.Render(doc, renderer)

	return html
}
