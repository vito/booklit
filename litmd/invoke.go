package litmd

import (
	"encoding/json"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type invokeParser struct {
}

func NewInlineInvokeParser() parser.InlineParser {
	return &invokeParser{}
}

func NewBlockArgParser() parser.BlockParser {
	return &blockArgParser{}
}

func (b *invokeParser) Trigger() []byte {
	return []byte{'\\', '{', '}'}
}

var funcRegexp = regexp.MustCompile(`\\([a-z-]+)`)

var inlineArgStateKey = parser.NewContextKey()

type inlineArgState struct {
	ast.BaseBlock

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

func (b *invokeParser) Parse(parent ast.Node, reader text.Reader, pc parser.Context) ast.Node {
	line, segment := reader.PeekLine()
	log.Println("PARSE", string(line))
	switch line[0] {
	// starting an inline invoke call
	case '\\':
		matches := reader.FindSubMatch(funcRegexp)
		if matches == nil {
			return nil
		}

		function := string(matches[1])

		return &invokeNode{
			Function: function,
		}

	// beginning an inline argument?
	case '{':
		log.Println("OPEN UP", string(line))
		if string(line) == "{\n" {
			return nil
		}

		state := &inlineArgState{
			Segment: text.NewSegment(segment.Start, segment.Start+1),
		}

		pushInlineArgState(pc, state)

		reader.Advance(1)
		return state

	// ending an inline argument
	case '}':
		log.Printf("ENDING INLINE ARG (parent %T)", parent)

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

		arg := &invokeArgumentNode{}
		processInlineArg(reader.Source(), parent, arg, last, pc)

		last.Dump(reader.Source(), 0)

		last.Parent().RemoveChild(last.Parent(), last)
		return arg

	default:
		panic("impossible?")
	}
}

func processInlineArg(source []byte, parent ast.Node, arg *invokeArgumentNode, last *inlineArgState, pc parser.Context) {
	for c := last.NextSibling(); c != nil; {
		log.Printf("SIB")
		c.Dump(source, 1)
		next := c.NextSibling()
		parent.RemoveChild(parent, c)
		arg.AppendChild(arg, c)
		c = next
	}
}

type blockArgParser struct {
	ast.BaseBlock
}

func (b *blockArgParser) Trigger() []byte {
	return []byte{'{', '}'}
}
func (b *blockArgParser) Open(parent ast.Node, reader text.Reader, pc parser.Context) (ast.Node, parser.State) {
	log.Printf("OPEN %T", parent)
	// nothing to de

	line, _ := reader.PeekLine()
	log.Println("BLOCK ARG? (must be {)", string(line))
	switch string(line) {
	case "{":
		return &invokeArgumentNode{}, parser.HasChildren
	case "}":
		return nil, parser.Close
	default:
		println("NOT A BLOCK ARG")
		return nil, parser.NoChildren
	}

}

func (b *blockArgParser) Continue(node ast.Node, reader text.Reader, pc parser.Context) parser.State {
	line, _ := reader.PeekLine()
	log.Printf("CONTINUE? (peek %q, node %T)", string(line), node)
	if string(line) == "}" {
		println("EXITING BLOCK ARG")
		return parser.Close
	}
	return parser.Continue | parser.HasChildren
}

func (b *blockArgParser) Close(node ast.Node, reader text.Reader, pc parser.Context) {
	log.Printf("CLOSE %T", node)
	// nothing to do
}

func (b *blockArgParser) CanInterruptParagraph() bool {
	return true
}

func (b *blockArgParser) CanAcceptIndentedLine() bool {
	return false
}
