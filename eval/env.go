package eval

import (
    "fmt"
)

type Env struct {
	scope   map[string]*Object
	next    *Env
}

func NewEnv(outer *Env) *Env {
    return &Env{make(map[string]*Object), outer}
}

//func (e * Env) Push() {
//    // Creates a new scope (pushes a stack frame)
//    // The new scope is the first
//    e.scopes = append([]Scope{make(Scope)}, e.scopes...)
//}
//
//func (e * Env) Pop() {
//    e.scopes = e.scopes[1:]
//}

func (e * Env) NewVar(vname string) {  // Creates a new variable in the current scope
    e.scope[vname] = nil
}

func (e * Env) GetVar(vname string) *Object {
    for ; e != nil; e = e.next {
        if item, ok := e.scope[vname]; ok {
            return item
        }
    }
    panic(fmt.Sprintf("Variable %s not defined!", vname))
}

func (e * Env) SetVar(vname string, val *Object) {
    for ; e != nil; e = e.next {
        if _, ok := e.scope[vname]; ok {
            e.scope[vname] = val
            return
        }
    }
    panic(fmt.Sprintf("Variable %s not defined!", vname))
}

func (e * Env) SetVarX(vname string, val *Object) {
    // In the current scope, set the var named vname
    // Don't check for variable existence
    e.scope[vname] = val
}

func StandardEnv() *Env {
    return &Env{map[string]*Object{
        // Constants
        "nil": GYSP_NIL,
        "true": GYSP_TRUE,
        "false": GYSP_FALSE,

        // Primitives
        "print": NewObject(OBJECT_PRIM, func (args ...*Object) *Object {
            fmt.Print(args)
            return GYSP_NIL
        }),
        "println": NewObject(OBJECT_PRIM, func (args ...*Object) *Object {
            fmt.Println(args)
            return GYSP_NIL
        }),
        "debug": NewObject(OBJECT_PRIM, func (args ...*Object) *Object {
            for _, o := range args {
                fmt.Println(o)
            }
            return GYSP_NIL
        }),
        "if": NewObject(OBJECT_PRIM, func (args ...*Object) *Object {
            if args[0] != GYSP_NIL {
                return args[1]
            }
            if len(args) > 2 {
                return args[2]
            }
            return GYSP_NIL
        }),
    }, nil}
}

//func (e * Env) Freeze() *Env {
//    // Freezes the current environment for function objects
//    frozen := new(Env)
//    frozen.scopes = make([]Scope, len(e.scopes))
//    copy(frozen.scopes, e.scopes)
//    return frozen
//}
