package builtins

import (
	"fmt"

	"github.com/vito/booklit"
	"github.com/vito/booklit/ast"
)

func init() {
	Register("Table", tableFunc)
	Register("TableRow", tableRowFunc)
	Register("Row", tableRowFunc)
}

// tableFunc — `<Table><Row>...</Row></Table>`. Rows come from <Row> /
// <TableRow> children; non-row children are dropped.
func tableFunc(ctx *Context, _ map[string]ast.Node, children []ast.Node) (booklit.Content, error) {
	t := booklit.Table{}
	for _, child := range children {
		val, err := ctx.Evaluate(child)
		if err != nil {
			return nil, err
		}
		if val == nil {
			continue
		}
		row, ok := val.(rowContent)
		if !ok {
			continue
		}
		t.Rows = append(t.Rows, row.Cells)
	}
	return t, nil
}

// rowContent wraps a single row's cells. <Table> recognizes this type
// when collecting rows from child evaluations.
type rowContent struct {
	Cells []booklit.Content
}

func (r rowContent) IsFlow() bool                  { return false }
func (r rowContent) String() string                { return fmt.Sprintf("row(%d cells)", len(r.Cells)) }
func (r rowContent) Visit(v booklit.Visitor) error { return nil }

// tableRowFunc — `<Row><Cell>a</Cell><Cell>b</Cell></Row>` or just
// `<TableRow>...</TableRow>`. Cells come from <Item> / <Cell> children
// or, lacking those, the row's evaluated children become a single cell
// each at the AST top level.
func tableRowFunc(ctx *Context, _ map[string]ast.Node, children []ast.Node) (booklit.Content, error) {
	var cells []booklit.Content
	for _, child := range children {
		val, err := ctx.Evaluate(child)
		if err != nil {
			return nil, err
		}
		if val == nil {
			continue
		}
		// Both <Item> (carrying itemContent) and <Cell> are accepted.
		if it, ok := val.(itemContent); ok {
			cells = append(cells, it.Content)
			continue
		}
	}
	return rowContent{Cells: cells}, nil
}
