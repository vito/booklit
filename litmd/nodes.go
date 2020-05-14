package litmd

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

var KindInvokeBlock = ast.NewNodeKind("InvokeBlock")
var KindInvokeInline = ast.NewNodeKind("InvokeInline")
var KindInvokeInlineArgument = ast.NewNodeKind("InvokeInlineArgument")
var KindInvokeBlockArgument = ast.NewNodeKind("InvokeBlockArgument")
var KindInvokeBlockArgumentPreformatted = ast.NewNodeKind("InvokeBlockArgumentPreformatted")
var KindInvokeBlockArgumentVerbatim = ast.NewNodeKind("InvokeBlockArgumentVerbatim")

type InvokeBlock struct {
	ast.BaseBlock

	Function string
}

func (node *InvokeBlock) Kind() ast.NodeKind {
	return KindInvokeBlock
}

func (node *InvokeBlock) Dump(source []byte, level int) {
	ast.DumpHelper(node, source, level, map[string]string{
		"Function": node.Function,
	}, nil)
}

type InvokeInline struct {
	ast.BaseInline

	Function string
}

func (node *InvokeInline) Kind() ast.NodeKind {
	return KindInvokeInline
}

func (node *InvokeInline) Dump(source []byte, level int) {
	ast.DumpHelper(node, source, level, map[string]string{
		"Function": node.Function,
	}, nil)
}

func (b *InvokeInline) HasBlankPreviousLines() bool {
	panic("can not call with inline nodes.")
}

// SetBlankPreviousLines implements Node.SetBlankPreviousLines.
func (b *InvokeInline) SetBlankPreviousLines(v bool) {
	panic("can not call with inline nodes.")
}

// Lines implements Node.Lines
func (b *InvokeInline) Lines() *text.Segments {
	panic("can not call with inline nodes.")
}

// SetLines implements Node.SetLines
func (b *InvokeInline) SetLines(v *text.Segments) {
	panic("can not call with inline nodes.")
}

type InvokeInlineArgument struct {
	ast.BaseInline
}

func (node *InvokeInlineArgument) Kind() ast.NodeKind {
	return KindInvokeInlineArgument
}

func (node *InvokeInlineArgument) Dump(source []byte, level int) {
	ast.DumpHelper(node, source, level, map[string]string{}, nil)
}

func (b *InvokeInlineArgument) HasBlankPreviousLines() bool {
	panic("can not call with inline nodes.")
}

// SetBlankPreviousLines implements Node.SetBlankPreviousLines.
func (b *InvokeInlineArgument) SetBlankPreviousLines(v bool) {
	panic("can not call with inline nodes.")
}

// Lines implements Node.Lines
func (b *InvokeInlineArgument) Lines() *text.Segments {
	panic("can not call with inline nodes.")
}

// SetLines implements Node.SetLines
func (b *InvokeInlineArgument) SetLines(v *text.Segments) {
	panic("can not call with inline nodes.")
}

type InvokeBlockArgument struct {
	ast.BaseBlock
}

func (node *InvokeBlockArgument) Kind() ast.NodeKind {
	return KindInvokeBlockArgument
}

func (node *InvokeBlockArgument) Dump(source []byte, level int) {
	ast.DumpHelper(node, source, level, map[string]string{}, nil)
}

