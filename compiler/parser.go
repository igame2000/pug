package compiler

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

func toSlice(v interface{}) []interface{} {
	if v == nil {
		return nil
	}
	return v.([]interface{})
}

func pos(p position) Position {
	return Position{
		Line:   p.line,
		Col:    p.col,
		Offset: p.offset,
	}
}

func binary(first, rest interface{}, curpos position) (Expression, error) {
	restElem := toSlice(rest)

	if len(restElem) == 0 {
		return first.(Expression), nil
	}

	var cur Expression = first.(Expression)

	for _, x := range restElem {
		elem := toSlice(x)
		cur = &BinaryExpression{X: cur, Y: elem[3].(Expression), Op: elem[1].(string), GraphNode: NewNode(pos(curpos))}
	}

	return cur, nil
}

var g = &grammar{
	rules: []*rule{
		{
			name: "Input",
			pos:  position{line: 37, col: 1, offset: 708},
			expr: &actionExpr{
				pos: position{line: 37, col: 10, offset: 717},
				run: (*parser).callonInput1,
				expr: &seqExpr{
					pos: position{line: 37, col: 10, offset: 717},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 37, col: 10, offset: 717},
							label: "l",
							expr: &ruleRefExpr{
								pos:  position{line: 37, col: 12, offset: 719},
								name: "List",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 37, col: 17, offset: 724},
							name: "EOF",
						},
					},
				},
			},
		},
		{
			name: "List",
			pos:  position{line: 41, col: 1, offset: 802},
			expr: &choiceExpr{
				pos: position{line: 41, col: 9, offset: 810},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 41, col: 9, offset: 810},
						run: (*parser).callonList2,
						expr: &seqExpr{
							pos: position{line: 41, col: 9, offset: 810},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 41, col: 9, offset: 810},
									label: "node",
									expr: &ruleRefExpr{
										pos:  position{line: 41, col: 14, offset: 815},
										name: "ListNode",
									},
								},
								&labeledExpr{
									pos:   position{line: 41, col: 23, offset: 824},
									label: "list",
									expr: &ruleRefExpr{
										pos:  position{line: 41, col: 28, offset: 829},
										name: "List",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 49, col: 5, offset: 1015},
						run: (*parser).callonList8,
						expr: &andExpr{
							pos: position{line: 49, col: 5, offset: 1015},
							expr: &choiceExpr{
								pos: position{line: 49, col: 7, offset: 1017},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 49, col: 7, offset: 1017},
										name: "Outdent",
									},
									&ruleRefExpr{
										pos:  position{line: 49, col: 17, offset: 1027},
										name: "EOF",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "IndentedList",
			pos:  position{line: 53, col: 1, offset: 1089},
			expr: &actionExpr{
				pos: position{line: 53, col: 17, offset: 1105},
				run: (*parser).callonIndentedList1,
				expr: &seqExpr{
					pos: position{line: 53, col: 17, offset: 1105},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 53, col: 17, offset: 1105},
							name: "Indent",
						},
						&labeledExpr{
							pos:   position{line: 53, col: 24, offset: 1112},
							label: "list",
							expr: &ruleRefExpr{
								pos:  position{line: 53, col: 29, offset: 1117},
								name: "List",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 53, col: 34, offset: 1122},
							name: "Outdent",
						},
					},
				},
			},
		},
		{
			name: "IndentedRawText",
			pos:  position{line: 57, col: 1, offset: 1154},
			expr: &actionExpr{
				pos: position{line: 57, col: 20, offset: 1173},
				run: (*parser).callonIndentedRawText1,
				expr: &seqExpr{
					pos: position{line: 57, col: 20, offset: 1173},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 57, col: 20, offset: 1173},
							name: "Indent",
						},
						&labeledExpr{
							pos:   position{line: 57, col: 27, offset: 1180},
							label: "t",
							expr: &ruleRefExpr{
								pos:  position{line: 57, col: 29, offset: 1182},
								name: "RawText",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 57, col: 37, offset: 1190},
							name: "Outdent",
						},
					},
				},
			},
		},
		{
			name: "RawText",
			pos:  position{line: 61, col: 1, offset: 1274},
			expr: &choiceExpr{
				pos: position{line: 61, col: 12, offset: 1285},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 61, col: 12, offset: 1285},
						run: (*parser).callonRawText2,
						expr: &seqExpr{
							pos: position{line: 61, col: 12, offset: 1285},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 61, col: 12, offset: 1285},
									name: "Indent",
								},
								&labeledExpr{
									pos:   position{line: 61, col: 19, offset: 1292},
									label: "rt",
									expr: &ruleRefExpr{
										pos:  position{line: 61, col: 22, offset: 1295},
										name: "RawText",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 61, col: 30, offset: 1303},
									name: "Outdent",
								},
								&labeledExpr{
									pos:   position{line: 61, col: 38, offset: 1311},
									label: "tail",
									expr: &ruleRefExpr{
										pos:  position{line: 61, col: 43, offset: 1316},
										name: "RawText",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 63, col: 5, offset: 1372},
						run: (*parser).callonRawText10,
						expr: &choiceExpr{
							pos: position{line: 63, col: 6, offset: 1373},
							alternatives: []interface{}{
								&andExpr{
									pos: position{line: 63, col: 6, offset: 1373},
									expr: &ruleRefExpr{
										pos:  position{line: 63, col: 7, offset: 1374},
										name: "Outdent",
									},
								},
								&andExpr{
									pos: position{line: 63, col: 17, offset: 1384},
									expr: &ruleRefExpr{
										pos:  position{line: 63, col: 18, offset: 1385},
										name: "EOF",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 65, col: 5, offset: 1413},
						run: (*parser).callonRawText16,
						expr: &seqExpr{
							pos: position{line: 65, col: 5, offset: 1413},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 65, col: 5, offset: 1413},
									label: "head",
									expr: &anyMatcher{
										line: 65, col: 10, offset: 1418,
									},
								},
								&labeledExpr{
									pos:   position{line: 65, col: 12, offset: 1420},
									label: "tail",
									expr: &ruleRefExpr{
										pos:  position{line: 65, col: 17, offset: 1425},
										name: "RawText",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:        "ListNode",
			displayName: "\"listnode\"",
			pos:         position{line: 69, col: 1, offset: 1490},
			expr: &choiceExpr{
				pos: position{line: 70, col: 3, offset: 1515},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 70, col: 3, offset: 1515},
						name: "Comment",
					},
					&ruleRefExpr{
						pos:  position{line: 71, col: 3, offset: 1527},
						name: "Import",
					},
					&ruleRefExpr{
						pos:  position{line: 72, col: 3, offset: 1538},
						name: "Extend",
					},
					&ruleRefExpr{
						pos:  position{line: 73, col: 3, offset: 1549},
						name: "PipeText",
					},
					&ruleRefExpr{
						pos:  position{line: 74, col: 3, offset: 1562},
						name: "If",
					},
					&ruleRefExpr{
						pos:  position{line: 75, col: 3, offset: 1569},
						name: "Each",
					},
					&ruleRefExpr{
						pos:  position{line: 76, col: 3, offset: 1578},
						name: "DocType",
					},
					&ruleRefExpr{
						pos:  position{line: 77, col: 3, offset: 1590},
						name: "Mixin",
					},
					&ruleRefExpr{
						pos:  position{line: 78, col: 3, offset: 1600},
						name: "MixinCall",
					},
					&ruleRefExpr{
						pos:  position{line: 79, col: 3, offset: 1614},
						name: "Assignment",
					},
					&ruleRefExpr{
						pos:  position{line: 80, col: 3, offset: 1629},
						name: "Block",
					},
					&ruleRefExpr{
						pos:  position{line: 81, col: 3, offset: 1639},
						name: "Tag",
					},
					&actionExpr{
						pos: position{line: 82, col: 3, offset: 1647},
						run: (*parser).callonListNode14,
						expr: &seqExpr{
							pos: position{line: 82, col: 4, offset: 1648},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 82, col: 4, offset: 1648},
									name: "_",
								},
								&charClassMatcher{
									pos:        position{line: 82, col: 6, offset: 1650},
									val:        "[\\n]",
									chars:      []rune{'\n'},
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
			name:        "DocType",
			displayName: "\"doctype\"",
			pos:         position{line: 85, col: 1, offset: 1685},
			expr: &actionExpr{
				pos: position{line: 85, col: 21, offset: 1705},
				run: (*parser).callonDocType1,
				expr: &seqExpr{
					pos: position{line: 85, col: 21, offset: 1705},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 85, col: 21, offset: 1705},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 85, col: 23, offset: 1707},
							val:        "doctype",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 85, col: 33, offset: 1717},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 85, col: 35, offset: 1719},
							label: "val",
							expr: &ruleRefExpr{
								pos:  position{line: 85, col: 39, offset: 1723},
								name: "LineText",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 85, col: 48, offset: 1732},
							name: "NL",
						},
					},
				},
			},
		},
		{
			name: "Tag",
			pos:  position{line: 91, col: 1, offset: 1825},
			expr: &actionExpr{
				pos: position{line: 91, col: 8, offset: 1832},
				run: (*parser).callonTag1,
				expr: &seqExpr{
					pos: position{line: 91, col: 8, offset: 1832},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 91, col: 8, offset: 1832},
							label: "tag",
							expr: &ruleRefExpr{
								pos:  position{line: 91, col: 12, offset: 1836},
								name: "TagHeader",
							},
						},
						&labeledExpr{
							pos:   position{line: 91, col: 22, offset: 1846},
							label: "list",
							expr: &zeroOrOneExpr{
								pos: position{line: 91, col: 27, offset: 1851},
								expr: &ruleRefExpr{
									pos:  position{line: 91, col: 27, offset: 1851},
									name: "IndentedList",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "TagHeader",
			pos:  position{line: 101, col: 1, offset: 1974},
			expr: &choiceExpr{
				pos: position{line: 101, col: 14, offset: 1987},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 101, col: 14, offset: 1987},
						run: (*parser).callonTagHeader2,
						expr: &seqExpr{
							pos: position{line: 101, col: 14, offset: 1987},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 101, col: 14, offset: 1987},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 101, col: 16, offset: 1989},
									label: "name",
									expr: &ruleRefExpr{
										pos:  position{line: 101, col: 21, offset: 1994},
										name: "TagName",
									},
								},
								&labeledExpr{
									pos:   position{line: 101, col: 29, offset: 2002},
									label: "attrs",
									expr: &zeroOrOneExpr{
										pos: position{line: 101, col: 35, offset: 2008},
										expr: &ruleRefExpr{
											pos:  position{line: 101, col: 35, offset: 2008},
											name: "TagAttributes",
										},
									},
								},
								&labeledExpr{
									pos:   position{line: 101, col: 50, offset: 2023},
									label: "tl",
									expr: &zeroOrOneExpr{
										pos: position{line: 101, col: 53, offset: 2026},
										expr: &seqExpr{
											pos: position{line: 101, col: 54, offset: 2027},
											exprs: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 101, col: 54, offset: 2027},
													name: "__",
												},
												&zeroOrOneExpr{
													pos: position{line: 101, col: 57, offset: 2030},
													expr: &ruleRefExpr{
														pos:  position{line: 101, col: 57, offset: 2030},
														name: "TextList",
													},
												},
											},
										},
									},
								},
								&ruleRefExpr{
									pos:  position{line: 101, col: 69, offset: 2042},
									name: "NL",
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 110, col: 5, offset: 2265},
						run: (*parser).callonTagHeader17,
						expr: &seqExpr{
							pos: position{line: 110, col: 5, offset: 2265},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 110, col: 5, offset: 2265},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 110, col: 7, offset: 2267},
									label: "name",
									expr: &ruleRefExpr{
										pos:  position{line: 110, col: 12, offset: 2272},
										name: "TagName",
									},
								},
								&labeledExpr{
									pos:   position{line: 110, col: 20, offset: 2280},
									label: "attrs",
									expr: &zeroOrOneExpr{
										pos: position{line: 110, col: 26, offset: 2286},
										expr: &ruleRefExpr{
											pos:  position{line: 110, col: 26, offset: 2286},
											name: "TagAttributes",
										},
									},
								},
								&litMatcher{
									pos:        position{line: 110, col: 41, offset: 2301},
									val:        ".",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 110, col: 45, offset: 2305},
									name: "NL",
								},
								&labeledExpr{
									pos:   position{line: 110, col: 48, offset: 2308},
									label: "text",
									expr: &zeroOrOneExpr{
										pos: position{line: 110, col: 53, offset: 2313},
										expr: &ruleRefExpr{
											pos:  position{line: 110, col: 53, offset: 2313},
											name: "IndentedRawText",
										},
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 119, col: 5, offset: 2539},
						run: (*parser).callonTagHeader30,
						expr: &seqExpr{
							pos: position{line: 119, col: 5, offset: 2539},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 119, col: 5, offset: 2539},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 119, col: 7, offset: 2541},
									label: "name",
									expr: &ruleRefExpr{
										pos:  position{line: 119, col: 12, offset: 2546},
										name: "TagName",
									},
								},
								&labeledExpr{
									pos:   position{line: 119, col: 20, offset: 2554},
									label: "attrs",
									expr: &zeroOrOneExpr{
										pos: position{line: 119, col: 26, offset: 2560},
										expr: &ruleRefExpr{
											pos:  position{line: 119, col: 26, offset: 2560},
											name: "TagAttributes",
										},
									},
								},
								&labeledExpr{
									pos:   position{line: 119, col: 41, offset: 2575},
									label: "unescaped",
									expr: &zeroOrOneExpr{
										pos: position{line: 119, col: 51, offset: 2585},
										expr: &litMatcher{
											pos:        position{line: 119, col: 51, offset: 2585},
											val:        "!",
											ignoreCase: false,
										},
									},
								},
								&litMatcher{
									pos:        position{line: 119, col: 56, offset: 2590},
									val:        "=",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 119, col: 60, offset: 2594},
									name: "__",
								},
								&labeledExpr{
									pos:   position{line: 119, col: 63, offset: 2597},
									label: "expr",
									expr: &zeroOrOneExpr{
										pos: position{line: 119, col: 68, offset: 2602},
										expr: &ruleRefExpr{
											pos:  position{line: 119, col: 68, offset: 2602},
											name: "Expression",
										},
									},
								},
								&ruleRefExpr{
									pos:  position{line: 119, col: 80, offset: 2614},
									name: "NL",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "TagName",
			pos:  position{line: 134, col: 1, offset: 3019},
			expr: &choiceExpr{
				pos: position{line: 134, col: 12, offset: 3030},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 134, col: 12, offset: 3030},
						run: (*parser).callonTagName2,
						expr: &seqExpr{
							pos: position{line: 134, col: 12, offset: 3030},
							exprs: []interface{}{
								&charClassMatcher{
									pos:        position{line: 134, col: 12, offset: 3030},
									val:        "[_a-zA-Z]",
									chars:      []rune{'_'},
									ranges:     []rune{'a', 'z', 'A', 'Z'},
									ignoreCase: false,
									inverted:   false,
								},
								&zeroOrMoreExpr{
									pos: position{line: 134, col: 22, offset: 3040},
									expr: &charClassMatcher{
										pos:        position{line: 134, col: 22, offset: 3040},
										val:        "[_-:a-zA-Z0-9]",
										ranges:     []rune{'_', ':', 'a', 'z', 'A', 'Z', '0', '9'},
										ignoreCase: false,
										inverted:   false,
									},
								},
							},
						},
					},
					&andExpr{
						pos: position{line: 136, col: 5, offset: 3091},
						expr: &ruleRefExpr{
							pos:  position{line: 136, col: 6, offset: 3092},
							name: "TagAttributeClass",
						},
					},
					&actionExpr{
						pos: position{line: 136, col: 26, offset: 3112},
						run: (*parser).callonTagName9,
						expr: &andExpr{
							pos: position{line: 136, col: 26, offset: 3112},
							expr: &ruleRefExpr{
								pos:  position{line: 136, col: 27, offset: 3113},
								name: "TagAttributeID",
							},
						},
					},
				},
			},
		},
		{
			name: "TagAttributes",
			pos:  position{line: 140, col: 1, offset: 3153},
			expr: &choiceExpr{
				pos: position{line: 140, col: 18, offset: 3170},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 140, col: 18, offset: 3170},
						run: (*parser).callonTagAttributes2,
						expr: &seqExpr{
							pos: position{line: 140, col: 18, offset: 3170},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 140, col: 18, offset: 3170},
									label: "head",
									expr: &choiceExpr{
										pos: position{line: 140, col: 24, offset: 3176},
										alternatives: []interface{}{
											&ruleRefExpr{
												pos:  position{line: 140, col: 24, offset: 3176},
												name: "TagAttributeClass",
											},
											&ruleRefExpr{
												pos:  position{line: 140, col: 44, offset: 3196},
												name: "TagAttributeID",
											},
										},
									},
								},
								&labeledExpr{
									pos:   position{line: 140, col: 60, offset: 3212},
									label: "tail",
									expr: &zeroOrOneExpr{
										pos: position{line: 140, col: 65, offset: 3217},
										expr: &ruleRefExpr{
											pos:  position{line: 140, col: 65, offset: 3217},
											name: "TagAttributes",
										},
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 148, col: 5, offset: 3382},
						run: (*parser).callonTagAttributes11,
						expr: &seqExpr{
							pos: position{line: 148, col: 5, offset: 3382},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 148, col: 5, offset: 3382},
									label: "head",
									expr: &seqExpr{
										pos: position{line: 148, col: 11, offset: 3388},
										exprs: []interface{}{
											&litMatcher{
												pos:        position{line: 148, col: 11, offset: 3388},
												val:        "(",
												ignoreCase: false,
											},
											&ruleRefExpr{
												pos:  position{line: 148, col: 15, offset: 3392},
												name: "_",
											},
											&seqExpr{
												pos: position{line: 148, col: 18, offset: 3395},
												exprs: []interface{}{
													&ruleRefExpr{
														pos:  position{line: 148, col: 18, offset: 3395},
														name: "TagAttribute",
													},
													&zeroOrMoreExpr{
														pos: position{line: 148, col: 31, offset: 3408},
														expr: &seqExpr{
															pos: position{line: 148, col: 32, offset: 3409},
															exprs: []interface{}{
																&ruleRefExpr{
																	pos:  position{line: 148, col: 32, offset: 3409},
																	name: "__",
																},
																&ruleRefExpr{
																	pos:  position{line: 148, col: 35, offset: 3412},
																	name: "TagAttribute",
																},
															},
														},
													},
												},
											},
											&ruleRefExpr{
												pos:  position{line: 148, col: 51, offset: 3428},
												name: "_",
											},
											&litMatcher{
												pos:        position{line: 148, col: 53, offset: 3430},
												val:        ")",
												ignoreCase: false,
											},
										},
									},
								},
								&labeledExpr{
									pos:   position{line: 148, col: 58, offset: 3435},
									label: "tail",
									expr: &zeroOrOneExpr{
										pos: position{line: 148, col: 63, offset: 3440},
										expr: &ruleRefExpr{
											pos:  position{line: 148, col: 63, offset: 3440},
											name: "TagAttributes",
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
			name: "TagAttributeClass",
			pos:  position{line: 172, col: 1, offset: 3886},
			expr: &actionExpr{
				pos: position{line: 172, col: 22, offset: 3907},
				run: (*parser).callonTagAttributeClass1,
				expr: &seqExpr{
					pos: position{line: 172, col: 22, offset: 3907},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 172, col: 22, offset: 3907},
							val:        ".",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 172, col: 26, offset: 3911},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 172, col: 31, offset: 3916},
								name: "ClassName",
							},
						},
					},
				},
			},
		},
		{
			name: "TagAttributeID",
			pos:  position{line: 176, col: 1, offset: 4065},
			expr: &actionExpr{
				pos: position{line: 176, col: 19, offset: 4083},
				run: (*parser).callonTagAttributeID1,
				expr: &seqExpr{
					pos: position{line: 176, col: 19, offset: 4083},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 176, col: 19, offset: 4083},
							val:        "#",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 176, col: 23, offset: 4087},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 176, col: 28, offset: 4092},
								name: "IdName",
							},
						},
					},
				},
			},
		},
		{
			name: "TagAttribute",
			pos:  position{line: 180, col: 1, offset: 4235},
			expr: &choiceExpr{
				pos: position{line: 180, col: 17, offset: 4251},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 180, col: 17, offset: 4251},
						run: (*parser).callonTagAttribute2,
						expr: &seqExpr{
							pos: position{line: 180, col: 17, offset: 4251},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 180, col: 17, offset: 4251},
									label: "name",
									expr: &ruleRefExpr{
										pos:  position{line: 180, col: 22, offset: 4256},
										name: "TagAttributeName",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 180, col: 39, offset: 4273},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 180, col: 41, offset: 4275},
									val:        "=",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 180, col: 45, offset: 4279},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 180, col: 47, offset: 4281},
									label: "value",
									expr: &ruleRefExpr{
										pos:  position{line: 180, col: 53, offset: 4287},
										name: "Expression",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 182, col: 5, offset: 4423},
						run: (*parser).callonTagAttribute11,
						expr: &labeledExpr{
							pos:   position{line: 182, col: 5, offset: 4423},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 182, col: 10, offset: 4428},
								name: "TagAttributeName",
							},
						},
					},
				},
			},
		},
		{
			name: "TagAttributeName",
			pos:  position{line: 186, col: 1, offset: 4542},
			expr: &choiceExpr{
				pos: position{line: 186, col: 21, offset: 4562},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 186, col: 21, offset: 4562},
						run: (*parser).callonTagAttributeName2,
						expr: &seqExpr{
							pos: position{line: 186, col: 21, offset: 4562},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 186, col: 21, offset: 4562},
									val:        "(",
									ignoreCase: false,
								},
								&labeledExpr{
									pos:   position{line: 186, col: 25, offset: 4566},
									label: "tn",
									expr: &ruleRefExpr{
										pos:  position{line: 186, col: 28, offset: 4569},
										name: "TagAttributeNameLiteral",
									},
								},
								&litMatcher{
									pos:        position{line: 186, col: 52, offset: 4593},
									val:        ")",
									ignoreCase: false,
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 188, col: 5, offset: 4620},
						run: (*parser).callonTagAttributeName8,
						expr: &seqExpr{
							pos: position{line: 188, col: 5, offset: 4620},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 188, col: 5, offset: 4620},
									val:        "[",
									ignoreCase: false,
								},
								&labeledExpr{
									pos:   position{line: 188, col: 9, offset: 4624},
									label: "tn",
									expr: &ruleRefExpr{
										pos:  position{line: 188, col: 12, offset: 4627},
										name: "TagAttributeNameLiteral",
									},
								},
								&litMatcher{
									pos:        position{line: 188, col: 36, offset: 4651},
									val:        "]",
									ignoreCase: false,
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 190, col: 5, offset: 4678},
						run: (*parser).callonTagAttributeName14,
						expr: &labeledExpr{
							pos:   position{line: 190, col: 5, offset: 4678},
							label: "tn",
							expr: &ruleRefExpr{
								pos:  position{line: 190, col: 8, offset: 4681},
								name: "TagAttributeNameLiteral",
							},
						},
					},
				},
			},
		},
		{
			name: "ClassName",
			pos:  position{line: 194, col: 1, offset: 4727},
			expr: &actionExpr{
				pos: position{line: 194, col: 14, offset: 4740},
				run: (*parser).callonClassName1,
				expr: &oneOrMoreExpr{
					pos: position{line: 194, col: 14, offset: 4740},
					expr: &charClassMatcher{
						pos:        position{line: 194, col: 14, offset: 4740},
						val:        "[_-:a-zA-Z0-9]",
						ranges:     []rune{'_', ':', 'a', 'z', 'A', 'Z', '0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
				},
			},
		},
		{
			name: "IdName",
			pos:  position{line: 198, col: 1, offset: 4790},
			expr: &actionExpr{
				pos: position{line: 198, col: 11, offset: 4800},
				run: (*parser).callonIdName1,
				expr: &oneOrMoreExpr{
					pos: position{line: 198, col: 11, offset: 4800},
					expr: &charClassMatcher{
						pos:        position{line: 198, col: 11, offset: 4800},
						val:        "[_-:a-zA-Z0-9]",
						ranges:     []rune{'_', ':', 'a', 'z', 'A', 'Z', '0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
				},
			},
		},
		{
			name: "TagAttributeNameLiteral",
			pos:  position{line: 202, col: 1, offset: 4850},
			expr: &actionExpr{
				pos: position{line: 202, col: 28, offset: 4877},
				run: (*parser).callonTagAttributeNameLiteral1,
				expr: &seqExpr{
					pos: position{line: 202, col: 28, offset: 4877},
					exprs: []interface{}{
						&charClassMatcher{
							pos:        position{line: 202, col: 28, offset: 4877},
							val:        "[@_a-zA-Z]",
							chars:      []rune{'@', '_'},
							ranges:     []rune{'a', 'z', 'A', 'Z'},
							ignoreCase: false,
							inverted:   false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 202, col: 39, offset: 4888},
							expr: &charClassMatcher{
								pos:        position{line: 202, col: 39, offset: 4888},
								val:        "[._-:a-zA-Z0-9]",
								chars:      []rune{'.'},
								ranges:     []rune{'_', ':', 'a', 'z', 'A', 'Z', '0', '9'},
								ignoreCase: false,
								inverted:   false,
							},
						},
					},
				},
			},
		},
		{
			name: "If",
			pos:  position{line: 207, col: 1, offset: 4948},
			expr: &actionExpr{
				pos: position{line: 207, col: 7, offset: 4954},
				run: (*parser).callonIf1,
				expr: &seqExpr{
					pos: position{line: 207, col: 7, offset: 4954},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 207, col: 7, offset: 4954},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 207, col: 9, offset: 4956},
							val:        "if",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 207, col: 14, offset: 4961},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 207, col: 17, offset: 4964},
							label: "expr",
							expr: &ruleRefExpr{
								pos:  position{line: 207, col: 22, offset: 4969},
								name: "Expression",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 207, col: 33, offset: 4980},
							name: "_",
						},
						&ruleRefExpr{
							pos:  position{line: 207, col: 35, offset: 4982},
							name: "NL",
						},
						&labeledExpr{
							pos:   position{line: 207, col: 38, offset: 4985},
							label: "block",
							expr: &ruleRefExpr{
								pos:  position{line: 207, col: 44, offset: 4991},
								name: "IndentedList",
							},
						},
						&labeledExpr{
							pos:   position{line: 207, col: 57, offset: 5004},
							label: "elseNode",
							expr: &zeroOrOneExpr{
								pos: position{line: 207, col: 66, offset: 5013},
								expr: &ruleRefExpr{
									pos:  position{line: 207, col: 66, offset: 5013},
									name: "Else",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Else",
			pos:  position{line: 215, col: 1, offset: 5222},
			expr: &choiceExpr{
				pos: position{line: 215, col: 9, offset: 5230},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 215, col: 9, offset: 5230},
						run: (*parser).callonElse2,
						expr: &seqExpr{
							pos: position{line: 215, col: 9, offset: 5230},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 215, col: 9, offset: 5230},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 215, col: 11, offset: 5232},
									val:        "else",
									ignoreCase: false,
								},
								&labeledExpr{
									pos:   position{line: 215, col: 18, offset: 5239},
									label: "node",
									expr: &ruleRefExpr{
										pos:  position{line: 215, col: 23, offset: 5244},
										name: "If",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 217, col: 5, offset: 5272},
						run: (*parser).callonElse8,
						expr: &seqExpr{
							pos: position{line: 217, col: 5, offset: 5272},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 217, col: 5, offset: 5272},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 217, col: 7, offset: 5274},
									val:        "else",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 217, col: 14, offset: 5281},
									name: "_",
								},
								&ruleRefExpr{
									pos:  position{line: 217, col: 16, offset: 5283},
									name: "NL",
								},
								&labeledExpr{
									pos:   position{line: 217, col: 19, offset: 5286},
									label: "block",
									expr: &ruleRefExpr{
										pos:  position{line: 217, col: 25, offset: 5292},
										name: "IndentedList",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Each",
			pos:  position{line: 221, col: 1, offset: 5330},
			expr: &actionExpr{
				pos: position{line: 221, col: 9, offset: 5338},
				run: (*parser).callonEach1,
				expr: &seqExpr{
					pos: position{line: 221, col: 9, offset: 5338},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 221, col: 9, offset: 5338},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 221, col: 11, offset: 5340},
							val:        "each",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 221, col: 18, offset: 5347},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 221, col: 21, offset: 5350},
							label: "v1",
							expr: &ruleRefExpr{
								pos:  position{line: 221, col: 24, offset: 5353},
								name: "Variable",
							},
						},
						&labeledExpr{
							pos:   position{line: 221, col: 33, offset: 5362},
							label: "v2",
							expr: &zeroOrOneExpr{
								pos: position{line: 221, col: 36, offset: 5365},
								expr: &seqExpr{
									pos: position{line: 221, col: 37, offset: 5366},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 221, col: 37, offset: 5366},
											name: "_",
										},
										&litMatcher{
											pos:        position{line: 221, col: 39, offset: 5368},
											val:        ",",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 221, col: 43, offset: 5372},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 221, col: 45, offset: 5374},
											name: "Variable",
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 221, col: 56, offset: 5385},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 221, col: 58, offset: 5387},
							val:        "in",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 221, col: 63, offset: 5392},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 221, col: 65, offset: 5394},
							label: "expr",
							expr: &ruleRefExpr{
								pos:  position{line: 221, col: 70, offset: 5399},
								name: "Expression",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 221, col: 81, offset: 5410},
							name: "_",
						},
						&ruleRefExpr{
							pos:  position{line: 221, col: 83, offset: 5412},
							name: "NL",
						},
						&labeledExpr{
							pos:   position{line: 221, col: 86, offset: 5415},
							label: "block",
							expr: &ruleRefExpr{
								pos:  position{line: 221, col: 92, offset: 5421},
								name: "IndentedList",
							},
						},
					},
				},
			},
		},
		{
			name: "Assignment",
			pos:  position{line: 232, col: 1, offset: 5706},
			expr: &actionExpr{
				pos: position{line: 232, col: 15, offset: 5720},
				run: (*parser).callonAssignment1,
				expr: &seqExpr{
					pos: position{line: 232, col: 15, offset: 5720},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 232, col: 15, offset: 5720},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 232, col: 17, offset: 5722},
							val:        "-",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 232, col: 21, offset: 5726},
							name: "_",
						},
						&choiceExpr{
							pos: position{line: 232, col: 24, offset: 5729},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 232, col: 24, offset: 5729},
									val:        "var",
									ignoreCase: false,
								},
								&litMatcher{
									pos:        position{line: 232, col: 32, offset: 5737},
									val:        "let",
									ignoreCase: false,
								},
								&litMatcher{
									pos:        position{line: 232, col: 40, offset: 5745},
									val:        "const",
									ignoreCase: false,
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 232, col: 49, offset: 5754},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 232, col: 52, offset: 5757},
							label: "vr",
							expr: &ruleRefExpr{
								pos:  position{line: 232, col: 55, offset: 5760},
								name: "Variable",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 232, col: 64, offset: 5769},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 232, col: 66, offset: 5771},
							val:        "=",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 232, col: 70, offset: 5775},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 232, col: 72, offset: 5777},
							label: "expr",
							expr: &ruleRefExpr{
								pos:  position{line: 232, col: 77, offset: 5782},
								name: "Expression",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 232, col: 88, offset: 5793},
							name: "_",
						},
						&ruleRefExpr{
							pos:  position{line: 232, col: 90, offset: 5795},
							name: "NL",
						},
					},
				},
			},
		},
		{
			name: "Mixin",
			pos:  position{line: 237, col: 1, offset: 5927},
			expr: &actionExpr{
				pos: position{line: 237, col: 10, offset: 5936},
				run: (*parser).callonMixin1,
				expr: &seqExpr{
					pos: position{line: 237, col: 10, offset: 5936},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 237, col: 10, offset: 5936},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 237, col: 12, offset: 5938},
							val:        "mixin",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 237, col: 20, offset: 5946},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 237, col: 23, offset: 5949},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 237, col: 28, offset: 5954},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 237, col: 39, offset: 5965},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 237, col: 41, offset: 5967},
							label: "args",
							expr: &zeroOrOneExpr{
								pos: position{line: 237, col: 46, offset: 5972},
								expr: &ruleRefExpr{
									pos:  position{line: 237, col: 46, offset: 5972},
									name: "MixinArguments",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 237, col: 62, offset: 5988},
							name: "NL",
						},
						&labeledExpr{
							pos:   position{line: 237, col: 65, offset: 5991},
							label: "list",
							expr: &ruleRefExpr{
								pos:  position{line: 237, col: 70, offset: 5996},
								name: "IndentedList",
							},
						},
					},
				},
			},
		},
		{
			name: "MixinArguments",
			pos:  position{line: 245, col: 1, offset: 6205},
			expr: &choiceExpr{
				pos: position{line: 245, col: 19, offset: 6223},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 245, col: 19, offset: 6223},
						run: (*parser).callonMixinArguments2,
						expr: &seqExpr{
							pos: position{line: 245, col: 19, offset: 6223},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 245, col: 19, offset: 6223},
									val:        "(",
									ignoreCase: false,
								},
								&labeledExpr{
									pos:   position{line: 245, col: 23, offset: 6227},
									label: "head",
									expr: &ruleRefExpr{
										pos:  position{line: 245, col: 28, offset: 6232},
										name: "MixinArgument",
									},
								},
								&labeledExpr{
									pos:   position{line: 245, col: 42, offset: 6246},
									label: "tail",
									expr: &zeroOrMoreExpr{
										pos: position{line: 245, col: 47, offset: 6251},
										expr: &seqExpr{
											pos: position{line: 245, col: 48, offset: 6252},
											exprs: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 245, col: 48, offset: 6252},
													name: "_",
												},
												&litMatcher{
													pos:        position{line: 245, col: 50, offset: 6254},
													val:        ",",
													ignoreCase: false,
												},
												&ruleRefExpr{
													pos:  position{line: 245, col: 54, offset: 6258},
													name: "_",
												},
												&ruleRefExpr{
													pos:  position{line: 245, col: 56, offset: 6260},
													name: "MixinArgument",
												},
											},
										},
									},
								},
								&litMatcher{
									pos:        position{line: 245, col: 72, offset: 6276},
									val:        ")",
									ignoreCase: false,
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 258, col: 5, offset: 6538},
						run: (*parser).callonMixinArguments15,
						expr: &seqExpr{
							pos: position{line: 258, col: 5, offset: 6538},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 258, col: 5, offset: 6538},
									val:        "(",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 258, col: 9, offset: 6542},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 258, col: 11, offset: 6544},
									val:        ")",
									ignoreCase: false,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "MixinArgument",
			pos:  position{line: 262, col: 1, offset: 6571},
			expr: &actionExpr{
				pos: position{line: 262, col: 18, offset: 6588},
				run: (*parser).callonMixinArgument1,
				expr: &seqExpr{
					pos: position{line: 262, col: 18, offset: 6588},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 262, col: 18, offset: 6588},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 262, col: 23, offset: 6593},
								name: "Variable",
							},
						},
						&labeledExpr{
							pos:   position{line: 262, col: 32, offset: 6602},
							label: "def",
							expr: &zeroOrOneExpr{
								pos: position{line: 262, col: 36, offset: 6606},
								expr: &seqExpr{
									pos: position{line: 262, col: 37, offset: 6607},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 262, col: 37, offset: 6607},
											name: "_",
										},
										&litMatcher{
											pos:        position{line: 262, col: 39, offset: 6609},
											val:        "=",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 262, col: 43, offset: 6613},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 262, col: 45, offset: 6615},
											name: "Expression",
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
			name: "MixinCall",
			pos:  position{line: 273, col: 1, offset: 6839},
			expr: &actionExpr{
				pos: position{line: 273, col: 14, offset: 6852},
				run: (*parser).callonMixinCall1,
				expr: &seqExpr{
					pos: position{line: 273, col: 14, offset: 6852},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 273, col: 14, offset: 6852},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 273, col: 16, offset: 6854},
							val:        "+",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 273, col: 20, offset: 6858},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 273, col: 25, offset: 6863},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 273, col: 36, offset: 6874},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 273, col: 38, offset: 6876},
							label: "args",
							expr: &zeroOrOneExpr{
								pos: position{line: 273, col: 43, offset: 6881},
								expr: &ruleRefExpr{
									pos:  position{line: 273, col: 43, offset: 6881},
									name: "CallArguments",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 273, col: 58, offset: 6896},
							name: "NL",
						},
					},
				},
			},
		},
		{
			name: "CallArguments",
			pos:  position{line: 281, col: 1, offset: 7067},
			expr: &actionExpr{
				pos: position{line: 281, col: 18, offset: 7084},
				run: (*parser).callonCallArguments1,
				expr: &seqExpr{
					pos: position{line: 281, col: 18, offset: 7084},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 281, col: 18, offset: 7084},
							val:        "(",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 281, col: 22, offset: 7088},
							label: "head",
							expr: &zeroOrOneExpr{
								pos: position{line: 281, col: 27, offset: 7093},
								expr: &ruleRefExpr{
									pos:  position{line: 281, col: 27, offset: 7093},
									name: "Expression",
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 281, col: 39, offset: 7105},
							label: "tail",
							expr: &zeroOrMoreExpr{
								pos: position{line: 281, col: 44, offset: 7110},
								expr: &seqExpr{
									pos: position{line: 281, col: 45, offset: 7111},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 281, col: 45, offset: 7111},
											name: "_",
										},
										&litMatcher{
											pos:        position{line: 281, col: 47, offset: 7113},
											val:        ",",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 281, col: 51, offset: 7117},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 281, col: 53, offset: 7119},
											name: "Expression",
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 281, col: 66, offset: 7132},
							val:        ")",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Import",
			pos:  position{line: 302, col: 1, offset: 7454},
			expr: &actionExpr{
				pos: position{line: 302, col: 11, offset: 7464},
				run: (*parser).callonImport1,
				expr: &seqExpr{
					pos: position{line: 302, col: 11, offset: 7464},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 302, col: 11, offset: 7464},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 302, col: 13, offset: 7466},
							val:        "include",
							ignoreCase: false,
						},
						&zeroOrOneExpr{
							pos: position{line: 302, col: 23, offset: 7476},
							expr: &litMatcher{
								pos:        position{line: 302, col: 23, offset: 7476},
								val:        "s",
								ignoreCase: false,
							},
						},
						&ruleRefExpr{
							pos:  position{line: 302, col: 28, offset: 7481},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 302, col: 31, offset: 7484},
							label: "file",
							expr: &ruleRefExpr{
								pos:  position{line: 302, col: 36, offset: 7489},
								name: "String",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 302, col: 43, offset: 7496},
							name: "NL",
						},
					},
				},
			},
		},
		{
			name: "Extend",
			pos:  position{line: 306, col: 1, offset: 7579},
			expr: &actionExpr{
				pos: position{line: 306, col: 11, offset: 7589},
				run: (*parser).callonExtend1,
				expr: &seqExpr{
					pos: position{line: 306, col: 11, offset: 7589},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 306, col: 11, offset: 7589},
							val:        "extend",
							ignoreCase: false,
						},
						&zeroOrOneExpr{
							pos: position{line: 306, col: 20, offset: 7598},
							expr: &litMatcher{
								pos:        position{line: 306, col: 20, offset: 7598},
								val:        "s",
								ignoreCase: false,
							},
						},
						&ruleRefExpr{
							pos:  position{line: 306, col: 25, offset: 7603},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 306, col: 28, offset: 7606},
							label: "file",
							expr: &ruleRefExpr{
								pos:  position{line: 306, col: 33, offset: 7611},
								name: "String",
							},
						},
					},
				},
			},
		},
		{
			name: "Block",
			pos:  position{line: 310, col: 1, offset: 7698},
			expr: &actionExpr{
				pos: position{line: 310, col: 10, offset: 7707},
				run: (*parser).callonBlock1,
				expr: &seqExpr{
					pos: position{line: 310, col: 10, offset: 7707},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 310, col: 10, offset: 7707},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 310, col: 12, offset: 7709},
							val:        "block",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 310, col: 20, offset: 7717},
							label: "mod",
							expr: &zeroOrOneExpr{
								pos: position{line: 310, col: 24, offset: 7721},
								expr: &seqExpr{
									pos: position{line: 310, col: 25, offset: 7722},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 310, col: 25, offset: 7722},
											name: "__",
										},
										&choiceExpr{
											pos: position{line: 310, col: 29, offset: 7726},
											alternatives: []interface{}{
												&litMatcher{
													pos:        position{line: 310, col: 29, offset: 7726},
													val:        "append",
													ignoreCase: false,
												},
												&litMatcher{
													pos:        position{line: 310, col: 40, offset: 7737},
													val:        "prepend",
													ignoreCase: false,
												},
											},
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 310, col: 53, offset: 7750},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 310, col: 56, offset: 7753},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 310, col: 61, offset: 7758},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 310, col: 72, offset: 7769},
							name: "NL",
						},
						&labeledExpr{
							pos:   position{line: 310, col: 75, offset: 7772},
							label: "list",
							expr: &zeroOrOneExpr{
								pos: position{line: 310, col: 80, offset: 7777},
								expr: &ruleRefExpr{
									pos:  position{line: 310, col: 80, offset: 7777},
									name: "IndentedList",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Comment",
			pos:  position{line: 331, col: 1, offset: 8142},
			expr: &actionExpr{
				pos: position{line: 331, col: 12, offset: 8153},
				run: (*parser).callonComment1,
				expr: &seqExpr{
					pos: position{line: 331, col: 12, offset: 8153},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 331, col: 12, offset: 8153},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 331, col: 14, offset: 8155},
							val:        "//",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 331, col: 19, offset: 8160},
							label: "silent",
							expr: &zeroOrOneExpr{
								pos: position{line: 331, col: 26, offset: 8167},
								expr: &litMatcher{
									pos:        position{line: 331, col: 26, offset: 8167},
									val:        "-",
									ignoreCase: false,
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 331, col: 31, offset: 8172},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 331, col: 33, offset: 8174},
							label: "comment",
							expr: &ruleRefExpr{
								pos:  position{line: 331, col: 41, offset: 8182},
								name: "LineText",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 331, col: 50, offset: 8191},
							name: "NL",
						},
					},
				},
			},
		},
		{
			name: "LineText",
			pos:  position{line: 336, col: 1, offset: 8325},
			expr: &actionExpr{
				pos: position{line: 336, col: 13, offset: 8337},
				run: (*parser).callonLineText1,
				expr: &zeroOrMoreExpr{
					pos: position{line: 336, col: 13, offset: 8337},
					expr: &charClassMatcher{
						pos:        position{line: 336, col: 13, offset: 8337},
						val:        "[^\\n]",
						chars:      []rune{'\n'},
						ignoreCase: false,
						inverted:   true,
					},
				},
			},
		},
		{
			name: "PipeText",
			pos:  position{line: 341, col: 1, offset: 8386},
			expr: &actionExpr{
				pos: position{line: 341, col: 13, offset: 8398},
				run: (*parser).callonPipeText1,
				expr: &seqExpr{
					pos: position{line: 341, col: 13, offset: 8398},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 341, col: 13, offset: 8398},
							name: "_",
						},
						&choiceExpr{
							pos: position{line: 341, col: 16, offset: 8401},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 341, col: 16, offset: 8401},
									val:        "|",
									ignoreCase: false,
								},
								&litMatcher{
									pos:        position{line: 341, col: 22, offset: 8407},
									val:        "<",
									ignoreCase: false,
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 341, col: 27, offset: 8412},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 341, col: 29, offset: 8414},
							label: "tl",
							expr: &ruleRefExpr{
								pos:  position{line: 341, col: 32, offset: 8417},
								name: "TextList",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 341, col: 41, offset: 8426},
							name: "NL",
						},
					},
				},
			},
		},
		{
			name: "TextList",
			pos:  position{line: 345, col: 1, offset: 8451},
			expr: &choiceExpr{
				pos: position{line: 345, col: 13, offset: 8463},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 345, col: 13, offset: 8463},
						run: (*parser).callonTextList2,
						expr: &seqExpr{
							pos: position{line: 345, col: 13, offset: 8463},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 345, col: 13, offset: 8463},
									label: "intr",
									expr: &ruleRefExpr{
										pos:  position{line: 345, col: 18, offset: 8468},
										name: "Interpolation",
									},
								},
								&labeledExpr{
									pos:   position{line: 345, col: 32, offset: 8482},
									label: "tl",
									expr: &ruleRefExpr{
										pos:  position{line: 345, col: 35, offset: 8485},
										name: "TextList",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 360, col: 5, offset: 8806},
						run: (*parser).callonTextList8,
						expr: &andExpr{
							pos: position{line: 360, col: 5, offset: 8806},
							expr: &ruleRefExpr{
								pos:  position{line: 360, col: 6, offset: 8807},
								name: "NL",
							},
						},
					},
					&actionExpr{
						pos: position{line: 362, col: 5, offset: 8872},
						run: (*parser).callonTextList11,
						expr: &seqExpr{
							pos: position{line: 362, col: 5, offset: 8872},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 362, col: 5, offset: 8872},
									label: "ch",
									expr: &anyMatcher{
										line: 362, col: 8, offset: 8875,
									},
								},
								&labeledExpr{
									pos:   position{line: 362, col: 10, offset: 8877},
									label: "tl",
									expr: &ruleRefExpr{
										pos:  position{line: 362, col: 13, offset: 8880},
										name: "TextList",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Indent",
			pos:  position{line: 379, col: 1, offset: 9276},
			expr: &litMatcher{
				pos:        position{line: 379, col: 11, offset: 9286},
				val:        "\x01",
				ignoreCase: false,
			},
		},
		{
			name: "Outdent",
			pos:  position{line: 380, col: 1, offset: 9295},
			expr: &litMatcher{
				pos:        position{line: 380, col: 12, offset: 9306},
				val:        "\x02",
				ignoreCase: false,
			},
		},
		{
			name: "Interpolation",
			pos:  position{line: 382, col: 1, offset: 9316},
			expr: &actionExpr{
				pos: position{line: 382, col: 18, offset: 9333},
				run: (*parser).callonInterpolation1,
				expr: &seqExpr{
					pos: position{line: 382, col: 18, offset: 9333},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 382, col: 18, offset: 9333},
							label: "mod",
							expr: &choiceExpr{
								pos: position{line: 382, col: 23, offset: 9338},
								alternatives: []interface{}{
									&litMatcher{
										pos:        position{line: 382, col: 23, offset: 9338},
										val:        "#",
										ignoreCase: false,
									},
									&litMatcher{
										pos:        position{line: 382, col: 29, offset: 9344},
										val:        "!",
										ignoreCase: false,
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 382, col: 34, offset: 9349},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 382, col: 38, offset: 9353},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 382, col: 40, offset: 9355},
							label: "expr",
							expr: &ruleRefExpr{
								pos:  position{line: 382, col: 45, offset: 9360},
								name: "Expression",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 382, col: 56, offset: 9371},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 382, col: 58, offset: 9373},
							val:        "}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Expression",
			pos:  position{line: 390, col: 1, offset: 9557},
			expr: &ruleRefExpr{
				pos:  position{line: 390, col: 15, offset: 9571},
				name: "ExpressionTernery",
			},
		},
		{
			name: "ExpressionTernery",
			pos:  position{line: 392, col: 1, offset: 9590},
			expr: &actionExpr{
				pos: position{line: 392, col: 22, offset: 9611},
				run: (*parser).callonExpressionTernery1,
				expr: &seqExpr{
					pos: position{line: 392, col: 22, offset: 9611},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 392, col: 22, offset: 9611},
							label: "cnd",
							expr: &ruleRefExpr{
								pos:  position{line: 392, col: 26, offset: 9615},
								name: "ExpressionBinOp",
							},
						},
						&labeledExpr{
							pos:   position{line: 392, col: 42, offset: 9631},
							label: "rest",
							expr: &zeroOrOneExpr{
								pos: position{line: 392, col: 47, offset: 9636},
								expr: &seqExpr{
									pos: position{line: 392, col: 48, offset: 9637},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 392, col: 48, offset: 9637},
											name: "_",
										},
										&litMatcher{
											pos:        position{line: 392, col: 50, offset: 9639},
											val:        "?",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 392, col: 54, offset: 9643},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 392, col: 56, offset: 9645},
											name: "ExpressionTernery",
										},
										&ruleRefExpr{
											pos:  position{line: 392, col: 74, offset: 9663},
											name: "_",
										},
										&litMatcher{
											pos:        position{line: 392, col: 76, offset: 9665},
											val:        ":",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 392, col: 80, offset: 9669},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 392, col: 82, offset: 9671},
											name: "ExpressionTernery",
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
			name: "ExpressionBinOp",
			pos:  position{line: 412, col: 1, offset: 10043},
			expr: &actionExpr{
				pos: position{line: 412, col: 20, offset: 10062},
				run: (*parser).callonExpressionBinOp1,
				expr: &seqExpr{
					pos: position{line: 412, col: 20, offset: 10062},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 412, col: 20, offset: 10062},
							label: "first",
							expr: &ruleRefExpr{
								pos:  position{line: 412, col: 26, offset: 10068},
								name: "ExpressionCmpOp",
							},
						},
						&labeledExpr{
							pos:   position{line: 412, col: 42, offset: 10084},
							label: "rest",
							expr: &zeroOrMoreExpr{
								pos: position{line: 412, col: 47, offset: 10089},
								expr: &seqExpr{
									pos: position{line: 412, col: 49, offset: 10091},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 412, col: 49, offset: 10091},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 412, col: 51, offset: 10093},
											name: "CmpOp",
										},
										&ruleRefExpr{
											pos:  position{line: 412, col: 57, offset: 10099},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 412, col: 59, offset: 10101},
											name: "ExpressionBinOp",
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
			name: "ExpressionCmpOp",
			pos:  position{line: 416, col: 1, offset: 10161},
			expr: &actionExpr{
				pos: position{line: 416, col: 20, offset: 10180},
				run: (*parser).callonExpressionCmpOp1,
				expr: &seqExpr{
					pos: position{line: 416, col: 20, offset: 10180},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 416, col: 20, offset: 10180},
							label: "first",
							expr: &ruleRefExpr{
								pos:  position{line: 416, col: 26, offset: 10186},
								name: "ExpressionAddOp",
							},
						},
						&labeledExpr{
							pos:   position{line: 416, col: 42, offset: 10202},
							label: "rest",
							expr: &zeroOrMoreExpr{
								pos: position{line: 416, col: 47, offset: 10207},
								expr: &seqExpr{
									pos: position{line: 416, col: 49, offset: 10209},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 416, col: 49, offset: 10209},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 416, col: 51, offset: 10211},
											name: "CmpOp",
										},
										&ruleRefExpr{
											pos:  position{line: 416, col: 57, offset: 10217},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 416, col: 59, offset: 10219},
											name: "ExpressionCmpOp",
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
			name: "ExpressionAddOp",
			pos:  position{line: 420, col: 1, offset: 10279},
			expr: &actionExpr{
				pos: position{line: 420, col: 20, offset: 10298},
				run: (*parser).callonExpressionAddOp1,
				expr: &seqExpr{
					pos: position{line: 420, col: 20, offset: 10298},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 420, col: 20, offset: 10298},
							label: "first",
							expr: &ruleRefExpr{
								pos:  position{line: 420, col: 26, offset: 10304},
								name: "ExpressionMulOp",
							},
						},
						&labeledExpr{
							pos:   position{line: 420, col: 42, offset: 10320},
							label: "rest",
							expr: &zeroOrMoreExpr{
								pos: position{line: 420, col: 47, offset: 10325},
								expr: &seqExpr{
									pos: position{line: 420, col: 49, offset: 10327},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 420, col: 49, offset: 10327},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 420, col: 51, offset: 10329},
											name: "AddOp",
										},
										&ruleRefExpr{
											pos:  position{line: 420, col: 57, offset: 10335},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 420, col: 59, offset: 10337},
											name: "ExpressionAddOp",
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
			name: "ExpressionMulOp",
			pos:  position{line: 424, col: 1, offset: 10397},
			expr: &actionExpr{
				pos: position{line: 424, col: 20, offset: 10416},
				run: (*parser).callonExpressionMulOp1,
				expr: &seqExpr{
					pos: position{line: 424, col: 20, offset: 10416},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 424, col: 20, offset: 10416},
							label: "first",
							expr: &ruleRefExpr{
								pos:  position{line: 424, col: 26, offset: 10422},
								name: "ExpressionUnaryOp",
							},
						},
						&labeledExpr{
							pos:   position{line: 424, col: 44, offset: 10440},
							label: "rest",
							expr: &zeroOrMoreExpr{
								pos: position{line: 424, col: 49, offset: 10445},
								expr: &seqExpr{
									pos: position{line: 424, col: 51, offset: 10447},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 424, col: 51, offset: 10447},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 424, col: 53, offset: 10449},
											name: "MulOp",
										},
										&ruleRefExpr{
											pos:  position{line: 424, col: 59, offset: 10455},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 424, col: 61, offset: 10457},
											name: "ExpressionMulOp",
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
			name: "ExpressionUnaryOp",
			pos:  position{line: 428, col: 1, offset: 10517},
			expr: &choiceExpr{
				pos: position{line: 428, col: 22, offset: 10538},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 428, col: 22, offset: 10538},
						run: (*parser).callonExpressionUnaryOp2,
						expr: &seqExpr{
							pos: position{line: 428, col: 22, offset: 10538},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 428, col: 22, offset: 10538},
									label: "op",
									expr: &ruleRefExpr{
										pos:  position{line: 428, col: 25, offset: 10541},
										name: "UnaryOp",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 428, col: 33, offset: 10549},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 428, col: 35, offset: 10551},
									label: "ex",
									expr: &ruleRefExpr{
										pos:  position{line: 428, col: 38, offset: 10554},
										name: "ExpressionFactor",
									},
								},
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 430, col: 5, offset: 10677},
						name: "ExpressionFactor",
					},
				},
			},
		},
		{
			name: "ExpressionFactor",
			pos:  position{line: 432, col: 1, offset: 10695},
			expr: &choiceExpr{
				pos: position{line: 432, col: 21, offset: 10715},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 432, col: 21, offset: 10715},
						run: (*parser).callonExpressionFactor2,
						expr: &seqExpr{
							pos: position{line: 432, col: 21, offset: 10715},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 432, col: 21, offset: 10715},
									val:        "(",
									ignoreCase: false,
								},
								&labeledExpr{
									pos:   position{line: 432, col: 25, offset: 10719},
									label: "e",
									expr: &ruleRefExpr{
										pos:  position{line: 432, col: 27, offset: 10721},
										name: "Expression",
									},
								},
								&litMatcher{
									pos:        position{line: 432, col: 38, offset: 10732},
									val:        ")",
									ignoreCase: false,
								},
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 434, col: 5, offset: 10758},
						name: "StringExpression",
					},
					&ruleRefExpr{
						pos:  position{line: 434, col: 24, offset: 10777},
						name: "NumberExpression",
					},
					&ruleRefExpr{
						pos:  position{line: 434, col: 43, offset: 10796},
						name: "BooleanExpression",
					},
					&ruleRefExpr{
						pos:  position{line: 434, col: 63, offset: 10816},
						name: "NilExpression",
					},
					&ruleRefExpr{
						pos:  position{line: 434, col: 79, offset: 10832},
						name: "MemberExpression",
					},
				},
			},
		},
		{
			name: "StringExpression",
			pos:  position{line: 436, col: 1, offset: 10850},
			expr: &actionExpr{
				pos: position{line: 436, col: 21, offset: 10870},
				run: (*parser).callonStringExpression1,
				expr: &labeledExpr{
					pos:   position{line: 436, col: 21, offset: 10870},
					label: "s",
					expr: &ruleRefExpr{
						pos:  position{line: 436, col: 23, offset: 10872},
						name: "String",
					},
				},
			},
		},
		{
			name: "NumberExpression",
			pos:  position{line: 440, col: 1, offset: 10967},
			expr: &actionExpr{
				pos: position{line: 440, col: 21, offset: 10987},
				run: (*parser).callonNumberExpression1,
				expr: &seqExpr{
					pos: position{line: 440, col: 21, offset: 10987},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 440, col: 21, offset: 10987},
							expr: &litMatcher{
								pos:        position{line: 440, col: 21, offset: 10987},
								val:        "-",
								ignoreCase: false,
							},
						},
						&ruleRefExpr{
							pos:  position{line: 440, col: 26, offset: 10992},
							name: "Integer",
						},
						&labeledExpr{
							pos:   position{line: 440, col: 34, offset: 11000},
							label: "dec",
							expr: &zeroOrOneExpr{
								pos: position{line: 440, col: 38, offset: 11004},
								expr: &seqExpr{
									pos: position{line: 440, col: 40, offset: 11006},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 440, col: 40, offset: 11006},
											val:        ".",
											ignoreCase: false,
										},
										&oneOrMoreExpr{
											pos: position{line: 440, col: 44, offset: 11010},
											expr: &ruleRefExpr{
												pos:  position{line: 440, col: 44, offset: 11010},
												name: "DecimalDigit",
											},
										},
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 440, col: 61, offset: 11027},
							label: "ex",
							expr: &zeroOrOneExpr{
								pos: position{line: 440, col: 64, offset: 11030},
								expr: &ruleRefExpr{
									pos:  position{line: 440, col: 64, offset: 11030},
									name: "Exponent",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "NilExpression",
			pos:  position{line: 450, col: 1, offset: 11360},
			expr: &actionExpr{
				pos: position{line: 450, col: 18, offset: 11377},
				run: (*parser).callonNilExpression1,
				expr: &ruleRefExpr{
					pos:  position{line: 450, col: 18, offset: 11377},
					name: "Null",
				},
			},
		},
		{
			name: "BooleanExpression",
			pos:  position{line: 454, col: 1, offset: 11448},
			expr: &actionExpr{
				pos: position{line: 454, col: 22, offset: 11469},
				run: (*parser).callonBooleanExpression1,
				expr: &labeledExpr{
					pos:   position{line: 454, col: 22, offset: 11469},
					label: "b",
					expr: &ruleRefExpr{
						pos:  position{line: 454, col: 24, offset: 11471},
						name: "Bool",
					},
				},
			},
		},
		{
			name: "MemberExpression",
			pos:  position{line: 458, col: 1, offset: 11563},
			expr: &actionExpr{
				pos: position{line: 458, col: 21, offset: 11583},
				run: (*parser).callonMemberExpression1,
				expr: &seqExpr{
					pos: position{line: 458, col: 21, offset: 11583},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 458, col: 21, offset: 11583},
							label: "field",
							expr: &ruleRefExpr{
								pos:  position{line: 458, col: 27, offset: 11589},
								name: "Field",
							},
						},
						&labeledExpr{
							pos:   position{line: 458, col: 33, offset: 11595},
							label: "member",
							expr: &zeroOrMoreExpr{
								pos: position{line: 458, col: 40, offset: 11602},
								expr: &choiceExpr{
									pos: position{line: 458, col: 41, offset: 11603},
									alternatives: []interface{}{
										&seqExpr{
											pos: position{line: 458, col: 42, offset: 11604},
											exprs: []interface{}{
												&litMatcher{
													pos:        position{line: 458, col: 42, offset: 11604},
													val:        ".",
													ignoreCase: false,
												},
												&ruleRefExpr{
													pos:  position{line: 458, col: 46, offset: 11608},
													name: "Identifier",
												},
											},
										},
										&seqExpr{
											pos: position{line: 458, col: 61, offset: 11623},
											exprs: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 458, col: 61, offset: 11623},
													name: "_",
												},
												&ruleRefExpr{
													pos:  position{line: 458, col: 63, offset: 11625},
													name: "CallArguments",
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
			name: "Field",
			pos:  position{line: 475, col: 1, offset: 12083},
			expr: &actionExpr{
				pos: position{line: 475, col: 10, offset: 12092},
				run: (*parser).callonField1,
				expr: &labeledExpr{
					pos:   position{line: 475, col: 10, offset: 12092},
					label: "variable",
					expr: &ruleRefExpr{
						pos:  position{line: 475, col: 19, offset: 12101},
						name: "Variable",
					},
				},
			},
		},
		{
			name: "UnaryOp",
			pos:  position{line: 479, col: 1, offset: 12210},
			expr: &actionExpr{
				pos: position{line: 479, col: 12, offset: 12221},
				run: (*parser).callonUnaryOp1,
				expr: &choiceExpr{
					pos: position{line: 479, col: 14, offset: 12223},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 479, col: 14, offset: 12223},
							val:        "+",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 479, col: 20, offset: 12229},
							val:        "-",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 479, col: 26, offset: 12235},
							val:        "!",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "AddOp",
			pos:  position{line: 483, col: 1, offset: 12275},
			expr: &actionExpr{
				pos: position{line: 483, col: 10, offset: 12284},
				run: (*parser).callonAddOp1,
				expr: &choiceExpr{
					pos: position{line: 483, col: 12, offset: 12286},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 483, col: 12, offset: 12286},
							val:        "+",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 483, col: 18, offset: 12292},
							val:        "-",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "MulOp",
			pos:  position{line: 487, col: 1, offset: 12332},
			expr: &actionExpr{
				pos: position{line: 487, col: 10, offset: 12341},
				run: (*parser).callonMulOp1,
				expr: &choiceExpr{
					pos: position{line: 487, col: 12, offset: 12343},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 487, col: 12, offset: 12343},
							val:        "*",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 487, col: 18, offset: 12349},
							val:        "/",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 487, col: 24, offset: 12355},
							val:        "%",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "CmpOp",
			pos:  position{line: 491, col: 1, offset: 12395},
			expr: &actionExpr{
				pos: position{line: 491, col: 10, offset: 12404},
				run: (*parser).callonCmpOp1,
				expr: &choiceExpr{
					pos: position{line: 491, col: 12, offset: 12406},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 491, col: 12, offset: 12406},
							val:        "==",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 491, col: 19, offset: 12413},
							val:        "!=",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 491, col: 26, offset: 12420},
							val:        "<",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 491, col: 32, offset: 12426},
							val:        ">",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 491, col: 38, offset: 12432},
							val:        "<=",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 491, col: 45, offset: 12439},
							val:        ">=",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "BinOp",
			pos:  position{line: 495, col: 1, offset: 12480},
			expr: &actionExpr{
				pos: position{line: 495, col: 10, offset: 12489},
				run: (*parser).callonBinOp1,
				expr: &choiceExpr{
					pos: position{line: 495, col: 12, offset: 12491},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 495, col: 12, offset: 12491},
							val:        "&&",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 495, col: 19, offset: 12498},
							val:        "||",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name:        "String",
			displayName: "\"string\"",
			pos:         position{line: 499, col: 1, offset: 12539},
			expr: &actionExpr{
				pos: position{line: 499, col: 20, offset: 12558},
				run: (*parser).callonString1,
				expr: &seqExpr{
					pos: position{line: 499, col: 20, offset: 12558},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 499, col: 20, offset: 12558},
							name: "Quote",
						},
						&zeroOrMoreExpr{
							pos: position{line: 499, col: 26, offset: 12564},
							expr: &choiceExpr{
								pos: position{line: 499, col: 28, offset: 12566},
								alternatives: []interface{}{
									&seqExpr{
										pos: position{line: 499, col: 28, offset: 12566},
										exprs: []interface{}{
											&notExpr{
												pos: position{line: 499, col: 28, offset: 12566},
												expr: &ruleRefExpr{
													pos:  position{line: 499, col: 29, offset: 12567},
													name: "EscapedChar",
												},
											},
											&anyMatcher{
												line: 499, col: 41, offset: 12579,
											},
										},
									},
									&seqExpr{
										pos: position{line: 499, col: 45, offset: 12583},
										exprs: []interface{}{
											&litMatcher{
												pos:        position{line: 499, col: 45, offset: 12583},
												val:        "\\",
												ignoreCase: false,
											},
											&ruleRefExpr{
												pos:  position{line: 499, col: 50, offset: 12588},
												name: "EscapeSequence",
											},
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 499, col: 68, offset: 12606},
							name: "Quote",
						},
					},
				},
			},
		},
		{
			name:        "Quote",
			displayName: "\"quote\"",
			pos:         position{line: 503, col: 1, offset: 12658},
			expr: &litMatcher{
				pos:        position{line: 503, col: 18, offset: 12675},
				val:        "\"",
				ignoreCase: false,
			},
		},
		{
			name: "EscapedChar",
			pos:  position{line: 505, col: 1, offset: 12680},
			expr: &charClassMatcher{
				pos:        position{line: 505, col: 16, offset: 12695},
				val:        "[\\x00-\\x1f\"\\\\]",
				chars:      []rune{'"', '\\'},
				ranges:     []rune{'\x00', '\x1f'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "EscapeSequence",
			pos:  position{line: 506, col: 1, offset: 12710},
			expr: &choiceExpr{
				pos: position{line: 506, col: 19, offset: 12728},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 506, col: 19, offset: 12728},
						name: "SingleCharEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 506, col: 38, offset: 12747},
						name: "UnicodeEscape",
					},
				},
			},
		},
		{
			name: "SingleCharEscape",
			pos:  position{line: 507, col: 1, offset: 12761},
			expr: &charClassMatcher{
				pos:        position{line: 507, col: 21, offset: 12781},
				val:        "[\"\\\\/bfnrt]",
				chars:      []rune{'"', '\\', '/', 'b', 'f', 'n', 'r', 't'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "UnicodeEscape",
			pos:  position{line: 508, col: 1, offset: 12793},
			expr: &seqExpr{
				pos: position{line: 508, col: 18, offset: 12810},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 508, col: 18, offset: 12810},
						val:        "u",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 508, col: 22, offset: 12814},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 508, col: 31, offset: 12823},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 508, col: 40, offset: 12832},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 508, col: 49, offset: 12841},
						name: "HexDigit",
					},
				},
			},
		},
		{
			name: "Integer",
			pos:  position{line: 510, col: 1, offset: 12851},
			expr: &choiceExpr{
				pos: position{line: 510, col: 12, offset: 12862},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 510, col: 12, offset: 12862},
						val:        "0",
						ignoreCase: false,
					},
					&seqExpr{
						pos: position{line: 510, col: 18, offset: 12868},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 510, col: 18, offset: 12868},
								name: "NonZeroDecimalDigit",
							},
							&zeroOrMoreExpr{
								pos: position{line: 510, col: 38, offset: 12888},
								expr: &ruleRefExpr{
									pos:  position{line: 510, col: 38, offset: 12888},
									name: "DecimalDigit",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Exponent",
			pos:  position{line: 511, col: 1, offset: 12902},
			expr: &seqExpr{
				pos: position{line: 511, col: 13, offset: 12914},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 511, col: 13, offset: 12914},
						val:        "e",
						ignoreCase: true,
					},
					&zeroOrOneExpr{
						pos: position{line: 511, col: 18, offset: 12919},
						expr: &charClassMatcher{
							pos:        position{line: 511, col: 18, offset: 12919},
							val:        "[+-]",
							chars:      []rune{'+', '-'},
							ignoreCase: false,
							inverted:   false,
						},
					},
					&oneOrMoreExpr{
						pos: position{line: 511, col: 24, offset: 12925},
						expr: &ruleRefExpr{
							pos:  position{line: 511, col: 24, offset: 12925},
							name: "DecimalDigit",
						},
					},
				},
			},
		},
		{
			name: "DecimalDigit",
			pos:  position{line: 512, col: 1, offset: 12939},
			expr: &charClassMatcher{
				pos:        position{line: 512, col: 17, offset: 12955},
				val:        "[0-9]",
				ranges:     []rune{'0', '9'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "NonZeroDecimalDigit",
			pos:  position{line: 513, col: 1, offset: 12961},
			expr: &charClassMatcher{
				pos:        position{line: 513, col: 24, offset: 12984},
				val:        "[1-9]",
				ranges:     []rune{'1', '9'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "HexDigit",
			pos:  position{line: 514, col: 1, offset: 12990},
			expr: &charClassMatcher{
				pos:        position{line: 514, col: 13, offset: 13002},
				val:        "[0-9a-f]i",
				ranges:     []rune{'0', '9', 'a', 'f'},
				ignoreCase: true,
				inverted:   false,
			},
		},
		{
			name: "Bool",
			pos:  position{line: 515, col: 1, offset: 13012},
			expr: &choiceExpr{
				pos: position{line: 515, col: 9, offset: 13020},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 515, col: 9, offset: 13020},
						run: (*parser).callonBool2,
						expr: &litMatcher{
							pos:        position{line: 515, col: 9, offset: 13020},
							val:        "true",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 515, col: 39, offset: 13050},
						run: (*parser).callonBool4,
						expr: &litMatcher{
							pos:        position{line: 515, col: 39, offset: 13050},
							val:        "false",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Null",
			pos:  position{line: 516, col: 1, offset: 13080},
			expr: &actionExpr{
				pos: position{line: 516, col: 9, offset: 13088},
				run: (*parser).callonNull1,
				expr: &choiceExpr{
					pos: position{line: 516, col: 10, offset: 13089},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 516, col: 10, offset: 13089},
							val:        "null",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 516, col: 19, offset: 13098},
							val:        "nil",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Variable",
			pos:  position{line: 518, col: 1, offset: 13126},
			expr: &actionExpr{
				pos: position{line: 518, col: 13, offset: 13138},
				run: (*parser).callonVariable1,
				expr: &labeledExpr{
					pos:   position{line: 518, col: 13, offset: 13138},
					label: "ident",
					expr: &ruleRefExpr{
						pos:  position{line: 518, col: 19, offset: 13144},
						name: "Identifier",
					},
				},
			},
		},
		{
			name: "Identifier",
			pos:  position{line: 522, col: 1, offset: 13238},
			expr: &actionExpr{
				pos: position{line: 522, col: 15, offset: 13252},
				run: (*parser).callonIdentifier1,
				expr: &seqExpr{
					pos: position{line: 522, col: 15, offset: 13252},
					exprs: []interface{}{
						&charClassMatcher{
							pos:        position{line: 522, col: 15, offset: 13252},
							val:        "[a-zA-Z_]",
							chars:      []rune{'_'},
							ranges:     []rune{'a', 'z', 'A', 'Z'},
							ignoreCase: false,
							inverted:   false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 522, col: 25, offset: 13262},
							expr: &charClassMatcher{
								pos:        position{line: 522, col: 25, offset: 13262},
								val:        "[a-zA-Z0-9_]",
								chars:      []rune{'_'},
								ranges:     []rune{'a', 'z', 'A', 'Z', '0', '9'},
								ignoreCase: false,
								inverted:   false,
							},
						},
					},
				},
			},
		},
		{
			name: "EmptyLine",
			pos:  position{line: 526, col: 1, offset: 13310},
			expr: &actionExpr{
				pos: position{line: 526, col: 14, offset: 13323},
				run: (*parser).callonEmptyLine1,
				expr: &seqExpr{
					pos: position{line: 526, col: 14, offset: 13323},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 526, col: 14, offset: 13323},
							name: "_",
						},
						&charClassMatcher{
							pos:        position{line: 526, col: 16, offset: 13325},
							val:        "[\\n]",
							chars:      []rune{'\n'},
							ignoreCase: false,
							inverted:   false,
						},
					},
				},
			},
		},
		{
			name:        "_",
			displayName: "\"whitespace\"",
			pos:         position{line: 530, col: 1, offset: 13353},
			expr: &actionExpr{
				pos: position{line: 530, col: 19, offset: 13371},
				run: (*parser).callon_1,
				expr: &zeroOrMoreExpr{
					pos: position{line: 530, col: 19, offset: 13371},
					expr: &charClassMatcher{
						pos:        position{line: 530, col: 19, offset: 13371},
						val:        "[ \\t]",
						chars:      []rune{' ', '\t'},
						ignoreCase: false,
						inverted:   false,
					},
				},
			},
		},
		{
			name:        "__",
			displayName: "\"whitespace\"",
			pos:         position{line: 531, col: 1, offset: 13398},
			expr: &actionExpr{
				pos: position{line: 531, col: 20, offset: 13417},
				run: (*parser).callon__1,
				expr: &charClassMatcher{
					pos:        position{line: 531, col: 20, offset: 13417},
					val:        "[ \\t]",
					chars:      []rune{' ', '\t'},
					ignoreCase: false,
					inverted:   false,
				},
			},
		},
		{
			name: "NL",
			pos:  position{line: 532, col: 1, offset: 13444},
			expr: &choiceExpr{
				pos: position{line: 532, col: 7, offset: 13450},
				alternatives: []interface{}{
					&charClassMatcher{
						pos:        position{line: 532, col: 7, offset: 13450},
						val:        "[\\n]",
						chars:      []rune{'\n'},
						ignoreCase: false,
						inverted:   false,
					},
					&andExpr{
						pos: position{line: 532, col: 14, offset: 13457},
						expr: &ruleRefExpr{
							pos:  position{line: 532, col: 15, offset: 13458},
							name: "EOF",
						},
					},
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 533, col: 1, offset: 13462},
			expr: &notExpr{
				pos: position{line: 533, col: 8, offset: 13469},
				expr: &anyMatcher{
					line: 533, col: 9, offset: 13470,
				},
			},
		},
	},
}

