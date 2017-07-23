package ast

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

func ifaceStr(val interface{}) string {
	str := ""
	for _, seg := range val.([]interface{}) {
		str = str + string(seg.([]byte))
	}

	return str
}

func ifaceSequences(val interface{}) []Sequence {
	seq := []Sequence{}
	for _, node := range val.([]interface{}) {
		seq = append(seq, node.(Sequence))
	}

	return seq
}

func ifaceNodes(val interface{}) []Node {
	nodes := []Node{}
	for _, node := range val.([]interface{}) {
		nodes = append(nodes, node.(Node))
	}

	return nodes
}

var g = &grammar{
	rules: []*rule{
		{
			name: "Booklit",
			pos:  position{line: 34, col: 1, offset: 581},
			expr: &actionExpr{
				pos: position{line: 34, col: 12, offset: 592},
				run: (*parser).callonBooklit1,
				expr: &seqExpr{
					pos: position{line: 34, col: 12, offset: 592},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 34, col: 12, offset: 592},
							label: "node",
							expr: &ruleRefExpr{
								pos:  position{line: 34, col: 17, offset: 597},
								name: "Paragraphs",
							},
						},
						&notExpr{
							pos: position{line: 34, col: 28, offset: 608},
							expr: &anyMatcher{
								line: 34, col: 29, offset: 609,
							},
						},
					},
				},
			},
		},
		{
			name: "Paragraphs",
			pos:  position{line: 38, col: 1, offset: 635},
			expr: &actionExpr{
				pos: position{line: 38, col: 15, offset: 649},
				run: (*parser).callonParagraphs1,
				expr: &seqExpr{
					pos: position{line: 38, col: 15, offset: 649},
					exprs: []interface{}{
						&zeroOrMoreExpr{
							pos: position{line: 38, col: 15, offset: 649},
							expr: &ruleRefExpr{
								pos:  position{line: 38, col: 15, offset: 649},
								name: "CommentSpacing",
							},
						},
						&labeledExpr{
							pos:   position{line: 38, col: 31, offset: 665},
							label: "paragraphs",
							expr: &oneOrMoreExpr{
								pos: position{line: 38, col: 42, offset: 676},
								expr: &actionExpr{
									pos: position{line: 38, col: 43, offset: 677},
									run: (*parser).callonParagraphs7,
									expr: &seqExpr{
										pos: position{line: 38, col: 43, offset: 677},
										exprs: []interface{}{
											&labeledExpr{
												pos:   position{line: 38, col: 43, offset: 677},
												label: "p",
												expr: &ruleRefExpr{
													pos:  position{line: 38, col: 45, offset: 679},
													name: "Paragraph",
												},
											},
											&zeroOrMoreExpr{
												pos: position{line: 38, col: 55, offset: 689},
												expr: &ruleRefExpr{
													pos:  position{line: 38, col: 55, offset: 689},
													name: "CommentSpacing",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "CommentSpacing",
			pos:  position{line: 42, col: 1, offset: 777},
			expr: &choiceExpr{
				pos: position{line: 42, col: 19, offset: 795},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 42, col: 19, offset: 795},
						val:        "\n",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 42, col: 26, offset: 802},
						name: "Comment",
					},
				},
			},
		},
		{
			name: "Paragraph",
			pos:  position{line: 44, col: 1, offset: 811},
			expr: &actionExpr{
				pos: position{line: 44, col: 14, offset: 824},
				run: (*parser).callonParagraph1,
				expr: &labeledExpr{
					pos:   position{line: 44, col: 14, offset: 824},
					label: "lines",
					expr: &oneOrMoreExpr{
						pos: position{line: 44, col: 20, offset: 830},
						expr: &actionExpr{
							pos: position{line: 44, col: 21, offset: 831},
							run: (*parser).callonParagraph4,
							expr: &seqExpr{
								pos: position{line: 44, col: 21, offset: 831},
								exprs: []interface{}{
									&labeledExpr{
										pos:   position{line: 44, col: 21, offset: 831},
										label: "l",
										expr: &ruleRefExpr{
											pos:  position{line: 44, col: 23, offset: 833},
											name: "Line",
										},
									},
									&litMatcher{
										pos:        position{line: 44, col: 28, offset: 838},
										val:        "\n",
										ignoreCase: false,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Line",
			pos:  position{line: 48, col: 1, offset: 915},
			expr: &actionExpr{
				pos: position{line: 48, col: 9, offset: 923},
				run: (*parser).callonLine1,
				expr: &seqExpr{
					pos: position{line: 48, col: 9, offset: 923},
					exprs: []interface{}{
						&zeroOrMoreExpr{
							pos: position{line: 48, col: 9, offset: 923},
							expr: &charClassMatcher{
								pos:        position{line: 48, col: 9, offset: 923},
								val:        "[ \\t]",
								chars:      []rune{' ', '\t'},
								ignoreCase: false,
								inverted:   false,
							},
						},
						&labeledExpr{
							pos:   position{line: 48, col: 16, offset: 930},
							label: "words",
							expr: &oneOrMoreExpr{
								pos: position{line: 48, col: 22, offset: 936},
								expr: &ruleRefExpr{
									pos:  position{line: 48, col: 23, offset: 937},
									name: "Word",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Word",
			pos:  position{line: 52, col: 1, offset: 991},
			expr: &actionExpr{
				pos: position{line: 52, col: 9, offset: 999},
				run: (*parser).callonWord1,
				expr: &seqExpr{
					pos: position{line: 52, col: 9, offset: 999},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 52, col: 9, offset: 999},
							expr: &ruleRefExpr{
								pos:  position{line: 52, col: 9, offset: 999},
								name: "Comment",
							},
						},
						&labeledExpr{
							pos:   position{line: 52, col: 18, offset: 1008},
							label: "val",
							expr: &choiceExpr{
								pos: position{line: 52, col: 23, offset: 1013},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 52, col: 23, offset: 1013},
										name: "String",
									},
									&ruleRefExpr{
										pos:  position{line: 52, col: 32, offset: 1022},
										name: "Invoke",
									},
									&ruleRefExpr{
										pos:  position{line: 52, col: 41, offset: 1031},
										name: "Interpolated",
									},
								},
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 52, col: 55, offset: 1045},
							expr: &ruleRefExpr{
								pos:  position{line: 52, col: 55, offset: 1045},
								name: "Comment",
							},
						},
					},
				},
			},
		},
		{
			name: "Comment",
			pos:  position{line: 56, col: 1, offset: 1077},
			expr: &seqExpr{
				pos: position{line: 56, col: 12, offset: 1088},
				exprs: []interface{}{
					&zeroOrMoreExpr{
						pos: position{line: 56, col: 12, offset: 1088},
						expr: &charClassMatcher{
							pos:        position{line: 56, col: 12, offset: 1088},
							val:        "[ \\t]",
							chars:      []rune{' ', '\t'},
							ignoreCase: false,
							inverted:   false,
						},
					},
					&litMatcher{
						pos:        position{line: 56, col: 19, offset: 1095},
						val:        "{-",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 56, col: 24, offset: 1100},
						expr: &choiceExpr{
							pos: position{line: 56, col: 25, offset: 1101},
							alternatives: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 56, col: 25, offset: 1101},
									name: "Comment",
								},
								&seqExpr{
									pos: position{line: 56, col: 35, offset: 1111},
									exprs: []interface{}{
										&notExpr{
											pos: position{line: 56, col: 35, offset: 1111},
											expr: &litMatcher{
												pos:        position{line: 56, col: 36, offset: 1112},
												val:        "-}",
												ignoreCase: false,
											},
										},
										&anyMatcher{
											line: 56, col: 41, offset: 1117,
										},
									},
								},
							},
						},
					},
					&litMatcher{
						pos:        position{line: 56, col: 45, offset: 1121},
						val:        "-}",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "Interpolated",
			pos:  position{line: 58, col: 1, offset: 1127},
			expr: &actionExpr{
				pos: position{line: 58, col: 17, offset: 1143},
				run: (*parser).callonInterpolated1,
				expr: &seqExpr{
					pos: position{line: 58, col: 17, offset: 1143},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 58, col: 17, offset: 1143},
							val:        "{",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 58, col: 21, offset: 1147},
							label: "word",
							expr: &zeroOrOneExpr{
								pos: position{line: 58, col: 26, offset: 1152},
								expr: &ruleRefExpr{
									pos:  position{line: 58, col: 26, offset: 1152},
									name: "Word",
								},
							},
						},
						&litMatcher{
							pos:        position{line: 58, col: 32, offset: 1158},
							val:        "}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "WrappedLine",
			pos:  position{line: 66, col: 1, offset: 1249},
			expr: &actionExpr{
				pos: position{line: 66, col: 16, offset: 1264},
				run: (*parser).callonWrappedLine1,
				expr: &seqExpr{
					pos: position{line: 66, col: 16, offset: 1264},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 66, col: 16, offset: 1264},
							label: "firstWord",
							expr: &ruleRefExpr{
								pos:  position{line: 66, col: 26, offset: 1274},
								name: "Word",
							},
						},
						&labeledExpr{
							pos:   position{line: 66, col: 31, offset: 1279},
							label: "words",
							expr: &zeroOrMoreExpr{
								pos: position{line: 66, col: 37, offset: 1285},
								expr: &choiceExpr{
									pos: position{line: 66, col: 38, offset: 1286},
									alternatives: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 66, col: 38, offset: 1286},
											name: "Word",
										},
										&ruleRefExpr{
											pos:  position{line: 66, col: 45, offset: 1293},
											name: "Split",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Split",
			pos:  position{line: 71, col: 1, offset: 1424},
			expr: &actionExpr{
				pos: position{line: 71, col: 10, offset: 1433},
				run: (*parser).callonSplit1,
				expr: &seqExpr{
					pos: position{line: 71, col: 10, offset: 1433},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 71, col: 10, offset: 1433},
							val:        "\n",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 71, col: 15, offset: 1438},
							expr: &charClassMatcher{
								pos:        position{line: 71, col: 15, offset: 1438},
								val:        "[ \\t]",
								chars:      []rune{' ', '\t'},
								ignoreCase: false,
								inverted:   false,
							},
						},
					},
				},
			},
		},
		{
			name: "String",
			pos:  position{line: 73, col: 1, offset: 1474},
			expr: &choiceExpr{
				pos: position{line: 73, col: 11, offset: 1484},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 73, col: 11, offset: 1484},
						run: (*parser).callonString2,
						expr: &labeledExpr{
							pos:   position{line: 73, col: 11, offset: 1484},
							label: "str",
							expr: &oneOrMoreExpr{
								pos: position{line: 73, col: 15, offset: 1488},
								expr: &charClassMatcher{
									pos:        position{line: 73, col: 15, offset: 1488},
									val:        "[^\\\\{}\\n]",
									chars:      []rune{'\\', '{', '}', '\n'},
									ignoreCase: false,
									inverted:   true,
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 73, col: 59, offset: 1532},
						run: (*parser).callonString6,
						expr: &seqExpr{
							pos: position{line: 73, col: 59, offset: 1532},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 73, col: 59, offset: 1532},
									val:        "\\",
									ignoreCase: false,
								},
								&charClassMatcher{
									pos:        position{line: 73, col: 64, offset: 1537},
									val:        "[\\\\{}]",
									chars:      []rune{'\\', '{', '}'},
									ignoreCase: false,
									inverted:   false,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "VerbatimString",
			pos:  position{line: 75, col: 1, offset: 1580},
			expr: &choiceExpr{
				pos: position{line: 75, col: 19, offset: 1598},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 75, col: 19, offset: 1598},
						run: (*parser).callonVerbatimString2,
						expr: &labeledExpr{
							pos:   position{line: 75, col: 19, offset: 1598},
							label: "str",
							expr: &oneOrMoreExpr{
								pos: position{line: 75, col: 23, offset: 1602},
								expr: &charClassMatcher{
									pos:        position{line: 75, col: 23, offset: 1602},
									val:        "[^\\n}]",
									chars:      []rune{'\n', '}'},
									ignoreCase: false,
									inverted:   true,
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 77, col: 5, offset: 1645},
						run: (*parser).callonVerbatimString6,
						expr: &seqExpr{
							pos: position{line: 77, col: 5, offset: 1645},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 77, col: 5, offset: 1645},
									val:        "}",
									ignoreCase: false,
								},
								&notExpr{
									pos: position{line: 77, col: 9, offset: 1649},
									expr: &litMatcher{
										pos:        position{line: 77, col: 10, offset: 1650},
										val:        "}}",
										ignoreCase: false,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "PreformattedLine",
			pos:  position{line: 81, col: 1, offset: 1689},
			expr: &actionExpr{
				pos: position{line: 81, col: 21, offset: 1709},
				run: (*parser).callonPreformattedLine1,
				expr: &seqExpr{
					pos: position{line: 81, col: 21, offset: 1709},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 81, col: 21, offset: 1709},
							label: "indent",
							expr: &ruleRefExpr{
								pos:  position{line: 81, col: 28, offset: 1716},
								name: "Indent",
							},
						},
						&labeledExpr{
							pos:   position{line: 81, col: 35, offset: 1723},
							label: "words",
							expr: &zeroOrMoreExpr{
								pos: position{line: 81, col: 41, offset: 1729},
								expr: &ruleRefExpr{
									pos:  position{line: 81, col: 41, offset: 1729},
									name: "Word",
								},
							},
						},
						&litMatcher{
							pos:        position{line: 81, col: 47, offset: 1735},
							val:        "\n",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Preformatted",
			pos:  position{line: 87, col: 1, offset: 1860},
			expr: &actionExpr{
				pos: position{line: 87, col: 17, offset: 1876},
				run: (*parser).callonPreformatted1,
				expr: &seqExpr{
					pos: position{line: 87, col: 17, offset: 1876},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 87, col: 17, offset: 1876},
							val:        "\n",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 87, col: 22, offset: 1881},
							label: "lines",
							expr: &zeroOrMoreExpr{
								pos: position{line: 87, col: 28, offset: 1887},
								expr: &ruleRefExpr{
									pos:  position{line: 87, col: 28, offset: 1887},
									name: "PreformattedLine",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "VerbatimLine",
			pos:  position{line: 92, col: 1, offset: 1999},
			expr: &actionExpr{
				pos: position{line: 92, col: 17, offset: 2015},
				run: (*parser).callonVerbatimLine1,
				expr: &seqExpr{
					pos: position{line: 92, col: 17, offset: 2015},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 92, col: 17, offset: 2015},
							label: "indent",
							expr: &ruleRefExpr{
								pos:  position{line: 92, col: 24, offset: 2022},
								name: "Indent",
							},
						},
						&labeledExpr{
							pos:   position{line: 92, col: 31, offset: 2029},
							label: "words",
							expr: &zeroOrMoreExpr{
								pos: position{line: 92, col: 37, offset: 2035},
								expr: &ruleRefExpr{
									pos:  position{line: 92, col: 37, offset: 2035},
									name: "VerbatimString",
								},
							},
						},
						&litMatcher{
							pos:        position{line: 92, col: 53, offset: 2051},
							val:        "\n",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Verbatim",
			pos:  position{line: 98, col: 1, offset: 2176},
			expr: &actionExpr{
				pos: position{line: 98, col: 13, offset: 2188},
				run: (*parser).callonVerbatim1,
				expr: &seqExpr{
					pos: position{line: 98, col: 13, offset: 2188},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 98, col: 13, offset: 2188},
							val:        "\n",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 98, col: 18, offset: 2193},
							label: "lines",
							expr: &zeroOrMoreExpr{
								pos: position{line: 98, col: 24, offset: 2199},
								expr: &ruleRefExpr{
									pos:  position{line: 98, col: 24, offset: 2199},
									name: "VerbatimLine",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Invoke",
			pos:  position{line: 103, col: 1, offset: 2307},
			expr: &actionExpr{
				pos: position{line: 103, col: 11, offset: 2317},
				run: (*parser).callonInvoke1,
				expr: &seqExpr{
					pos: position{line: 103, col: 11, offset: 2317},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 103, col: 11, offset: 2317},
							val:        "\\",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 103, col: 16, offset: 2322},
							label: "name",
							expr: &oneOrMoreExpr{
								pos: position{line: 103, col: 22, offset: 2328},
								expr: &charClassMatcher{
									pos:        position{line: 103, col: 22, offset: 2328},
									val:        "[a-z-]",
									chars:      []rune{'-'},
									ranges:     []rune{'a', 'z'},
									ignoreCase: false,
									inverted:   false,
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 103, col: 31, offset: 2337},
							label: "args",
							expr: &zeroOrMoreExpr{
								pos: position{line: 103, col: 37, offset: 2343},
								expr: &ruleRefExpr{
									pos:  position{line: 103, col: 37, offset: 2343},
									name: "Argument",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "VerbatimArg",
			pos:  position{line: 110, col: 1, offset: 2448},
			expr: &actionExpr{
				pos: position{line: 110, col: 16, offset: 2463},
				run: (*parser).callonVerbatimArg1,
				expr: &seqExpr{
					pos: position{line: 110, col: 16, offset: 2463},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 110, col: 16, offset: 2463},
							val:        "{{{",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 110, col: 22, offset: 2469},
							label: "node",
							expr: &ruleRefExpr{
								pos:  position{line: 110, col: 27, offset: 2474},
								name: "Verbatim",
							},
						},
						&zeroOrMoreExpr{
							pos: position{line: 110, col: 36, offset: 2483},
							expr: &charClassMatcher{
								pos:        position{line: 110, col: 36, offset: 2483},
								val:        "[ \\t]",
								chars:      []rune{' ', '\t'},
								ignoreCase: false,
								inverted:   false,
							},
						},
						&litMatcher{
							pos:        position{line: 110, col: 43, offset: 2490},
							val:        "}}}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "PreformattedArg",
			pos:  position{line: 114, col: 1, offset: 2520},
			expr: &actionExpr{
				pos: position{line: 114, col: 20, offset: 2539},
				run: (*parser).callonPreformattedArg1,
				expr: &seqExpr{
					pos: position{line: 114, col: 20, offset: 2539},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 114, col: 20, offset: 2539},
							val:        "{{",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 114, col: 25, offset: 2544},
							label: "node",
							expr: &ruleRefExpr{
								pos:  position{line: 114, col: 30, offset: 2549},
								name: "Preformatted",
							},
						},
						&zeroOrMoreExpr{
							pos: position{line: 114, col: 43, offset: 2562},
							expr: &charClassMatcher{
								pos:        position{line: 114, col: 43, offset: 2562},
								val:        "[ \\t]",
								chars:      []rune{' ', '\t'},
								ignoreCase: false,
								inverted:   false,
							},
						},
						&litMatcher{
							pos:        position{line: 114, col: 50, offset: 2569},
							val:        "}}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Arg",
			pos:  position{line: 118, col: 1, offset: 2598},
			expr: &actionExpr{
				pos: position{line: 118, col: 8, offset: 2605},
				run: (*parser).callonArg1,
				expr: &seqExpr{
					pos: position{line: 118, col: 8, offset: 2605},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 118, col: 8, offset: 2605},
							val:        "{",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 118, col: 12, offset: 2609},
							label: "node",
							expr: &zeroOrOneExpr{
								pos: position{line: 118, col: 17, offset: 2614},
								expr: &choiceExpr{
									pos: position{line: 118, col: 18, offset: 2615},
									alternatives: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 118, col: 18, offset: 2615},
											name: "WrappedLine",
										},
										&ruleRefExpr{
											pos:  position{line: 118, col: 32, offset: 2629},
											name: "ParaArg",
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 118, col: 42, offset: 2639},
							val:        "}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "ParaArg",
			pos:  position{line: 126, col: 1, offset: 2730},
			expr: &actionExpr{
				pos: position{line: 126, col: 12, offset: 2741},
				run: (*parser).callonParaArg1,
				expr: &seqExpr{
					pos: position{line: 126, col: 12, offset: 2741},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 126, col: 12, offset: 2741},
							label: "paras",
							expr: &ruleRefExpr{
								pos:  position{line: 126, col: 18, offset: 2747},
								name: "Paragraphs",
							},
						},
						&zeroOrMoreExpr{
							pos: position{line: 126, col: 29, offset: 2758},
							expr: &charClassMatcher{
								pos:        position{line: 126, col: 29, offset: 2758},
								val:        "[ \\t]",
								chars:      []rune{' ', '\t'},
								ignoreCase: false,
								inverted:   false,
							},
						},
					},
				},
			},
		},
		{
			name: "Argument",
			pos:  position{line: 130, col: 1, offset: 2790},
			expr: &choiceExpr{
				pos: position{line: 130, col: 13, offset: 2802},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 130, col: 13, offset: 2802},
						name: "VerbatimArg",
					},
					&ruleRefExpr{
						pos:  position{line: 130, col: 27, offset: 2816},
						name: "PreformattedArg",
					},
					&ruleRefExpr{
						pos:  position{line: 130, col: 45, offset: 2834},
						name: "Arg",
					},
				},
			},
		},
		{
			name: "Indent",
			pos:  position{line: 132, col: 1, offset: 2839},
			expr: &actionExpr{
				pos: position{line: 132, col: 11, offset: 2849},
				run: (*parser).callonIndent1,
				expr: &zeroOrMoreExpr{
					pos: position{line: 132, col: 11, offset: 2849},
					expr: &charClassMatcher{
						pos:        position{line: 132, col: 11, offset: 2849},
						val:        "[ \\t]",
						chars:      []rune{' ', '\t'},
						ignoreCase: false,
						inverted:   false,
					},
				},
			},
		},
	},
}

func (c *current) onBooklit1(node interface{}) (interface{}, error) {
	return node, nil
}

func (p *parser) callonBooklit1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBooklit1(stack["node"])
}

func (c *current) onParagraphs7(p interface{}) (interface{}, error) {
	return p, nil
}

func (p *parser) callonParagraphs7() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onParagraphs7(stack["p"])
}

func (c *current) onParagraphs1(paragraphs interface{}) (interface{}, error) {
	return Sequence(ifaceNodes(paragraphs)), nil
}

func (p *parser) callonParagraphs1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onParagraphs1(stack["paragraphs"])
}

func (c *current) onParagraph4(l interface{}) (interface{}, error) {
	return l, nil
}

func (p *parser) callonParagraph4() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onParagraph4(stack["l"])
}

func (c *current) onParagraph1(lines interface{}) (interface{}, error) {
	return Paragraph(ifaceSequences(lines)), nil
}

func (p *parser) callonParagraph1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onParagraph1(stack["lines"])
}

func (c *current) onLine1(words interface{}) (interface{}, error) {
	return Sequence(ifaceNodes(words)), nil
}

func (p *parser) callonLine1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLine1(stack["words"])
}

func (c *current) onWord1(val interface{}) (interface{}, error) {
	return val, nil
}

func (p *parser) callonWord1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onWord1(stack["val"])
}

func (c *current) onInterpolated1(word interface{}) (interface{}, error) {
	if word == nil {
		return Sequence{}, nil
	} else {
		return word, nil
	}
}

func (p *parser) callonInterpolated1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onInterpolated1(stack["word"])
}

func (c *current) onWrappedLine1(firstWord, words interface{}) (interface{}, error) {
	allWords := append([]interface{}{firstWord}, words.([]interface{})...)
	return Sequence(ifaceNodes(allWords)), nil
}

func (p *parser) callonWrappedLine1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onWrappedLine1(stack["firstWord"], stack["words"])
}

func (c *current) onSplit1() (interface{}, error) {
	return String(" "), nil
}

func (p *parser) callonSplit1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSplit1()
}

func (c *current) onString2(str interface{}) (interface{}, error) {
	return String(c.text), nil
}

func (p *parser) callonString2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onString2(stack["str"])
}

func (c *current) onString6() (interface{}, error) {
	return String(c.text[1:]), nil
}

func (p *parser) callonString6() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onString6()
}

func (c *current) onVerbatimString2(str interface{}) (interface{}, error) {
	return String(c.text), nil
}

func (p *parser) callonVerbatimString2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onVerbatimString2(stack["str"])
}

func (c *current) onVerbatimString6() (interface{}, error) {
	return String(c.text), nil
}

func (p *parser) callonVerbatimString6() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onVerbatimString6()
}

func (c *current) onPreformattedLine1(indent, words interface{}) (interface{}, error) {
	line := []Node{String(indent.(string))}
	line = append(line, ifaceNodes(words)...)
	return Sequence(line), nil
}

func (p *parser) callonPreformattedLine1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPreformattedLine1(stack["indent"], stack["words"])
}

func (c *current) onPreformatted1(lines interface{}) (interface{}, error) {
	delete(c.globalStore, "indent-skip")
	return Preformatted(ifaceSequences(lines)), nil
}

func (p *parser) callonPreformatted1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPreformatted1(stack["lines"])
}

