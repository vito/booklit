package litmd

// XXX: re-add delimiter stuff (linkBottom and etc)

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

var funcRegexp = regexp.MustCompile(`\\([a-z-]+)`)

type invokeInlineParser struct {
}

func NewInvokeInlineParser() parser.InlineParser {
	return &invokeInlineParser{}
}

func NewInvokeBlockParser(recurse parser.Parser) parser.BlockParser {
	return &invokeBlockParser{
		recurse: recurse,
	}
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
	log.Println("PARSE", string(line))
	switch line[0] {
	// starting an inline invoke call
	case '\\':
		matches := reader.FindSubMatch(funcRegexp)
		if matches == nil {
			return nil
		}

		function := string(matches[1])

		return &InvokeInline{
			Function: function,
		}

	// beginning an inline argument?
	case '{':
		if string(line) == "{\n" {
			log.Println("NOT AN INLINE ARGUMENT")
			return nil
		}

		log.Printf("????????????????????????? IM AN INLINE ARG STATE\n")
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

		arg := &InvokeInlineArgument{}
		processInlineArg(reader.Source(), parent, arg, last, pc)

		last.Dump(reader.Source(), 0)

		last.Parent().RemoveChild(last.Parent(), last)
		return arg

	default:
		panic("impossible invoke inline open: " + string(line))
	}
}

func processInlineArg(source []byte, parent ast.Node, arg *InvokeInlineArgument, last *inlineArgState, pc parser.Context) {
	for c := last.NextSibling(); c != nil; {
		log.Printf("SIB")
		c.Dump(source, 1)
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
	return []byte{'\\', '{', '}'}
}
func (b *invokeBlockParser) Open(parent ast.Node, reader text.Reader, pc parser.Context) (ast.Node, parser.State) {
	log.Printf("OPEN %T", parent)
	// nothing to de

	// startPos, startSeg := reader.Position()

	line, _ := reader.PeekLine()
	switch line[0] {
	case '\\':
		matches := reader.FindSubMatch(funcRegexp)
		if matches == nil {
			return nil, parser.NoChildren
		}

		function := string(matches[1])

		invoke := &InvokeBlock{
			Function: function,
		}

		// inCurly := 0

		// 	beforeArgsPos, beforeArgsSeg := reader.Position()
		// 	plb, _ := reader.PeekLine()
		// 	log.Printf("PEEK BEFORE SCANNING ARGS: %q\n", string(plb))

		// 	var args []byte
		// scanInlineArgs:
		// 	for {
		// 		pl, _ := reader.PeekLine()
		// 		if len(pl) == 0 {
		// 			// non-block invoke at EOF
		// 			reader.SetPosition(startPos, startSeg)
		// 			return nil, parser.NoChildren
		// 		}

		// 		switch pl[0] {
		// 		case '{':
		// 			if pl[1] == '\n' {
		// 				break scanInlineArgs
		// 			}

		// 			reader.Advance(1)
		// 			inCurly++
		// 		case '}':
		// 			if inCurly == 0 {
		// 				panic("ROGUE CLOSE CURLY")
		// 				reader.SetPosition(startPos, startSeg)
		// 				return nil, parser.NoChildren
		// 			}

		// 			reader.Advance(1)
		// 			inCurly--
		// 		default:
		// 			if inCurly == 0 {
		// 				panic("HIT CHARACTER WHILE NOT IN CURLY")
		// 				// \foo{bar} foo { is not a block
		// 				reader.SetPosition(startPos, startSeg)
		// 				return nil, parser.NoChildren
		// 			}

		// 			log.Printf("ADVANCING PAST %q\n", pl[0])

		// 			reader.Advance(1)
		// 		}

		// 		args = append(args, pl[0])
		// 	}

		// 	afterArgsPos, afterArgsSeg := reader.Position()

		// 	log.Printf("COLLECTED ARGS: %s\n", string(args))

		// 	if len(args) > 0 {
		// reader.SetPosition(beforeArgsPos, beforeArgsSeg)
		pl, _ := reader.PeekLine()
		log.Printf("PEEK BEFORE RECURSE: %q\n", string(pl))
		inlineArgs := b.recurse.Parse(reader)
		if inlineArgs != nil {
			para := inlineArgs.FirstChild()
			if para != nil {
				log.Println("vvvvvvvvvvvvvvvvvvvvv PARSED PARA OF ARGS vvvvvvvvvvvvvvvvvvvvv")
				para.Dump(reader.Source(), 0)
				// for n := para.FirstChild(); n != nil; n = n.NextSibling() {
				// 	log.Printf("APPEND CHILD: %#v\n", n)
				// 	argPara := ast.NewParagraph()
				// 	argPara.AppendChild(argPara, n)

				// 	invoke.AppendChild(invoke, argPara)
				// }
				// invoke.AppendChild(invoke, para)
			}
		}
		// }

		// reader.SetPosition(afterArgsPos, afterArgsSeg)

		if reader.Peek() == '{' {
			return invoke, parser.HasChildren
		} else {
			return invoke, parser.NoChildren
		}

		// for reader.Peek() == '{' {
		// 	// reader.Advance(1)

		// 	arg := b.recurse.Parse(reader, parser.WithContext(pc))
		// 	log.Println("RECURSED")
		// 	if arg == nil {
		// 		log.Println("GOT NIL")
		// 		break
		// 	}

		// 	doc := arg.(*ast.Document)

		// 	para := doc.FirstChild()
		// 	if para == nil {
		// 		log.Printf("NO PARA!")
		// 		break
		// 	}

		// 	log.Printf("CHILD: %T\n", para)
		// 	arg.Dump(reader.Source(), 1)

		// 	invoke.AppendChild(invoke, para)
		// }

		// confirmed block
		if reader.Peek() == '{' {
			log.Println("~~~~~~~~~~~~~~~~~~~ HAS BLOCK! ~~~~~~~~~~~~~~~~~~~~~~~~~")
			return invoke, parser.HasChildren
		} else {
			return invoke, parser.NoChildren
		}
	case '{':
		log.Printf("BLOCK ARGUMENT??? %q\n", string(line))
		if string(line) != "{\n" {
			log.Println("NOT A BLOCK ARGUMENT")
			return nil, parser.NoChildren
		}

		reader.Advance(1)
		log.Printf("TOTALLY A BLOCK ARGUMENT: %T\n", parent)
		return &InvokeBlockArgument{}, parser.HasChildren
	case '}':
		return nil, parser.Close
	default:
		panic("impossible invoke block open: " + string(line))
	}
}

func (b *invokeBlockParser) Continue(node ast.Node, reader text.Reader, pc parser.Context) parser.State {
	line, _ := reader.PeekLine()
	log.Printf("CONTINUE? (peek %q, node %T)", string(line), node)
	if string(line) == "}" {
		reader.Advance(1)
		println("EXITING BLOCK ARG")
		return parser.Close
	}
	return parser.Continue | parser.HasChildren
}

func (b *invokeBlockParser) Close(node ast.Node, reader text.Reader, pc parser.Context) {
	log.Printf("CLOSE %T", node)
	// nothing to do
}

func (b *invokeBlockParser) CanInterruptParagraph() bool {
	// XXX: confirm, this is a guess
	return false
}

func (b *invokeBlockParser) CanAcceptIndentedLine() bool {
	// XXX: confirm, this is a guess
	return true
}