func (c *current) onInput1(l interface{}) (interface{}, error) {
	return &Root{List: l.(*List), GraphNode: NewNode(pos(c.pos))}, nil
}

func (p *parser) callonInput1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onInput1(stack["l"])
}

func (c *current) onList2(node, list interface{}) (interface{}, error) {
	listItem := list.(*List)
	if node != nil {
		listItem.Nodes = append([]Node{node.(Node)}, listItem.Nodes...)
	}
	listItem.Position = pos(c.pos)

	return listItem, nil
}

func (p *parser) callonList2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onList2(stack["node"], stack["list"])
}

func (c *current) onList8() (interface{}, error) {
	return &List{GraphNode: NewNode(pos(c.pos))}, nil
}

func (p *parser) callonList8() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onList8()
}

func (c *current) onIndentedList1(list interface{}) (interface{}, error) {
	return list, nil
}

func (p *parser) callonIndentedList1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIndentedList1(stack["list"])
}

func (c *current) onIndentedRawText1(t interface{}) (interface{}, error) {
	return &Text{Value: t.(string), GraphNode: NewNode(pos(c.pos))}, nil
}

func (p *parser) callonIndentedRawText1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIndentedRawText1(stack["t"])
}

func (c *current) onRawText2(rt, tail interface{}) (interface{}, error) {
	return rt.(string) + tail.(string), nil
}