func (c *current) onVerbatimLine1(indent, words interface{}) (interface{}, error) {
	line := []Node{String(indent.(string))}
	line = append(line, ifaceNodes(words)...)
	return Sequence(line), nil
}

func (p *parser) callonVerbatimLine1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onVerbatimLine1(stack["indent"], stack["words"])
}

func (c *current) onVerbatim1(lines interface{}) (interface{}, error) {
	delete(c.globalStore, "indent-skip")
	return Preformatted(ifaceSequences(lines)), nil
}

func (p *parser) callonVerbatim1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onVerbatim1(stack["lines"])
}

func (c *current) onInvoke1(name, args interface{}) (interface{}, error) {
	return Invoke{
		Function:  ifaceStr(name),
		Arguments: ifaceNodes(args),
	}, nil
}

func (p *parser) callonInvoke1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onInvoke1(stack["name"], stack["args"])
}

func (c *current) onVerbatimArg1(node interface{}) (interface{}, error) {
	return node, nil
}

func (p *parser) callonVerbatimArg1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onVerbatimArg1(stack["node"])
}

func (c *current) onPreformattedArg1(node interface{}) (interface{}, error) {
	return node, nil
}

func (p *parser) callonPreformattedArg1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPreformattedArg1(stack["node"])
}

func (c *current) onArg1(node interface{}) (interface{}, error) {
	if node == nil {
		return String(""), nil
	} else {
		return node, nil
	}
}

