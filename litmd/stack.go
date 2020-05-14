package litmd

import (
	"fmt"

	bast "github.com/vito/booklit/ast"
)

type stack struct {
	seqs []bast.Sequence
}

func (stack *stack) push() {
	stack.seqs = append(stack.seqs, bast.Sequence{})
}

func (stack *stack) pop() bast.Sequence {
	end := stack.seqs[stack.last()]
	stack.seqs = stack.seqs[0:stack.last()]
	return end
}

func (stack *stack) append(node bast.Node) {
	end := stack.seqs[stack.last()]
	end = append(end, node)
	stack.seqs[stack.last()] = end
}

func (stack *stack) last() int {
	return len(stack.seqs) - 1
}

func (stack *stack) invoke(fun string, entering bool) {
	if entering {
		stack.push()
	} else {
		stack.append(bast.Invoke{
			Function:  fun,
			Arguments: stack.pop(),
		})
	}
}

func (stack *stack) dump() {
	for _, seq := range stack.seqs {
		fmt.Println(seq)
	}
}