func (p *parser) callonRawText2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onRawText2(stack["rt"], stack["tail"])
}

func (c *current) onRawText10() (interface{}, error) {
	return "", nil
}

func (p *parser) callonRawText10() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onRawText10()
}

func (c *current) onRawText16(head, tail interface{}) (interface{}, error) {
	return string(head.([]byte)) + tail.(string), nil
}

func (p *parser) callonRawText16() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onRawText16(stack["head"], stack["tail"])
}

func (c *current) onListNode14() (interface{}, error) {
	return nil, nil
}

func (p *parser) callonListNode14() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onListNode14()
}

func (c *current) onDocType1(val interface{}) (interface{}, error) {
	return &DocType{Value: val.(string), GraphNode: NewNode(pos(c.pos))}, nil
}

func (p *parser) callonDocType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDocType1(stack["val"])
}

func (c *current) onTag1(tag, list interface{}) (interface{}, error) {
	tagElem := tag.(*Tag)

	if list != nil {
		tagElem.Block = list.(*List)
	}

	return tagElem, nil
}

func (p *parser) callonTag1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTag1(stack["tag"], stack["list"])
}

func (c *current) onTagHeader2(name, attrs, tl interface{}) (interface{}, error) {
	tag := &Tag{Name: name.(string), GraphNode: NewNode(pos(c.pos))}
	if attrs != nil {
		tag.Attributes = attrs.([]*Attribute)
	}
	if tl != nil {
		tag.Text = toSlice(tl)[1].(*TextList)
	}
	return tag, nil
}