func (p *parser) callonArg1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onArg1(stack["node"])
}

func (c *current) onParaArg1(paras interface{}) (interface{}, error) {
	return paras, nil
}

func (p *parser) callonParaArg1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onParaArg1(stack["paras"])
}

func (c *current) onIndent1() (interface{}, error) {
	skip := len(c.text)

	i, found := c.globalStore["indent-skip"]
	if found {
		skip = i.(int)
	} else {
		c.globalStore["indent-skip"] = skip
	}

	if skip <= len(c.text) {
		return string(c.text[skip:]), nil
	} else {
		return "", nil
	}
}

func (p *parser) callonIndent1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIndent1()
}

var (
	// errNoRule is returned when the grammar to parse has no rule.
	errNoRule = errors.New("grammar has no rule")

	// errInvalidEncoding is returned when the source is not properly
	// utf8-encoded.
	errInvalidEncoding = errors.New("invalid encoding")
)

// Option is a function that can set an option on the parser. It returns
// the previous setting as an Option.
type Option func(*parser) Option

// Debug creates an Option to set the debug flag to b. When set to true,
// debugging information is printed to stdout while parsing.
//
// The default is false.
func Debug(b bool) Option {
	return func(p *parser) Option {
		old := p.debug
		p.debug = b
		return Debug(old)
	}
}

