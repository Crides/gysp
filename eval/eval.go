package eval

import (
    "fmt"
    "strings"
    "github.com/crides/gysp/parse"
)

// Moved here because of ``cycle import''
type WrapNode struct {      // Just wrap a object up
    Val     *Object
}

func WrapObject(o *Object) *WrapNode {
    return &WrapNode{o}
}

func (wn * WrapNode) NodeTyp() parse.NodeType {
    return parse.NODE_WRAP
}

func (wn * WrapNode) String() string {
    return wn.Val.String()
}

// Tree of Node --- eval --> Objects
type ObjectType int

const (
    // Primitive types
    OBJECT_NIL  ObjectType = iota   // val: nil
    OBJECT_BOOL     // val: bool
    OBJECT_INT      // val: int
    OBJECT_FLOAT    // val: float64
    OBJECT_CMPLX    // val: complex128

    // Collection types
    OBJECT_STR      // val: string
    OBJECT_LIST     // val: []Object
    OBJECT_DICT     // val: map[Object]*Object

    // Functions
    OBJECT_PRIM     // val: func(...Object) Object
    OBJECT_MACRO    // val: func(Node) Node; actually a primitive
    OBJECT_FUNC     // val: Func

    // Classes and objects
    OBJECT_CLASS    // Scheme for a class declaration;
                    // val: map[string]*Object -> map[var]initializers
    OBJECT_OBJ      // A instance of class
                    // val: map[string]*Object -> map[var]values
)

func (ot ObjectType) String() string {
    switch ot {
    case OBJECT_NIL:
        return "<nil>"
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
    case OBJECT_STR:
        return o.val.(string)
    default:
        return o.GoString()
    }
}

func (o * Object) GoString() string {
    switch o.typ {
    case OBJECT_NIL:
        return "nil"
    case OBJECT_BOOL, OBJECT_INT, OBJECT_FLOAT, OBJECT_CMPLX:
        return fmt.Sprint(o.val)

    case OBJECT_STR:
        return fmt.Sprintf("%q", o.val)
    case OBJECT_LIST:
        strs := make([]string, 0)
        for _, item := range o.val.([]*Object) {
            strs = append(strs, item.String())
        }
        return "[" + strings.Join(strs, " ") + "]"
    case OBJECT_DICT:
        strs := make([]string, 0)
        for key, val := range o.val.(map[Object]*Object) {
            strs = append(strs, key.String() + ": " + val.String())
        }
        return "{" + strings.Join(strs, ", ") + "}"
    }
    return "unknown"
}

// Some basic objects
var (
    GYSP_NIL = NewObject(OBJECT_NIL, nil)
    GYSP_TRUE = NewObject(OBJECT_BOOL, true)
    GYSP_FALSE = NewObject(OBJECT_BOOL, false)
)

// Some basic object (a class instance) operations
func (o * Object) Hasattr(attr string) bool {
    if o.typ != OBJECT_OBJ {
        panic("Hasattr() can only be used on objects!")
    }

    _, ok := o.val.(map[string]*Object)[attr]
    return ok
}

func (o * Object) Getattr(attr string) *Object {
    if o.typ != OBJECT_OBJ {
        panic("Getattr() can only be used on objects!")
    }

    if item, ok := o.val.(map[string]*Object)[attr]; ok {
        return item
    }
    panic(fmt.Sprintf("Object doesn't have attribute '%s'!", attr))
}

func (o * Object) Setattr(attr string, val *Object) {
    if o.typ != OBJECT_OBJ {
        panic("Setattr() can only be used on objects!")
    }

    if _, ok := o.val.(map[string]*Object)[attr]; ok {
        o.val.(map[string]*Object)[attr] = val
    }
    panic(fmt.Sprintf("Object doesn't have attribute '%s'!", attr))
}

// Gysp function type
type Func struct {
    vars    []string
    //rest    bool        // Whether the last argument is variadic
    env     *Env        // The outer environment; for implementing closures
    body    []*parse.CallNode
}

func EvalList(nodes []parse.Node, env *Env) []*Object {
    nodelen := len(nodes)
    objlist := make([]*Object, nodelen)
    for i := 0; i < nodelen; i ++ {
        objlist[i] = eval(nodes[i], env)
    }
    return objlist
}

func Eval(node parse.Node, env *Env) *Object {
    _prog, ok := node.(*parse.ListNode)
    if ! ok {
        fmt.Printf("%T\n", node)
        panic("Internal: argument to Eval() is not a ListNode!")
    }

    prog := _prog.List
    prog_len := len(prog)
    for i := 0; i < prog_len - 1; i ++ {
        eval(prog[i], env)
    }
    return eval(prog[prog_len - 1], env)
}

var DEFERED_CALL *Object = nil      // A deferred evaluation of function; used for tail-call
                                    // optimization. Is a CallNode wrapped as a macro

func defer_eval(node parse.Node, env *Env) *Object {
    switch node.(type) {
    case *parse.CallNode:
        DEFERED_CALL = NewMacro(func (args []parse.Node, env *Env) parse.Node {
            return node
        })
        return DEFERED_CALL
    default:
        return eval(node, env)
    }
}

func eval(node parse.Node, env *Env) *Object {
//start:      // For argument substitution in tail-call optimization
    switch n := node.(type) {
    // Literals
    case *parse.LiteralNode:
        obj_flag := OBJECT_NIL      // Dummy flag initializer
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
    case *WrapNode:         // Unfortunately it's here so no ``parse.''
        return n.Val
    case *parse.ListNode:
        return NewObject(OBJECT_LIST, EvalList(n.List, env))
    case *parse.DictNode:
        dict := make(map[Object]*Object)
        for k, v := range n.Dict {
            dict[*eval(k, env)] = eval(v, env)
        }
        return NewObject(OBJECT_DICT, dict)

    // Variable and references
    case *parse.SymNode:
        if n.Name == "/" || !strings.Contains(n.Name, "/") && !strings.Contains(n.Name, ".") {
            // Just a variable; no subs
            return env.GetVar(n.Name)
        }
        panic("Not implemented!")

    // Function calls
    case *parse.CallNode:
        _func := eval(n.Fun, env)       // Not really a function 'cause we don't know its type
        switch _func.Typ() {
        case OBJECT_PRIM:
            return _func.val.(func([]*Object) *Object)(EvalList(n.Arglist, env))
        case OBJECT_MACRO:
            return eval(_func.val.(func([]parse.Node, *Env) parse.Node)(n.Arglist, env), env)
        case OBJECT_FUNC:
            inner_env := NewEnv(env)
            fun := _func.val.(*Func)
            vars, body, args := fun.vars, fun.body, n.Arglist
            var_len, arg_len := len(vars), len(args)

            if var_len != arg_len {     // Check length of arguments
                panic(fmt.Sprintf("Expected %d arguments but %d were given", var_len, arg_len))
            }

            for i := 0; i < var_len; i ++ {
                inner_env.SetVarX(vars[i], eval(args[i], env))     // Create and set variable
            }

            // Run function body
            body_len := len(body)
            for i := 0; i < body_len - 1; i ++ {
                eval(fun.body[i], inner_env)
            }
            return eval(body[body_len - 1], inner_env)
            //node = body[body_len - 1]     // Substitude the ``node'' argument
            //goto start                    // And repeat the function again
        }
        panic(fmt.Sprintf("%s object can't be used as a function!", _func.Typ().String()))
    }
    panic("Not implemented!")
}
