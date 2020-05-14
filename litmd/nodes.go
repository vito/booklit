package litmd

import "github.com/yuin/goldmark/ast"

var KindInvoke = ast.NewNodeKind("Invoke")
var KindInvokeArgument = ast.NewNodeKind("InvokeArgument")
var KindInvokeArgumentPreformatted = ast.NewNodeKind("InvokeArgumentPreformatted")
var KindInvokeArgumentVerbatim = ast.NewNodeKind("InvokeArgumentVerbatim")

type invokeNode struct {
	ast.BaseBlock

	Function string
}

func (node *invokeNode) Kind() ast.NodeKind {
	return KindInvoke
}

func (node *invokeNode) Dump(source []byte, level int) {
	ast.DumpHelper(node, source, level, map[string]string{
		"Function": node.Function,
	}, nil)
}

type invokeArgumentNode struct {
	ast.BaseBlock
}

func (node *invokeArgumentNode) Kind() ast.NodeKind {
	return KindInvokeArgument
}

func (node *invokeArgumentNode) Dump(source []byte, level int) {
	ast.DumpHelper(node, source, level, map[string]string{}, nil)
}