func (p *parser) callonTagHeader2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTagHeader2(stack["name"], stack["attrs"], stack["tl"])
}

func (c *current) onTagHeader17(name, attrs, text interface{}) (interface{}, error) {
	tag := &Tag{Name: name.(string), GraphNode: NewNode(pos(c.pos))}
	if attrs != nil {
		tag.Attributes = attrs.([]*Attribute)
	}
	if text != nil {
		tag.Block = text.(*Text)
	}
	return tag, nil
}

func (p *parser) callonTagHeader17() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTagHeader17(stack["name"], stack["attrs"], stack["text"])
}

func (c *current) onTagHeader30(name, attrs, unescaped, expr interface{}) (interface{}, error) {
	tag := &Tag{Name: name.(string), GraphNode: NewNode(pos(c.pos))}
	if attrs != nil {
		tag.Attributes = attrs.([]*Attribute)
	}
	if expr != nil {
		intr := &Interpolation{Expr: expr.(Expression), GraphNode: NewNode(pos(c.pos))}
		if unescaped != nil {
			intr.Unescaped = true
		}
		tag.Block = &TextList{Nodes: []Node{intr}, GraphNode: NewNode(pos(c.pos))}
	}
	return tag, nil
}

func (p *parser) callonTagHeader30() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTagHeader30(stack["name"], stack["attrs"], stack["unescaped"], stack["expr"])
}

