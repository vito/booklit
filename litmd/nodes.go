package litmd

import (
	"github.com/yuin/goldmark/ast"
)

var KindInvoke = ast.NewNodeKind("Invoke")
var KindInvokeInlineArgument = ast.NewNodeKind("InvokeInlineArgument")
var KindInvokeBlockArgument = ast.NewNodeKind("InvokeBlockArgument")

type Invoke struct {
	ast.BaseInline

	Function string
}

func (node *Invoke) Kind() ast.NodeKind {
	return KindInvoke
}

func (node *Invoke) Dump(source []byte, level int) {
	ast.DumpHelper(node, source, level, map[string]string{
		"Function": node.Function,
	}, nil)
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

type InvokeBlockArgument struct {
	ast.BaseBlock
}

func (node *InvokeBlockArgument) Kind() ast.NodeKind {
	return KindInvokeBlockArgument
}

func (node *InvokeBlockArgument) Dump(source []byte, level int) {
	ast.DumpHelper(node, source, level, map[string]string{}, nil)
}