// Memoize creates an Option to set the memoize flag to b. When set to true,
// the parser will cache all results so each expression is evaluated only
// once. This guarantees linear parsing time even for pathological cases,
// at the expense of more memory and slower times for typical cases.
//
// The default is false.
func Memoize(b bool) Option {
	return func(p *parser) Option {
		old := p.memoize
		p.memoize = b
		return Memoize(old)
	}
}

// Recover creates an Option to set the recover flag to b. When set to
// true, this causes the parser to recover from panics and convert it
// to an error. Setting it to false can be useful while debugging to
// access the full stack trace.
//
// The default is true.
func Recover(b bool) Option {
	return func(p *parser) Option {
		old := p.recover
		p.recover = b
		return Recover(old)
	}
}

// GlobalStore creates an Option to set a key to a certain value in
// the globalStore.
func GlobalStore(key string, value interface{}) Option {
	return func(p *parser) Option {
		old := p.cur.globalStore[key]
		p.cur.globalStore[key] = value
		return GlobalStore(key, old)
	}
}

// ParseFile parses the file identified by filename.
func ParseFile(filename string, opts ...Option) (i interface{}, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = f.Close()
	}()
	return ParseReader(filename, f, opts...)
}

// ParseReader parses the data from r using filename as information in the
// error messages.
func ParseReader(filename string, r io.Reader, opts ...Option) (interface{}, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return Parse(filename, b, opts...)
}

