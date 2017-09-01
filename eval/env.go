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
