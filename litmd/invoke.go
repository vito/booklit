package litmd

// XXX: re-add delimiter stuff (linkBottom and etc)

import (
	"encoding/json"
	"os"
	"regexp"
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

var funcRegexp = regexp.MustCompile(`\\([a-z-]+)`)

type invokeInlineParser struct {
}

func NewInvokeInlineParser() parser.InlineParser {
	return &invokeInlineParser{}
}

func NewInvokeBlockParser() parser.BlockParser {
	return &invokeBlockParser{}
}

func (b *invokeInlineParser) Trigger() []byte {
	return []byte{'\\', '{', '}'}
}

var inlineArgStateKey = parser.NewContextKey()

type inlineArgState struct {
	// XXX: this was BaseBlock at one point
	ast.BaseInline

	Segment text.Segment

	Prev  *inlineArgState
	Next  *inlineArgState
	First *inlineArgState
	Last  *inlineArgState
}

func (s *inlineArgState) Text(source []byte) []byte {
	return s.Segment.Value(source)
}

func (s *inlineArgState) Dump(source []byte, level int) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent(strings.Repeat("    ", level), "    ")
	enc.Encode(s)
}

var kindLinkLabelState = ast.NewNodeKind("InlineArgState")

func (s *inlineArgState) Kind() ast.NodeKind {
	return kindLinkLabelState
}

func pushInlineArgState(pc parser.Context, v *inlineArgState) {
	tlist := pc.Get(inlineArgStateKey)
	var list *inlineArgState
	if tlist == nil {
		list = v
		v.First = v
		v.Last = v
		pc.Set(inlineArgStateKey, list)
	} else {
		list = tlist.(*inlineArgState)
		l := list.Last
		list.Last = v
		l.Next = v
		v.Prev = l
	}
}

func removeInlineArgState(pc parser.Context, d *inlineArgState) {
	tlist := pc.Get(inlineArgStateKey)
	var list *inlineArgState
	if tlist == nil {
		return
	}
	list = tlist.(*inlineArgState)

	if d.Prev == nil {
		list = d.Next
		if list != nil {
			list.First = d
			list.Last = d.Last
			list.Prev = nil
			pc.Set(inlineArgStateKey, list)
		} else {
			pc.Set(inlineArgStateKey, nil)
		}
	} else {
		d.Prev.Next = d.Next
		if d.Next != nil {
			d.Next.Prev = d.Prev
		}
	}
	if list != nil && d.Next == nil {
		list.Last = d.Prev
	}
	d.Next = nil
	d.Prev = nil
	d.First = nil
	d.Last = nil
}

func (b *invokeInlineParser) Parse(parent ast.Node, reader text.Reader, pc parser.Context) ast.Node {
	line, segment := reader.PeekLine()

	switch line[0] {
	// starting an inline invoke call
	case '\\':
		matches := reader.FindSubMatch(funcRegexp)
		if matches == nil {
			return nil
		}

		function := string(matches[1])

		return &Invoke{
			Function: function,
		}

	// beginning an inline argument?
	case '{':
		state := &inlineArgState{
			Segment: text.NewSegment(segment.Start, segment.Start+1),
		}

		pushInlineArgState(pc, state)

		reader.Advance(1)
		return state

	// ending an inline argument
	case '}':
		tlist := pc.Get(inlineArgStateKey)
		if tlist == nil {
			return nil
		}
		last := tlist.(*inlineArgState).Last
		if last == nil {
			return nil
		}
		reader.Advance(1)
		removeInlineArgState(pc, last)

		arg := &InvokeInlineArgument{}
		processInlineArg(reader.Source(), parent, arg, last, pc)

		last.Parent().RemoveChild(last.Parent(), last)
		return arg

	default:
		panic("impossible invoke inline open: " + string(line))
	}
}

func processInlineArg(source []byte, parent ast.Node, arg *InvokeInlineArgument, last *inlineArgState, pc parser.Context) {
	for c := last.NextSibling(); c != nil; {
		next := c.NextSibling()
		parent.RemoveChild(parent, c)
		arg.AppendChild(arg, c)
		c = next
	}
}

type invokeBlockParser struct {
	ast.BaseBlock

	recurse parser.Parser
}

func (b *invokeBlockParser) Trigger() []byte {
	return []byte{'{', '}'}
}
func (b *invokeBlockParser) Open(parent ast.Node, reader text.Reader, pc parser.Context) (ast.Node, parser.State) {
	line, _ := reader.PeekLine()
	switch line[0] {
	case '{':
		reader.Advance(1)
		return &InvokeBlockArgument{}, parser.HasChildren
	case '}':
		panic("TODO")
		return nil, parser.Close
	default:
		panic("impossible invoke block open: " + string(line))
	}
}

func (b *invokeBlockParser) Continue(node ast.Node, reader text.Reader, pc parser.Context) parser.State {
	if reader.Peek() == '}' {
		reader.Advance(1)
		return parser.Close
	}

	return parser.Continue | parser.HasChildren
}

func (b *invokeBlockParser) Close(node ast.Node, reader text.Reader, pc parser.Context) {
	// nothing to do
}

// CanInterruptParagraph is `true` so that an opening '{' on a line following
// the invoke call can be parsed as a block argument:
//
//     \foo{bar}
//     {
//       baz
//     }
//
// Without this set to true, the second line continues to be parsed as part of
// the initial paragraph.
func (b *invokeBlockParser) CanInterruptParagraph() bool {
	return true
}

func (b *invokeBlockParser) CanAcceptIndentedLine() bool {
	// XXX: confirm, this is a guess
	return true
}