// Parse parses the data from b using filename as information in the
// error messages.
func Parse(filename string, b []byte, opts ...Option) (interface{}, error) {
	return newParser(filename, b, opts...).parse(g)
}

// position records a position in the text.
type position struct {
	line, col, offset int
}

func (p position) String() string {
	return fmt.Sprintf("%d:%d [%d]", p.line, p.col, p.offset)
}

// savepoint stores all state required to go back to this point in the
// parser.
type savepoint struct {
	position
	rn rune
	w  int
}

type current struct {
	pos  position // start position of the match
	text []byte   // raw text of the match

	// the globalStore allows the parser to store arbitrary values
	globalStore map[string]interface{}
}

// the AST types...

type grammar struct {
	pos   position
	rules []*rule
}

type rule struct {
	pos         position
	name        string
	displayName string
	expr        interface{}
}

type choiceExpr struct {
	pos          position
	alternatives []interface{}
}

type actionExpr struct {
	pos  position
	expr interface{}
	run  func(*parser) (interface{}, error)
}

type seqExpr struct {
	pos   position
	exprs []interface{}
}

type labeledExpr struct {
	pos   position
	label string
	expr  interface{}
}

type expr struct {
	pos  position
	expr interface{}
}

type andExpr expr
type notExpr expr
type zeroOrOneExpr expr
type zeroOrMoreExpr expr
type oneOrMoreExpr expr

