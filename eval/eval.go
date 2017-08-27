package eval

import (
    "fmt"
    "strings"
    "github.com/crides/gysp/parse"
)

// Tree of Node --- Eval --> Objects
type ObjectType int

const (
    OBJECT_NIL  ObjectType = iota   // val: nil
    OBJECT_BOOL     // val: bool
    OBJECT_INT      // val: int
    OBJECT_FLOAT    // val: float64
    OBJECT_CMPLX    // val: complex128

    OBJECT_STR      // val: string
    OBJECT_LIST     // val: []Object
    OBJECT_DICT     // val: map[Object]*Object

    OBJECT_PRIM     // val: func(...Object) Object
    OBJECT_MACRO    // val: func(Node) Node; actually a primitive
    OBJECT_FUNC     // val: Func

    OBJECT_CLASS    // Scheme for a class declaration;
                    // val: map[string]interface{} -> map[var]initializers
    OBJECT_OBJ      // A instance of class
                    // val: map[string]interface{} -> map[var]values
)

func (ot ObjectType) String() string {
    switch ot {
    case OBJECT_NIL:
        return "nil"
    case OBJECT_BOOL:
        return "bool"
    case OBJECT_INT:
        return "int"
    case OBJECT_FLOAT:
        return "float"
    case OBJECT_CMPLX:
        return "complex"
    case OBJECT_STR:
        return "string"
    case OBJECT_LIST:
        return "list"
    case OBJECT_DICT:
        return "dict"
    case OBJECT_PRIM:
        return "built-in"
    case OBJECT_MACRO:
        return "macro"
    case OBJECT_FUNC:
        return "function"
    case OBJECT_CLASS:
        return fmt.Sprintf("<class %s>", ot)
    case OBJECT_OBJ:
        return fmt.Sprintf("<%s object>", ot)
    }
    panic(fmt.Sprintf("Unknown type %d!", ot))
}

type Object struct {
    typ     ObjectType
    val     interface{}
}

func NewObject(t ObjectType, v interface{}) *Object {
    return &Object{t, v}
}

func (o * Object) Typ() ObjectType {
    return o.typ
}

func (o * Object) Val() interface{} {
    return o.val
}

func (o * Object) String() string {
    switch o.typ {
    case OBJECT_NIL:
        return "nil"
    case OBJECT_BOOL, OBJECT_INT, OBJECT_FLOAT, OBJECT_CMPLX:
        return fmt.Sprint(o.val)

    case OBJECT_STR:
        return fmt.Sprintf("%q", o.val)
    case OBJECT_LIST:
        strs := make([]string, 0)
        for _, item := range o.val.(ListType) {
            strs = append(strs, item.String())
        }
        return "[" + strings.Join(strs, " ") + "]"
    case OBJECT_DICT:
        strs := make([]string, 0)
        for key, val := range o.val.(DictType) {
            strs = append(strs, key.String() + ": " + val.String())
        }
        return "{" + strings.Join(strs, ", ") + "}"
    }
    return "unknown"
}

var (
    GYSP_NIL = NewObject(OBJECT_NIL, nil)
    GYSP_TRUE = NewObject(OBJECT_BOOL, true)
    GYSP_FALSE = NewObject(OBJECT_BOOL, false)
)

type Func struct {      // Represents a gysp function
    vars    []string
    //rest    bool        // Whether the last argument is variadic
    env     *Env        // The outer environment; for implementing closures
    body    []parse.CallNode
}

type PrimFunc func(...*Object) *Object

type MacroFunc func(...parse.Node) parse.Node

type DictType map[Object]*Object

type ListType []*Object

type Generator struct {
    sources     [][]interface{}
    ptrs        []int
    sizes       []int
    size        int
}

func NewGenerator() *Generator {
    return new(Generator)
}

func (g * Generator) AddSource(a []interface{}) {
    g.sources = append(g.sources, a)
    g.ptrs = append(g.ptrs, 0)
    g.sizes = append(g.sizes, len(a))
    g.size ++
}

func (g * Generator) Generate(f func(...interface{})) {
    for g.ptrs[0] < g.sizes[0] {
        vals := make([]interface{}, g.size)
        for i := 0; i < g.size; i ++ {
            vals[i] = g.sources[i][g.ptrs[i]]
        }
        f(vals...)

        // Increment pointers
        g.ptrs[g.size - 1] ++       // Last pointer
        for i := g.size - 1; i > 0; i -- {
            if g.ptrs[i] >= g.sizes[i] {
                g.ptrs[i] = 0
                g.ptrs[i - 1] ++
            }
        }
    }
}

func Eval(node parse.Node, env *Env) *Object {
    switch n := node.(type) {
    // Literals
    case *parse.LiteralNode:
        obj_flag := OBJECT_NIL
        switch n.Val.(type) {
        case int:
            obj_flag = OBJECT_INT
        case float64:
            obj_flag = OBJECT_FLOAT
        case complex128:
            obj_flag = OBJECT_CMPLX
        case string:
            obj_flag = OBJECT_STR
        }
        return NewObject(obj_flag, n.Val)
    case *parse.ListNode:
        list := make(ListType, len(n.List))
        for i, node := range n.List {
            list[i] = Eval(node, env)
        }
        return NewObject(OBJECT_LIST, list)
    case *parse.DictNode:
        dict := make(DictType)
        for k, v := range n.Dict {
            dict[*Eval(k, env)] = Eval(v, env)
        }
        return NewObject(OBJECT_DICT, dict)

    // Variable and references
    case *parse.SymNode:
        if !strings.Contains(n.Name, "/") && !strings.Contains(n.Name, ".") {
            // Just a variable; no subs
            return env.GetVar(n.Name)
        }
        panic("Not implemented!")

    // Function calls
    case *parse.CallNode:
        fun := Eval(n.Fun, env)
        switch fun.Typ() {
        case OBJECT_PRIM:
            arglen := len(n.Arglist)
            arglist := make(ListType, arglen)
            for i := 0; i < arglen; i ++ {
                arglist[i] = Eval(n.Arglist[i], env)
            }
            return fun.val.(func(...*Object) *Object)(arglist...)
        case OBJECT_MACRO:
            return Eval(fun.val.(MacroFunc)(n.Arglist...), env)
        }
        panic(fmt.Sprintf("%s object can't be used as a function!", fun.Typ().String()))
    }
    panic("Not implemented!")
}