func (c *current) onTagName2() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonTagName2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTagName2()
}

func (c *current) onTagName9() (interface{}, error) {
	return "div", nil
}

func (p *parser) callonTagName9() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTagName9()
}

func (c *current) onTagAttributes2(head, tail interface{}) (interface{}, error) {
	tailElem := []*Attribute{}

	if tail != nil {
		tailElem = tail.([]*Attribute)
	}

	return append(head.([]*Attribute), tailElem...), nil
}

func (p *parser) callonTagAttributes2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTagAttributes2(stack["head"], stack["tail"])
}

func (c *current) onTagAttributes11(head, tail interface{}) (interface{}, error) {
	tailElem := []*Attribute{}

	if tail != nil {
		tailElem = tail.([]*Attribute)
	}

	vals := toSlice(toSlice(head)[2])

	if len(vals) == 0 {
		return tailElem, nil
	}

	headAttrs := vals[0].([]*Attribute)
	restAttrs := toSlice(vals[1])

	for _, a := range restAttrs {
		restAttr := toSlice(a)
		headAttrs = append(headAttrs, restAttr[1].([]*Attribute)...)
	}

	return append(headAttrs, tailElem...), nil
}

func (p *parser) callonTagAttributes11() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTagAttributes11(stack["head"], stack["tail"])
}