type ruleRefExpr struct {
	pos  position
	name string
}

type andCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type notCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type litMatcher struct {
	pos        position
	val        string
	ignoreCase bool
}

type charClassMatcher struct {
	pos        position
	val        string
	chars      []rune
	ranges     []rune
	classes    []*unicode.RangeTable
	ignoreCase bool
	inverted   bool
}

type anyMatcher position

// errList cumulates the errors found by the parser.
type errList []error

func (e *errList) add(err error) {
	*e = append(*e, err)
}

func (e errList) err() error {
	if len(e) == 0 {
		return nil
	}
	e.dedupe()
	return e
}

func (e *errList) dedupe() {
	var cleaned []error
	set := make(map[string]bool)
	for _, err := range *e {
		if msg := err.Error(); !set[msg] {
			set[msg] = true
			cleaned = append(cleaned, err)
		}
	}
	*e = cleaned
}

func (e errList) Error() string {
	switch len(e) {
	case 0:
		return ""
	case 1:
		return e[0].Error()
	default:
		var buf bytes.Buffer

		for i, err := range e {
			if i > 0 {
				buf.WriteRune('\n')
			}
			buf.WriteString(err.Error())
		}
		return buf.String()
	}
}

// parserError wraps an error with a prefix indicating the rule in which
// the error occurred. The original error is stored in the Inner field.
type parserError struct {
	Inner    error
	pos      position
	prefix   string
	expected []string
}

// Error returns the error message.
func (p *parserError) Error() string {
	return p.prefix + ": " + p.Inner.Error()
}

// newParser creates a parser with the specified input source and options.
func newParser(filename string, b []byte, opts ...Option) *parser {
	p := &parser{
		filename: filename,
		errs:     new(errList),
		data:     b,
		pt:       savepoint{position: position{line: 1}},
		recover:  true,
		cur: current{
			globalStore: make(map[string]interface{}),
		},
		maxFailPos:      position{col: 1, line: 1},
		maxFailExpected: make(map[string]struct{}),
	}
	p.setOptions(opts)
	return p
}

// setOptions applies the options to the parser.
func (p *parser) setOptions(opts []Option) {
	for _, opt := range opts {
		opt(p)
	}
}

type resultTuple struct {
	v   interface{}
	b   bool
	end savepoint
}

type parser struct {
	filename string
	pt       savepoint
	cur      current

	data []byte
	errs *errList

	depth   int
	recover bool
	debug   bool

	memoize bool
	// memoization table for the packrat algorithm:
	// map[offset in source] map[expression or rule] {value, match}
	memo map[int]map[interface{}]resultTuple

	// rules table, maps the rule identifier to the rule node
	rules map[string]*rule
	// variables stack, map of label to value
	vstack []map[string]interface{}
	// rule stack, allows identification of the current rule in errors
	rstack []*rule

	// stats
	exprCnt int

	// parse fail
	maxFailPos            position
	maxFailExpected       map[string]struct{}
	maxFailInvertExpected bool
}

// push a variable set on the vstack.
func (p *parser) pushV() {
	if cap(p.vstack) == len(p.vstack) {
		// create new empty slot in the stack
		p.vstack = append(p.vstack, nil)
	} else {
		// slice to 1 more
		p.vstack = p.vstack[:len(p.vstack)+1]
	}

	// get the last args set
	m := p.vstack[len(p.vstack)-1]
	if m != nil && len(m) == 0 {
		// empty map, all good
		return
	}

	m = make(map[string]interface{})
	p.vstack[len(p.vstack)-1] = m
}

// pop a variable set from the vstack.
func (p *parser) popV() {
	// if the map is not empty, clear it
	m := p.vstack[len(p.vstack)-1]
	if len(m) > 0 {
		// GC that map
		p.vstack[len(p.vstack)-1] = nil
	}
	p.vstack = p.vstack[:len(p.vstack)-1]
}

func (p *parser) print(prefix, s string) string {
	if !p.debug {
		return s
	}

	fmt.Printf("%s %d:%d:%d: %s [%#U]\n",
		prefix, p.pt.line, p.pt.col, p.pt.offset, s, p.pt.rn)
	return s
}

func (p *parser) in(s string) string {
	p.depth++
	return p.print(strings.Repeat(" ", p.depth)+">", s)
}

func (p *parser) out(s string) string {
	p.depth--
	return p.print(strings.Repeat(" ", p.depth)+"<", s)
}

func (p *parser) addErr(err error) {
	p.addErrAt(err, p.pt.position, []string{})
}

func (p *parser) addErrAt(err error, pos position, expected []string) {
	var buf bytes.Buffer
	if p.filename != "" {
		buf.WriteString(p.filename)
	}
	if buf.Len() > 0 {
		buf.WriteString(":")
	}
	buf.WriteString(fmt.Sprintf("%d:%d (%d)", pos.line, pos.col, pos.offset))
	if len(p.rstack) > 0 {
		if buf.Len() > 0 {
			buf.WriteString(": ")
		}
		rule := p.rstack[len(p.rstack)-1]
		if rule.displayName != "" {
			buf.WriteString("rule " + rule.displayName)
		} else {
			buf.WriteString("rule " + rule.name)
		}
	}
	pe := &parserError{Inner: err, pos: pos, prefix: buf.String(), expected: expected}
	p.errs.add(pe)
}

