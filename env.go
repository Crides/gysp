package eval

import (

)

type Scope map[string]Object

type Env struct {
    scopes    []Scope
}

func NewEnv() (e *Env) {
    e.Push()
    return
}

func (e * Env) Push() {
    // Creates a new scope (pushes a stack frame)
    // The new scope is the first
    e.scopes = append([]Scope{make(Scope)}, e.scopes...)
}

func (e * Env) Pop() {
    e.scopes = e.scopes[1:]
}

func (e * Env) NewVar(vname string) {  // Creates a new variable in the current scope
    e.scopes[0][vname] = nil
}

func (e * Env) GetVar(vname string) Object {
    for _, scope := range e.scopes {
        if item, ok := scope[vname]; ok {
            return item
        }
    }
    panic("Variable not defined!")
}

func (e * Env) SetVar(vname string, val Object) {
    for _, scope := range e.scopes {
        if item, ok := scope[vname]; ok {
            e.scopes[i][vname] = val
            return
        }
    }
    panic("Variable not defined!")
}

func (e * Env) SetVarX(vname string, val Object) {
    // In the current scope, set the var named vname
    // Don't check for variable existence
    e.scopes[0][vname] = val
}

func (e * Env) Freeze() *Env {
    // Freezes the current environment for function objects
    frozen := new(Env)
    frozen.scopes = make([]Scope, len(e.scopes))
    copy(frozen.scopes, e.scopes)
    return frozen
}