func (c *current) onTagAttributeClass1(name interface{}) (interface{}, error) {
	return []*Attribute{&Attribute{Name: "class", Value: &StringExpression{Value: name.(string)}, GraphNode: NewNode(pos(c.pos))}}, nil
}

func (p *parser) callonTagAttributeClass1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTagAttributeClass1(stack["name"])
}

func (c *current) onTagAttributeID1(name interface{}) (interface{}, error) {
	return []*Attribute{&Attribute{Name: "id", Value: &StringExpression{Value: name.(string)}, GraphNode: NewNode(pos(c.pos))}}, nil
}

func (p *parser) callonTagAttributeID1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTagAttributeID1(stack["name"])
}

func (c *current) onTagAttribute2(name, value interface{}) (interface{}, error) {
	return []*Attribute{&Attribute{Name: name.(string), Value: value.(Expression), GraphNode: NewNode(pos(c.pos))}}, nil
}

func (p *parser) callonTagAttribute2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTagAttribute2(stack["name"], stack["value"])
}

func (c *current) onTagAttribute11(name interface{}) (interface{}, error) {
	return []*Attribute{&Attribute{Name: name.(string), GraphNode: NewNode(pos(c.pos))}}, nil
}

func (p *parser) callonTagAttribute11() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTagAttribute11(stack["name"])
}