func (p *parser) failAt(fail bool, pos position, want string) {
	// process fail if parsing fails and not inverted or parsing succeeds and invert is set
	if fail == p.maxFailInvertExpected {
		if pos.offset < p.maxFailPos.offset {
			return
		}

		if pos.offset > p.maxFailPos.offset {
			p.maxFailPos = pos
			p.maxFailExpected = make(map[string]struct{})
		}

		if p.maxFailInvertExpected {
			want = "!" + want
		}
		p.maxFailExpected[want] = struct{}{}
	}
}

// read advances the parser to the next rune.
func (p *parser) read() {
	p.pt.offset += p.pt.w
	rn, n := utf8.DecodeRune(p.data[p.pt.offset:])
	p.pt.rn = rn
	p.pt.w = n
	p.pt.col++
	if rn == '\n' {
		p.pt.line++
		p.pt.col = 0
	}

	if rn == utf8.RuneError {
		if n == 1 {
			p.addErr(errInvalidEncoding)
		}
	}
}

// restore parser position to the savepoint pt.
func (p *parser) restore(pt savepoint) {
	if p.debug {
		defer p.out(p.in("restore"))
	}
	if pt.offset == p.pt.offset {
		return
	}
	p.pt = pt
}

// get the slice of bytes from the savepoint start to the current position.
func (p *parser) sliceFrom(start savepoint) []byte {
	return p.data[start.position.offset:p.pt.position.offset]
}

func (p *parser) getMemoized(node interface{}) (resultTuple, bool) {
	if len(p.memo) == 0 {
		return resultTuple{}, false
	}
	m := p.memo[p.pt.offset]
	if len(m) == 0 {
		return resultTuple{}, false
	}
	res, ok := m[node]
	return res, ok
}

func (p *parser) setMemoized(pt savepoint, node interface{}, tuple resultTuple) {
	if p.memo == nil {
		p.memo = make(map[int]map[interface{}]resultTuple)
	}
	m := p.memo[pt.offset]
	if m == nil {
		m = make(map[interface{}]resultTuple)
		p.memo[pt.offset] = m
	}
	m[node] = tuple
}

func (p *parser) buildRulesTable(g *grammar) {
	p.rules = make(map[string]*rule, len(g.rules))
	for _, r := range g.rules {
		p.rules[r.name] = r
	}
}

func (p *parser) parse(g *grammar) (val interface{}, err error) {
	if len(g.rules) == 0 {
		p.addErr(errNoRule)
		return nil, p.errs.err()
	}

	// TODO : not super critical but this could be generated
	p.buildRulesTable(g)

	if p.recover {
		// panic can be used in action code to stop parsing immediately
		// and return the panic as an error.
		defer func() {
			if e := recover(); e != nil {
				if p.debug {
					defer p.out(p.in("panic handler"))
				}
				val = nil
				switch e := e.(type) {
				case error:
					p.addErr(e)
				default:
					p.addErr(fmt.Errorf("%v", e))
				}
				err = p.errs.err()
			}
		}()
	}

	// start rule is rule [0]
	p.read() // advance to first rune
	val, ok := p.parseRule(g.rules[0])
	if !ok {
		if len(*p.errs) == 0 {
			// If parsing fails, but no errors have been recorded, the expected values
			// for the farthest parser position are returned as error.
			expected := make([]string, 0, len(p.maxFailExpected))
			eof := false
			if _, ok := p.maxFailExpected["!."]; ok {
				delete(p.maxFailExpected, "!.")
				eof = true
			}
			for k := range p.maxFailExpected {
				expected = append(expected, k)
			}
			sort.Strings(expected)
			if eof {
				expected = append(expected, "EOF")
			}
			p.addErrAt(errors.New("no match found, expected: "+listJoin(expected, ", ", "or")), p.maxFailPos, expected)
		}
		return nil, p.errs.err()
	}
	return val, p.errs.err()
}

func listJoin(list []string, sep string, lastSep string) string {
	switch len(list) {
	case 0:
		return ""
	case 1:
		return list[0]
	default:
		return fmt.Sprintf("%s %s %s", strings.Join(list[:len(list)-1], sep), lastSep, list[len(list)-1])
	}
}

func (p *parser) parseRule(rule *rule) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRule " + rule.name))
	}

	if p.memoize {
		res, ok := p.getMemoized(rule)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
	}

	start := p.pt
	p.rstack = append(p.rstack, rule)
	p.pushV()
	val, ok := p.parseExpr(rule.expr)
	p.popV()
	p.rstack = p.rstack[:len(p.rstack)-1]
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}

	if p.memoize {
		p.setMemoized(start, rule, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseExpr(expr interface{}) (interface{}, bool) {
	var pt savepoint

	if p.memoize {
		res, ok := p.getMemoized(expr)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
		pt = p.pt
	}

	p.exprCnt++
	var val interface{}
	var ok bool
	switch expr := expr.(type) {
	case *actionExpr:
		val, ok = p.parseActionExpr(expr)
	case *andCodeExpr:
		val, ok = p.parseAndCodeExpr(expr)
	case *andExpr:
		val, ok = p.parseAndExpr(expr)
	case *anyMatcher:
		val, ok = p.parseAnyMatcher(expr)
	case *charClassMatcher:
		val, ok = p.parseCharClassMatcher(expr)
	case *choiceExpr:
		val, ok = p.parseChoiceExpr(expr)
	case *labeledExpr:
		val, ok = p.parseLabeledExpr(expr)
	case *litMatcher:
		val, ok = p.parseLitMatcher(expr)
	case *notCodeExpr:
		val, ok = p.parseNotCodeExpr(expr)
	case *notExpr:
		val, ok = p.parseNotExpr(expr)
	case *oneOrMoreExpr:
		val, ok = p.parseOneOrMoreExpr(expr)
	case *ruleRefExpr:
		val, ok = p.parseRuleRefExpr(expr)
	case *seqExpr:
		val, ok = p.parseSeqExpr(expr)
	case *zeroOrMoreExpr:
		val, ok = p.parseZeroOrMoreExpr(expr)
	case *zeroOrOneExpr:
		val, ok = p.parseZeroOrOneExpr(expr)
	default:
		panic(fmt.Sprintf("unknown expression type %T", expr))
	}
	if p.memoize {
		p.setMemoized(pt, expr, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseActionExpr(act *actionExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseActionExpr"))
	}

	start := p.pt
	val, ok := p.parseExpr(act.expr)
	if ok {
		p.cur.pos = start.position
		p.cur.text = p.sliceFrom(start)
		actVal, err := act.run(p)
		if err != nil {
			p.addErrAt(err, start.position, []string{})
		}
		val = actVal
	}
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}
	return val, ok
}

func (p *parser) parseAndCodeExpr(and *andCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndCodeExpr"))
	}

	ok, err := and.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, ok
}

func (p *parser) parseAndExpr(and *andExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndExpr"))
	}

	pt := p.pt
	p.pushV()
	_, ok := p.parseExpr(and.expr)
	p.popV()
	p.restore(pt)
	return nil, ok
}

func (p *parser) parseAnyMatcher(any *anyMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAnyMatcher"))
	}

	if p.pt.rn != utf8.RuneError {
		start := p.pt
		p.read()
		p.failAt(true, start.position, ".")
		return p.sliceFrom(start), true
	}
	p.failAt(false, p.pt.position, ".")
	return nil, false
}

