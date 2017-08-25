package env

import (
    "fmt"
    "math"
    "math/rand"
    "time"
)

// The stack
func wraps(as ...interface{}) *Stack {       // Wrap as a stack
    return NewStack(as)
}

func wrapa(as ...interface{}) []interface{} {       // Wrap as a list
    return as
}

func wrap(a interface{}) interface{} {              // Just convert type
    return a
}

type Stack []interface{}        // A wrapper for []interface{}, also acts as the main stack

func NewStack(x []interface{}) *Stack {
    new_stack := new(Stack)
    if x != nil {
        *new_stack = x
    }
    return new_stack
}

func (x * Stack) Push(a interface{}) {
    *x = append(wrapa(a), *x...)
}

func (x * Stack) Pusha(as *Stack) {
    *x = append(*as, *x...)
}

func (x * Stack) Pop() (a interface{}) {
    a, *x = (*x)[0], (*x)[1:]
    return
}

func (x * Stack) Get(ind int) interface{} {
    return (*x)[ind]
}

// The unpacking methods for usage in operators' parameter passing
func (x * Stack) Get1() interface{} {
    return (*x)[0]
}

func (x * Stack) Get2() (interface{}, interface{}) {
    return (*x)[0], (*x)[1]
}

func (x * Stack) Get3() (interface{}, interface{}, interface{}) {
    return (*x)[0], (*x)[1], (*x)[2]
}

func (x * Stack) Get4() (interface{}, interface{}, interface{}, interface{}) {
    return (*x)[0], (*x)[1], (*x)[2], (*x)[3]
}

const (
    DUMP_VERTICAL   = iota
    DUMP_HORIZONTAL
    DUMP_STRING
)

func (x * Stack) dump(opt int) {
    switch opt {
    case DUMP_VERTICAL:
        fmt.Println("\x1b[32m--- Top ---\x1b[0m")
        for _, i := range *x {
            fmt.Println(to_v(i))
        }
        fmt.Println("\x1b[32m--- Bottom ---\x1b[0m")
    case DUMP_HORIZONTAL:
        fmt.Print("\x1b[32med: stack:\x1b[0m[")
        if x.Size() > 0 {
            fmt.Print(to_v(x.Get(x.Size() - 1)))
            for i := x.Size() - 2; i >= 0 ; i -- {
                fmt.Print(" " + to_v(x.Get(i)))
            }
        }
        fmt.Println("]")
    case DUMP_STRING:
        fmt.Println(to_s(x.Contents()))
    }
}

func (x * Stack) Clear() {
    *x = []interface{}{}
}

func (x * Stack) Size() int {
    return len(*x)
}

func (x * Stack) Contents() []interface{} {
    return []interface{}(*x)
}

func (x * Stack) Switch(a, b int) {
    (*x)[a], (*x)[b] = (*x)[b], (*x)[a]
}

func (x * Stack) Reverse() {    // Reverse in-place
    for i, j := 0, x.Size() - 1; i <= j; i, j = i + 1, j - 1 {
        x.Switch(i, j)
    }
}

// The environment
type Environ struct {
    stack       *Stack
    marks       []int
    memo        *Memo
    rand        *rand.Rand
    vars        []interface{}
    longVars    map[string]interface{}
    args        []string
}

func NewEnviron(args []string) *Environ {
    env := &Environ{NewStack(nil),
                    make([]int, 0),
                    nil,
                    rand.New(rand.NewSource(time.Now().UnixNano())),
                    make([]interface{}, 0),
                    make(map[string]interface{}),
                    args}
    env.InitVars()
    return env
}

func (env * Environ) SetMemo(m * Memo) {
    env.memo = m
}

func (env * Environ) GetMemo() *Memo {
    return env.memo
}

func (env * Environ) GetRand() *rand.Rand {
    return env.rand
}

func (env * Environ) Get(ind int) interface{} {
    return env.stack.Get(ind)
}

func (env * Environ) Pop() (a interface{}) {    // Wraps the Stack.pop() method to resize the env.marks.
    a = env.stack.Pop()
    size := env.stack.Size()
    for i, m := range env.marks {
        if m > size {
            env.marks[i] = size
        }
    }
    return
}

func (env * Environ) Push(a interface{}) {
    env.stack.Push(a)
}

func (env * Environ) Pusha(as *Stack) {
    env.stack.Pusha(as)
}

func (env * Environ) Dump(opt int) {
    env.stack.dump(opt)
}

func (env * Environ) Size() int {
    return env.stack.Size()
}

// Variables
func (env * Environ) GetVar(s string) interface{} {
    if len(s) == 1 {
        if c := s[0]; is_upper(rune(c)) {
            return env.vars[c - 'A']
        }
        panic("Invalid variable name")
    }
    if a, ok := env.longVars[s]; ok {
        return a
    }
    panic("Variable not defined!")
}

func (env * Environ) SetVar(s string, a interface{}) {
    if len(s) == 1 {
        if c := s[0]; is_upper(rune(c)) {
            env.vars[c - 'A'] = a
            return
        }
        panic("Invalid variable name")
    }
    env.longVars[s] = a
}

func (env * Environ) InitVars() {
}