func (c *current) onTagAttributeName2(tn interface{}) (interface{}, error) {
	return tn, nil
}

func (p *parser) callonTagAttributeName2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTagAttributeName2(stack["tn"])
}

func (c *current) onTagAttributeName8(tn interface{}) (interface{}, error) {
	return tn, nil
}

func (p *parser) callonTagAttributeName8() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTagAttributeName8(stack["tn"])
}

func (c *current) onTagAttributeName14(tn interface{}) (interface{}, error) {
	return tn, nil
}

func (p *parser) callonTagAttributeName14() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTagAttributeName14(stack["tn"])
}

func (c *current) onClassName1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonClassName1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onClassName1()
}

func (c *current) onIdName1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonIdName1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIdName1()
}

func (c *current) onTagAttributeNameLiteral1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonTagAttributeNameLiteral1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTagAttributeNameLiteral1()
}

func (c *current) onIf1(expr, block, elseNode interface{}) (interface{}, error) {
	ifElem := &If{Condition: expr.(Expression), PositiveBlock: block.(Node), GraphNode: NewNode(pos(c.pos))}
	if elseNode != nil {
		ifElem.NegativeBlock = elseNode.(Node)
	}
	return ifElem, nil
}

func (p *parser) callonIf1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIf1(stack["expr"], stack["block"], stack["elseNode"])
}

func (c *current) onElse2(node interface{}) (interface{}, error) {
	return node, nil
}

func (p *parser) callonElse2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onElse2(stack["node"])
}

func (c *current) onElse8(block interface{}) (interface{}, error) {
	return block, nil
}

func (p *parser) callonElse8() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onElse8(stack["block"])
}

func (c *current) onEach1(v1, v2, expr, block interface{}) (interface{}, error) {
	eachElem := &Each{GraphNode: NewNode(pos(c.pos)), Block: block.(Node), ElementVariable: v1.(*Variable), Container: expr.(Expression)}
	v2Slice := toSlice(v2)

	if len(v2Slice) != 0 {
		eachElem.IndexVariable = v2Slice[3].(*Variable)
	}

	return eachElem, nil
}

func (p *parser) callonEach1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onEach1(stack["v1"], stack["v2"], stack["expr"], stack["block"])
}

func (c *current) onAssignment1(vr, expr interface{}) (interface{}, error) {
	return &Assignment{Variable: vr.(*Variable), Expression: expr.(Expression), GraphNode: NewNode(pos(c.pos))}, nil
}

func (p *parser) callonAssignment1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onAssignment1(stack["vr"], stack["expr"])
}

func (c *current) onMixin1(name, args, list interface{}) (interface{}, error) {
	mixinElem := &Mixin{Name: name.(string), Block: list.(Node), GraphNode: NewNode(pos(c.pos))}
	if args != nil {
		mixinElem.Arguments = args.([]MixinArgument)
	}
	return mixinElem, nil
}

func (p *parser) callonMixin1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMixin1(stack["name"], stack["args"], stack["list"])
}

func (c *current) onMixinArguments2(head, tail interface{}) (interface{}, error) {
	args := []MixinArgument{head.(MixinArgument)}

	if tail != nil {
		tailSlice := toSlice(tail)

		for _, arg := range tailSlice {
			argSlice := toSlice(arg)
			args = append(args, argSlice[3].(MixinArgument))
		}
	}

	return args, nil
}

func (p *parser) callonMixinArguments2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMixinArguments2(stack["head"], stack["tail"])
}

func (c *current) onMixinArguments15() (interface{}, error) {
	return nil, nil
}

func (p *parser) callonMixinArguments15() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMixinArguments15()
}

func (c *current) onMixinArgument1(name, def interface{}) (interface{}, error) {
	argElem := MixinArgument{Name: name.(*Variable), GraphNode: NewNode(pos(c.pos))}

	if def != nil {
		defSlice := toSlice(def)
		argElem.Default = defSlice[3].(Expression)
	}

	return argElem, nil
}

func (p *parser) callonMixinArgument1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMixinArgument1(stack["name"], stack["def"])
}

func (c *current) onMixinCall1(name, args interface{}) (interface{}, error) {
	mcElem := &MixinCall{Name: name.(string), GraphNode: NewNode(pos(c.pos))}
	if args != nil {
		mcElem.Arguments = args.([]Expression)
	}
	return mcElem, nil
}

func (p *parser) callonMixinCall1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMixinCall1(stack["name"], stack["args"])
}

func (c *current) onCallArguments1(head, tail interface{}) (interface{}, error) {
	args := []Expression{}

	if head != nil {
		args = append(args, head.(Expression))
	}

	if tail != nil {
		tailSlice := toSlice(tail)

		for _, arg := range tailSlice {
			argSlice := toSlice(arg)
			args = append(args, argSlice[3].(Expression))
		}
	}

	return args, nil
}

func (p *parser) callonCallArguments1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCallArguments1(stack["head"], stack["tail"])
}

func (c *current) onImport1(file interface{}) (interface{}, error) {
	return &Import{File: file.(string), GraphNode: NewNode(pos(c.pos))}, nil
}

func (p *parser) callonImport1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onImport1(stack["file"])
}

func (c *current) onExtend1(file interface{}) (interface{}, error) {
	return &Extend{File: file.(string), GraphNode: NewNode(pos(c.pos))}, nil
}

func (p *parser) callonExtend1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExtend1(stack["file"])
}

func (c *current) onBlock1(mod, name, list interface{}) (interface{}, error) {
	block := &Block{Name: name.(string), GraphNode: NewNode(pos(c.pos))}

	if mod != nil {
		modSlice := toSlice(mod)

		if string(modSlice[1].([]byte)) == "append" {
			block.Modifier = "append"
		} else {
			block.Modifier = "prepend"
		}
	}

	if list != nil {
		block.Block = list.(*List)
	}

	return block, nil
}

func (p *parser) callonBlock1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBlock1(stack["mod"], stack["name"], stack["list"])
}

func (c *current) onComment1(silent, comment interface{}) (interface{}, error) {
	isSilent := silent != nil
	return &Comment{Value: comment.(string), Silent: isSilent, GraphNode: NewNode(pos(c.pos))}, nil
}

func (p *parser) callonComment1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onComment1(stack["silent"], stack["comment"])
}

func (c *current) onLineText1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonLineText1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLineText1()
}

func (c *current) onPipeText1(tl interface{}) (interface{}, error) {
	return tl, nil
}

func (p *parser) callonPipeText1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPipeText1(stack["tl"])
}

func (c *current) onTextList2(intr, tl interface{}) (interface{}, error) {
	intNode := intr.(*Interpolation)

	if tl != nil {
		tlnode := tl.(*TextList)
		return &TextList{
			Nodes:     append([]Node{intNode}, tlnode.Nodes...),
			GraphNode: NewNode(pos(c.pos)),
		}, nil
	}

	return TextList{
		Nodes:     []Node{intNode},
		GraphNode: NewNode(pos(c.pos)),
	}, nil
}

func (p *parser) callonTextList2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTextList2(stack["intr"], stack["tl"])
}

func (c *current) onTextList8() (interface{}, error) {
	return &TextList{GraphNode: NewNode(pos(c.pos))}, nil
}

func (p *parser) callonTextList8() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTextList8()
}

func (c *current) onTextList11(ch, tl interface{}) (interface{}, error) {
	tlnode := tl.(*TextList)
	chstr := string(ch.([]byte))

	if len(tlnode.Nodes) > 0 {
		if tn, ok := tlnode.Nodes[0].(*Text); ok {
			tlnode.Nodes[0] = &Text{Value: chstr + tn.Value}
			return tlnode, nil
		}
	}

	tlnode.Nodes = append([]Node{&Text{Value: chstr, GraphNode: NewNode(pos(c.pos))}}, tlnode.Nodes...)
	tlnode.Position = pos(c.pos)

	return tlnode, nil
}

func (p *parser) callonTextList11() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTextList11(stack["ch"], stack["tl"])
}

func (c *current) onInterpolation1(mod, expr interface{}) (interface{}, error) {
	intElem := &Interpolation{Expr: expr.(Expression), GraphNode: NewNode(pos(c.pos))}
	if string(mod.([]byte)) == "!" {
		intElem.Unescaped = true
	}
	return intElem, nil
}

func (p *parser) callonInterpolation1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onInterpolation1(stack["mod"], stack["expr"])
}

func (c *current) onExpressionTernery1(cnd, rest interface{}) (interface{}, error) {
	if rest == nil {
		return cnd, nil
	}

	restSlice := toSlice(rest)

	return &BinaryExpression{
		X: &BinaryExpression{
			X:         cnd.(Expression),
			Y:         restSlice[3].(Expression),
			Op:        "&&",
			GraphNode: NewNode(pos(c.pos)),
		},
		Y:         restSlice[7].(Expression),
		Op:        "||",
		GraphNode: NewNode(pos(c.pos)),
	}, nil
}

func (p *parser) callonExpressionTernery1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpressionTernery1(stack["cnd"], stack["rest"])
}

func (c *current) onExpressionBinOp1(first, rest interface{}) (interface{}, error) {
	return binary(first, rest, c.pos)
}

func (p *parser) callonExpressionBinOp1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpressionBinOp1(stack["first"], stack["rest"])
}

func (c *current) onExpressionCmpOp1(first, rest interface{}) (interface{}, error) {
	return binary(first, rest, c.pos)
}

func (p *parser) callonExpressionCmpOp1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpressionCmpOp1(stack["first"], stack["rest"])
}