func (p *parser) parseCharClassMatcher(chr *charClassMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseCharClassMatcher"))
	}

	cur := p.pt.rn
	start := p.pt
	// can't match EOF
	if cur == utf8.RuneError {
		p.failAt(false, start.position, chr.val)
		return nil, false
	}
	if chr.ignoreCase {
		cur = unicode.ToLower(cur)
	}

	// try to match in the list of available chars
	for _, rn := range chr.chars {
		if rn == cur {
			if chr.inverted {
				p.failAt(false, start.position, chr.val)
				return nil, false
			}
			p.read()
			p.failAt(true, start.position, chr.val)
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of ranges
	for i := 0; i < len(chr.ranges); i += 2 {
		if cur >= chr.ranges[i] && cur <= chr.ranges[i+1] {
			if chr.inverted {
				p.failAt(false, start.position, chr.val)
				return nil, false
			}
			p.read()
			p.failAt(true, start.position, chr.val)
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of Unicode classes
	for _, cl := range chr.classes {
		if unicode.Is(cl, cur) {
			if chr.inverted {
				p.failAt(false, start.position, chr.val)
				return nil, false
			}
			p.read()
			p.failAt(true, start.position, chr.val)
			return p.sliceFrom(start), true
		}
	}

	if chr.inverted {
		p.read()
		p.failAt(true, start.position, chr.val)
		return p.sliceFrom(start), true
	}
	p.failAt(false, start.position, chr.val)
	return nil, false
}

func (p *parser) parseChoiceExpr(ch *choiceExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseChoiceExpr"))
	}

	for _, alt := range ch.alternatives {
		p.pushV()
		val, ok := p.parseExpr(alt)
		p.popV()
		if ok {
			return val, ok
		}
	}
	return nil, false
}

func (p *parser) parseLabeledExpr(lab *labeledExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLabeledExpr"))
	}

	p.pushV()
	val, ok := p.parseExpr(lab.expr)
	p.popV()
	if ok && lab.label != "" {
		m := p.vstack[len(p.vstack)-1]
		m[lab.label] = val
	}
	return val, ok
}

func (p *parser) parseLitMatcher(lit *litMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLitMatcher"))
	}

	ignoreCase := ""
	if lit.ignoreCase {
		ignoreCase = "i"
	}
	val := fmt.Sprintf("%q%s", lit.val, ignoreCase)
	start := p.pt
	for _, want := range lit.val {
		cur := p.pt.rn
		if lit.ignoreCase {
			cur = unicode.ToLower(cur)
		}
		if cur != want {
			p.failAt(false, start.position, val)
			p.restore(start)
			return nil, false
		}
		p.read()
	}
	p.failAt(true, start.position, val)
	return p.sliceFrom(start), true
}

func (p *parser) parseNotCodeExpr(not *notCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotCodeExpr"))
	}

	ok, err := not.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, !ok
}

func (p *parser) parseNotExpr(not *notExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotExpr"))
	}

	pt := p.pt
	p.pushV()
	p.maxFailInvertExpected = !p.maxFailInvertExpected
	_, ok := p.parseExpr(not.expr)
	p.maxFailInvertExpected = !p.maxFailInvertExpected
	p.popV()
	p.restore(pt)
	return nil, !ok
}

func (p *parser) parseOneOrMoreExpr(expr *oneOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseOneOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			if len(vals) == 0 {
				// did not match once, no match
				return nil, false
			}
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseRuleRefExpr(ref *ruleRefExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRuleRefExpr " + ref.name))
	}

	if ref.name == "" {
		panic(fmt.Sprintf("%s: invalid rule: missing name", ref.pos))
	}

	rule := p.rules[ref.name]
	if rule == nil {
		p.addErr(fmt.Errorf("undefined rule: %s", ref.name))
		return nil, false
	}
	return p.parseRule(rule)
}

func (p *parser) parseSeqExpr(seq *seqExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseSeqExpr"))
	}

	var vals []interface{}

	pt := p.pt
	for _, expr := range seq.exprs {
		val, ok := p.parseExpr(expr)
		if !ok {
			p.restore(pt)
			return nil, false
		}
		vals = append(vals, val)
	}
	return vals, true
}

func (p *parser) parseZeroOrMoreExpr(expr *zeroOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseZeroOrOneExpr(expr *zeroOrOneExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrOneExpr"))
	}

	p.pushV()
	val, _ := p.parseExpr(expr.expr)
	p.popV()
	// whether it matched or not, consider it a match
	return val, true
}
