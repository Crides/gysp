package parse

import (
    "fmt"
    "strconv"
    re "regexp"

    "github.com/crides/gysp/color"
)

// []*Token --- Parse --> Tree of Node

type NodeType int

const (
    NODE_CALL   NodeType = iota
    NODE_LIST
    NODE_DICT
    NODE_SYM
    NODE_LIT
    NODE_WRAP       // Just wrap a object up
)

type Node interface {
    NodeTyp()   NodeType
    String()    string
}

type SymNode struct {       // Normal symbols
    Name    string
}

func NewSymNode(name string) *SymNode {
    return &SymNode{name}
}

func (sn * SymNode) NodeTyp() NodeType {
    return NODE_SYM
}

func (sn * SymNode) String() string {
    return color.Teal(sn.Name)
}

var NIL_NODE = NewSymNode("nil")

type CallNode struct {      // Function and macro calls
    Fun     Node       // Function head
    Arglist []Node     // Function argument list
}

func NewCallNode(fun Node) *CallNode {
    return &CallNode{fun, make([]Node, 0)}
}

func (cn * CallNode) NodeTyp() NodeType {
    return NODE_CALL
}

func (cn * CallNode) String() string {
    liststr := ""
    for i, arg := range cn.Arglist {
        if i == 0 {
            liststr += arg.String()
        } else {
            liststr += " " + arg.String()
        }
    }
    return color.Red(fmt.Sprintf("(%s %s)", cn.Fun, liststr))
}

func NewCallNodeFromList(node *ListNode) *CallNode {
    list := node.GetList()
    length := len(list)
    if length < 1 {
        panic("CallNode's number of items must be greater than 1!")
    }
    return &CallNode{list[0], list[1:]}
}

func (cn * CallNode) AddArg(node Node) {
    cn.Arglist = append(cn.Arglist, node)
}

func (cn * CallNode) GetArgs() []Node {
    return cn.Arglist
}

type ListNode struct {
    List    []Node     // The contents
}

func NewListNode() *ListNode {
    return &ListNode{make([]Node, 0)}
}

func (ln * ListNode) NodeTyp() NodeType {
    return NODE_LIST
}

func (ln * ListNode) String() string {
    return color.Green(fmt.Sprintf("%v", ln.List))
}

func (ln * ListNode) Add(node Node) {
    ln.List = append(ln.List, node)
}

func (ln * ListNode) GetList() []Node {
    return ln.List
}

type DictNode struct {
    Dict    map[Node]Node
}

func NewDictNode() *DictNode {
    return &DictNode{make(map[Node]Node)}
}

func NewDictNodeFromList(node *ListNode) *DictNode {
    list := node.GetList()
    length := len(list)
    if length % 2 != 0 {
        panic("DictNode's number of items must be multiples of 2!")
    }
    dn := NewDictNode()
    for i := 0; i < length; i += 2 {
        dn.Set(list[i], list[i + 1])
    }
    return dn
}

func (dn * DictNode) NodeTyp() NodeType {
    return NODE_DICT
}

func (dn * DictNode) String() string {
    return color.Yellow(fmt.Sprint(dn.Dict))
}

func (dn * DictNode) Set(key, val Node) {
    dn.Dict[key] = val
}

func (dn * DictNode) GetDict() map[Node]Node {
    return dn.Dict
}

type LiteralNode struct {       // A node that represents a literal other than lists and dicts
    Val     interface{}
}

func NewLiteralNode(val interface{}) *LiteralNode {
    return &LiteralNode{val}
}

func (ln * LiteralNode) NodeTyp() NodeType {
    return NODE_LIT
}

func (ln * LiteralNode) String() string {
    switch u := ln.Val.(type) {
    case int:
        return color.Blue(strconv.Itoa(u))
    case float64:
        return color.Mangenta(strconv.FormatFloat(u, 'f', -1, 64))
    case complex128:
        return color.Turqoise(fmt.Sprint(u))
    case string:
        return color.Yellow(u)
    }
    return ""
}

func (ln * LiteralNode) GetVal() interface{} {
    return ln.Val
}

func Parse(tokens []*Token) Node {
    node, _ := parse(tokens, TOKEN_NONE)
    return node
}

func parse(tokens []*Token, until TokenType) (Node, int) {
    root := NewListNode()
    // Return the next item & tokens read
    for i := 0; i < len(tokens); i ++ {
        token := tokens[i]
        switch t := token.Typ(); t {
        case FUNC_BEGIN, LIST_BEGIN, DICT_BEGIN:
            next, advance := parse(tokens[i + 1:], t + 1)  // Skip the left brac in the recur
            root.Add(next)
            i += advance + 1
        case FUNC_END, LIST_END, DICT_END:
            if t == until {
                switch t {
                case FUNC_END:
                    return NewCallNodeFromList(root), i
                case DICT_END:
                    return NewDictNodeFromList(root), i
                case LIST_END:
                    return root, i
                }
                panic("???")
            }
            panic("Unexpected bracket end!")
        case TOKEN:
            next := NewSymNode(token.Cont())
            //if until == TOKEN_NONE {
            //    return next, 1
            //}
            root.Add(next)
        case STRING:
            root.Add(NewLiteralNode(token.Cont()))
        case INTEGER:
            i, _ := strconv.Atoi(token.Cont())
            root.Add(NewLiteralNode(i))
        case FLOAT:
            f, _ := strconv.ParseFloat(token.Cont(), 64)
            root.Add(NewLiteralNode(f))
        case COMPLEX:
            subs := re.MustCompile(`([+-]?(?:\d*\.)?\d+)([+-]?(?:\d*\.)?\d+)j`).FindStringSubmatch(token.Cont())
            r, _ := strconv.ParseFloat(subs[1], 64)
            i, _ := strconv.ParseFloat(subs[2], 64)
            root.Add(NewLiteralNode(complex(r, i)))
        case QUOTE, QQUOTE, UNQUOTE, UNQUOTESP:
            symbol := ""
            switch t {
            case QUOTE:
                symbol = "quote"
            case QQUOTE:
                symbol = "quasiquote"
            case UNQUOTE:
                symbol = "unquote"
            case UNQUOTESP:
                symbol = "unquote-splice"
            }
            cn := NewCallNode(NewSymNode(symbol))
            next := tokens[i + 1]
            var (
                sub Node
                advance int
            )
            switch next.Typ() {
            case FUNC_BEGIN, LIST_BEGIN, DICT_BEGIN:
                sub, advance = parse(tokens[i + 2:], next.Typ() + 1)
                advance += 2        // Skipped the quote and left brac
            default:
                sub, advance = NewSymNode(tokens[i + 1].Cont()), 1
            }
            if next == nil {
                panic("Expected item after " + symbol)
            }
            cn.AddArg(sub)
            root.Add(cn)
            i += advance
        default:
            panic(fmt.Sprintf("Unknown token flag %d!", token.Typ()))
        }
    }
    if until == TOKEN_NONE {
        return root, len(tokens)
    }
    panic("Premature end of input: Expect closed parenthese!")
}