func (c *current) onExpressionAddOp1(first, rest interface{}) (interface{}, error) {
	return binary(first, rest, c.pos)
}

func (p *parser) callonExpressionAddOp1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpressionAddOp1(stack["first"], stack["rest"])
}

func (c *current) onExpressionMulOp1(first, rest interface{}) (interface{}, error) {
	return binary(first, rest, c.pos)
}

func (p *parser) callonExpressionMulOp1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpressionMulOp1(stack["first"], stack["rest"])
}

func (c *current) onExpressionUnaryOp2(op, ex interface{}) (interface{}, error) {
	return &UnaryExpression{X: ex.(Expression), Op: op.(string), GraphNode: NewNode(pos(c.pos))}, nil
}

func (p *parser) callonExpressionUnaryOp2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpressionUnaryOp2(stack["op"], stack["ex"])
}

func (c *current) onExpressionFactor2(e interface{}) (interface{}, error) {
	return e, nil
}

func (p *parser) callonExpressionFactor2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpressionFactor2(stack["e"])
}

func (c *current) onStringExpression1(s interface{}) (interface{}, error) {
	return &StringExpression{Value: s.(string), GraphNode: NewNode(pos(c.pos))}, nil
}

func (p *parser) callonStringExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onStringExpression1(stack["s"])
}

func (c *current) onNumberExpression1(dec, ex interface{}) (interface{}, error) {
	if dec != nil || ex != nil {
		val, err := strconv.ParseFloat(string(c.text), 64)
		return &FloatExpression{Value: val, GraphNode: NewNode(pos(c.pos))}, err
	} else {
		val, err := strconv.ParseInt(string(c.text), 10, 64)
		return &IntegerExpression{Value: val, GraphNode: NewNode(pos(c.pos))}, err
	}
}

func (p *parser) callonNumberExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNumberExpression1(stack["dec"], stack["ex"])
}

func (c *current) onNilExpression1() (interface{}, error) {
	return &NilExpression{GraphNode: NewNode(pos(c.pos))}, nil
}

func (p *parser) callonNilExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNilExpression1()
}

func (c *current) onBooleanExpression1(b interface{}) (interface{}, error) {
	return &BooleanExpression{Value: b.(bool), GraphNode: NewNode(pos(c.pos))}, nil
}

func (p *parser) callonBooleanExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBooleanExpression1(stack["b"])
}

func (c *current) onMemberExpression1(field, member interface{}) (interface{}, error) {
	memberSlice := toSlice(member)
	cur := field.(Expression)

	for _, m := range memberSlice {
		mSlice := toSlice(m)

		if mSlice[0] != nil && string(mSlice[0].([]byte)) == "." {
			cur = &MemberExpression{X: cur, Name: mSlice[1].(string), GraphNode: NewNode(pos(c.pos))}
		} else {
			cur = &FunctionCallExpression{X: cur, Arguments: mSlice[1].([]Expression), GraphNode: NewNode(pos(c.pos))}
		}
	}

	return cur, nil
}

func (p *parser) callonMemberExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMemberExpression1(stack["field"], stack["member"])
}

func (c *current) onField1(variable interface{}) (interface{}, error) {
	return &FieldExpression{Variable: variable.(*Variable), GraphNode: NewNode(pos(c.pos))}, nil
}

func (p *parser) callonField1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onField1(stack["variable"])
}

func (c *current) onUnaryOp1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonUnaryOp1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onUnaryOp1()
}

func (c *current) onAddOp1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonAddOp1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onAddOp1()
}

func (c *current) onMulOp1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonMulOp1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMulOp1()
}

func (c *current) onCmpOp1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonCmpOp1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCmpOp1()
}

func (c *current) onBinOp1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonBinOp1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBinOp1()
}

func (c *current) onString1() (interface{}, error) {
	return strconv.Unquote(string(c.text))
}

func (p *parser) callonString1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onString1()
}

func (c *current) onBool2() (interface{}, error) {
	return true, nil
}

func (p *parser) callonBool2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBool2()
}

func (c *current) onBool4() (interface{}, error) {
	return false, nil
}

func (p *parser) callonBool4() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBool4()
}

func (c *current) onNull1() (interface{}, error) {
	return nil, nil
}

func (p *parser) callonNull1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNull1()
}

func (c *current) onVariable1(ident interface{}) (interface{}, error) {
	return &Variable{Name: ident.(string), GraphNode: NewNode(pos(c.pos))}, nil
}

func (p *parser) callonVariable1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onVariable1(stack["ident"])
}

func (c *current) onIdentifier1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonIdentifier1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIdentifier1()
}

func (c *current) onEmptyLine1() (interface{}, error) {
	return nil, nil
}

func (p *parser) callonEmptyLine1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onEmptyLine1()
}

func (c *current) on_1() (interface{}, error) {
	return nil, nil
}

func (p *parser) callon_1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.on_1()
}

func (c *current) on__1() (interface{}, error) {
	return nil, nil
}

func (p *parser) callon__1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.on__1()
}

var (
	// errNoRule is returned when the grammar to parse has no rule.
	errNoRule = errors.New("grammar has no rule")

	// errInvalidEncoding is returned when the source is not properly
	// utf8-encoded.
	errInvalidEncoding = errors.New("invalid encoding")

	// errMaxExprCnt is used to signal that the maximum number of
	// expressions have been parsed.
	errMaxExprCnt = errors.New("max number of expresssions parsed")
)

// Option is a function that can set an option on the parser. It returns
// the previous setting as an Option.
type Option func(*parser) Option

// MaxExpressions creates an Option to stop parsing after the provided
// number of expressions have been parsed, if the value is 0 then the parser will
// parse for as many steps as needed (possibly an infinite number).
//
// The default for maxExprCnt is 0.
func MaxExpressions(maxExprCnt uint64) Option {
	return func(p *parser) Option {
		oldMaxExprCnt := p.maxExprCnt
		p.maxExprCnt = maxExprCnt
		return MaxExpressions(oldMaxExprCnt)
	}
}

// Statistics adds a user provided Stats struct to the parser to allow
// the user to process the results after the parsing has finished.
// Also the key for the "no match" counter is set.
//
// Example usage:
//
//     input := "input"
//     stats := Stats{}
//     _, err := Parse("input-file", []byte(input), Statistics(&stats, "no match"))
//     if err != nil {
//         log.Panicln(err)
//     }
//     b, err := json.MarshalIndent(stats.ChoiceAltCnt, "", "  ")
//     if err != nil {
//         log.Panicln(err)
//     }
//     fmt.Println(string(b))
//
func Statistics(stats *Stats, choiceNoMatch string) Option {
	return func(p *parser) Option {
		oldStats := p.Stats
		p.Stats = stats
		oldChoiceNoMatch := p.choiceNoMatch
		p.choiceNoMatch = choiceNoMatch
		if p.Stats.ChoiceAltCnt == nil {
			p.Stats.ChoiceAltCnt = make(map[string]map[string]int)
		}
		return Statistics(oldStats, oldChoiceNoMatch)
	}
}

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
		if closeErr := f.Close(); closeErr != nil {
			err = closeErr
		}
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

type recoveryExpr struct {
	pos          position
	expr         interface{}
	recoverExpr  interface{}
	failureLabel []string
}

type seqExpr struct {
	pos   position
	exprs []interface{}
}

type throwExpr struct {
	pos   position
	label string
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
	pos             position
	val             string
	basicLatinChars [128]bool
	chars           []rune
	ranges          []rune
	classes         []*unicode.RangeTable
	ignoreCase      bool
	inverted        bool
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
	stats := Stats{
		ChoiceAltCnt: make(map[string]map[string]int),
	}

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
		maxFailExpected: make([]string, 0, 20),
		Stats:           &stats,
	}
	p.setOptions(opts)

	if p.maxExprCnt == 0 {
		p.maxExprCnt = math.MaxUint64
	}

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

const choiceNoMatch = -1

// Stats stores some statistics, gathered during parsing
type Stats struct {
	// ExprCnt counts the number of expressions processed during parsing
	// This value is compared to the maximum number of expressions allowed
	// (set by the MaxExpressions option).
	ExprCnt uint64

	// ChoiceAltCnt is used to count for each ordered choice expression,
	// which alternative is used how may times.
	// These numbers allow to optimize the order of the ordered choice expression
	// to increase the performance of the parser
	//
	// The outer key of ChoiceAltCnt is composed of the name of the rule as well
	// as the line and the column of the ordered choice.
	// The inner key of ChoiceAltCnt is the number (one-based) of the matching alternative.
	// For each alternative the number of matches are counted. If an ordered choice does not
	// match, a special counter is incremented. The name of this counter is set with
	// the parser option Statistics.
	// For an alternative to be included in ChoiceAltCnt, it has to match at least once.
	ChoiceAltCnt map[string]map[string]int
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

	// parse fail
	maxFailPos            position
	maxFailExpected       []string
	maxFailInvertExpected bool

	// max number of expressions to be parsed
	maxExprCnt uint64

	*Stats

	choiceNoMatch string
	// recovery expression stack, keeps track of the currently available recovery expression, these are traversed in reverse
	recoveryStack []map[string]interface{}
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

// push a recovery expression with its labels to the recoveryStack
func (p *parser) pushRecovery(labels []string, expr interface{}) {
	if cap(p.recoveryStack) == len(p.recoveryStack) {
		// create new empty slot in the stack
		p.recoveryStack = append(p.recoveryStack, nil)
	} else {
		// slice to 1 more
		p.recoveryStack = p.recoveryStack[:len(p.recoveryStack)+1]
	}

	m := make(map[string]interface{}, len(labels))
	for _, fl := range labels {
		m[fl] = expr
	}
	p.recoveryStack[len(p.recoveryStack)-1] = m
}

// pop a recovery expression from the recoveryStack
func (p *parser) popRecovery() {
	// GC that map
	p.recoveryStack[len(p.recoveryStack)-1] = nil

	p.recoveryStack = p.recoveryStack[:len(p.recoveryStack)-1]
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
			p.maxFailExpected = p.maxFailExpected[:0]
		}

		if p.maxFailInvertExpected {
			want = "!" + want
		}
		p.maxFailExpected = append(p.maxFailExpected, want)
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
			maxFailExpectedMap := make(map[string]struct{}, len(p.maxFailExpected))
			for _, v := range p.maxFailExpected {
				maxFailExpectedMap[v] = struct{}{}
			}
			expected := make([]string, 0, len(maxFailExpectedMap))
			eof := false
			if _, ok := maxFailExpectedMap["!."]; ok {
				delete(maxFailExpectedMap, "!.")
				eof = true
			}
			for k := range maxFailExpectedMap {
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

	p.ExprCnt++
	if p.ExprCnt > p.maxExprCnt {
		panic(errMaxExprCnt)
	}

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
	case *recoveryExpr:
		val, ok = p.parseRecoveryExpr(expr)
	case *ruleRefExpr:
		val, ok = p.parseRuleRefExpr(expr)
	case *seqExpr:
		val, ok = p.parseSeqExpr(expr)
	case *throwExpr:
		val, ok = p.parseThrowExpr(expr)
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

func (p *parser) incChoiceAltCnt(ch *choiceExpr, altI int) {
	choiceIdent := fmt.Sprintf("%s %d:%d", p.rstack[len(p.rstack)-1].name, ch.pos.line, ch.pos.col)
	m := p.ChoiceAltCnt[choiceIdent]
	if m == nil {
		m = make(map[string]int)
		p.ChoiceAltCnt[choiceIdent] = m
	}
	// We increment altI by 1, so the keys do not start at 0
	alt := strconv.Itoa(altI + 1)
	if altI == choiceNoMatch {
		alt = p.choiceNoMatch
	}
	m[alt]++
}

func (p *parser) parseChoiceExpr(ch *choiceExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseChoiceExpr"))
	}

	for altI, alt := range ch.alternatives {
		// dummy assignment to prevent compile error if optimized
		_ = altI

		p.pushV()
		val, ok := p.parseExpr(alt)
		p.popV()
		if ok {
			p.incChoiceAltCnt(ch, altI)
			return val, ok
		}
	}
	p.incChoiceAltCnt(ch, choiceNoMatch)
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

func (p *parser) parseRecoveryExpr(recover *recoveryExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRecoveryExpr (" + strings.Join(recover.failureLabel, ",") + ")"))
	}

	p.pushRecovery(recover.failureLabel, recover.recoverExpr)
	val, ok := p.parseExpr(recover.expr)
	p.popRecovery()

	return val, ok
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

	vals := make([]interface{}, 0, len(seq.exprs))

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

func (p *parser) parseThrowExpr(expr *throwExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseThrowExpr"))
	}

	for i := len(p.recoveryStack) - 1; i >= 0; i-- {
		if recoverExpr, ok := p.recoveryStack[i][expr.label]; ok {
			if val, ok := p.parseExpr(recoverExpr); ok {
				return val, ok
			}
		}
	}

	return nil, false
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
