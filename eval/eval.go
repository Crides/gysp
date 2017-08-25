package eval

import (
    "fmt"
    "github.com/Irides-Chromium/gysp/parse"
)

// Tree of Node --- Eval --> Objects
type ObjectType int

const (
    OBJECT_NIL  ObjectType = iota   // val: nil
    OBJECT_BOOL     // val: bool
    OBJECT_INT      // val: int
    OBJECT_FLOAT    // val: float64
    OBJECT_CMPLX    // Complex; val: complex128

    OBJECT_STR      // String; val: string
    OBJECT_LIST     // val: []T
    OBJECT_DICT     // val: map[interface{}]interface{}

    OBJECT_PRIM     // val: func(...Object) Object
    OBJECT_FUNC     // val: ListNode + []CallNode
    OBJECT_MACRO    // val: Func(Node) Node

    OBJECT_CLASS    // Scheme for a class declaration;
                    // val: map[string]interface{} -> map[var]initializers
    OBJECT_OBJ      // A instance of class
                    // val: map[string]interface{} -> map[var]values
)

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

func (o * Object) Val() ObjectType {
    return o.val
}

type Func struct {      // Represents a gysp function
    vars    []string
    rest    bool        // Whether the last argument is variadic
    //env     *Env        // The outer environment; for implementing closures
    body    []CallNode
}
