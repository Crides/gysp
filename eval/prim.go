package eval

import (
    "fmt"
    "math"
    "github.com/crides/gysp/parse"
)

func NewMacro(m func([]parse.Node, *Env) parse.Node) *Object {
    return NewObject(OBJECT_MACRO, m)
}

func NewPrim(f func([]*Object) *Object) *Object {
    return NewObject(OBJECT_PRIM, f)
}

func NewInt(i int) *Object {
    return NewObject(OBJECT_INT, i)
}

// Helper methods
func convert_err(t1, t2 ObjectType) *Object {
    panic(fmt.Sprintf("Cannot convert %v to %v!", t1, t2))
    return nil      // Easier to suppress warnings
}

func nomethod_err1(meth string, t ObjectType) *Object {
    panic(fmt.Sprintf("No '%s' method for type '%v'!", meth, t))
    return nil
}

func nomethod_err2(meth string, t1, t2 ObjectType) *Object {
    panic(fmt.Sprintf("No '%s' method for type '%v' and '%v'!", meth, t1, t2))
    return nil
}

func to_float(o *Object) *Object {
    switch o.typ {
    case OBJECT_FLOAT:
        return o
    case OBJECT_INT:
        return NewObject(OBJECT_FLOAT, float64(o.val.(int)))
    }
    return convert_err(o.typ, OBJECT_FLOAT)
}

func to_cmplx(o *Object) *Object {
    switch o.typ {
    case OBJECT_CMPLX:
        return o
    case OBJECT_FLOAT:
        return NewObject(OBJECT_CMPLX, complex(o.val.(float64), 0))
    case OBJECT_INT:
        return NewObject(OBJECT_CMPLX, complex(float64(o.val.(int)), 0))
    }
    return convert_err(o.typ, OBJECT_CMPLX)
}

func add(a, b *Object) *Object {
    typ := a.typ
    if typ != b.typ {
        panic("intern-add: Arguments must be the same type!")
    }

    switch typ {
    case OBJECT_INT:
        return NewObject(OBJECT_INT, a.val.(int) + b.val.(int))
    case OBJECT_FLOAT:
        return NewObject(OBJECT_FLOAT, a.val.(float64) + b.val.(float64))
    case OBJECT_CMPLX:
        return NewObject(OBJECT_CMPLX, a.val.(complex128) + b.val.(complex128))
    case OBJECT_STR:
        return NewObject(OBJECT_STR, a.val.(string) + b.val.(string))
    }
    return GYSP_NIL
}

func negate(a *Object) *Object {
    switch a.typ {
    case OBJECT_INT:
        return NewObject(OBJECT_INT, -a.val.(int))
    case OBJECT_FLOAT:
        return NewObject(OBJECT_FLOAT, -a.val.(float64))
    case OBJECT_CMPLX:
        return NewObject(OBJECT_CMPLX, -a.val.(complex128))
    case OBJECT_STR:
        panic("No negation for type 'str'!")
    }
    return GYSP_NIL
}

func sub(a, b *Object) *Object {
    typ := a.typ
    if typ != b.typ {
        panic("intern-sub: Arguments must be the same type!")
    }

    switch typ {
    case OBJECT_INT, OBJECT_FLOAT, OBJECT_CMPLX, OBJECT_STR:
        return NewObject(OBJECT_INT, add(a, negate(b)))
    }
    return GYSP_NIL
}

func mul(a, b *Object) *Object {
    typ := a.typ
    if typ != b.typ {
        panic("intern-mul: Arguments must be the same type!")
    }

    switch typ {
    case OBJECT_INT:
        return NewObject(OBJECT_INT, a.val.(int) * b.val.(int))
    case OBJECT_FLOAT:
        return NewObject(OBJECT_FLOAT, a.val.(float64) * b.val.(float64))
    case OBJECT_CMPLX:
        return NewObject(OBJECT_CMPLX, a.val.(complex128) * b.val.(complex128))
    case OBJECT_STR:
        panic("No multiplication for type 'str'!")
    }
    return GYSP_NIL
}

func div(a, b *Object) *Object {
    typ := a.typ
    if typ != b.typ {
        panic("intern-div: Arguments must be the same type!")
    }

    switch typ {
    case OBJECT_INT:
        return NewObject(OBJECT_INT, a.val.(int) / b.val.(int))
    case OBJECT_FLOAT:
        return NewObject(OBJECT_FLOAT, a.val.(float64) / b.val.(float64))
    case OBJECT_CMPLX:
        return NewObject(OBJECT_CMPLX, a.val.(complex128) / b.val.(complex128))
    case OBJECT_STR:
        panic("No division for type 'str'!")
    }
    return GYSP_NIL
}

func mod(a, b *Object) *Object {
    typ := a.typ
    if typ != b.typ {
        panic("intern-mod: Arguments must be the same type!")
    }

    switch typ {
    case OBJECT_INT:
        return NewObject(OBJECT_INT, a.val.(int) % b.val.(int))
    case OBJECT_FLOAT:
        return NewObject(OBJECT_FLOAT, math.Mod(a.val.(float64), b.val.(float64)))
    case OBJECT_CMPLX:
        panic("No modulus for type 'complex'!")
    case OBJECT_STR:
        panic("No modulus for type 'str'!")
    }
    return GYSP_NIL
}

func pow(a, b *Object) *Object {
    typ := a.typ
    if typ != b.typ {
        panic("intern-pow: Arguments must be the same type!")
    }

    switch typ {
    case OBJECT_INT:
        return NewObject(OBJECT_INT, a.val.(int) % b.val.(int))
    case OBJECT_FLOAT:
        return NewObject(OBJECT_FLOAT, math.Pow(a.val.(float64), b.val.(float64)))
    case OBJECT_CMPLX:
        panic("No power for type 'complex'!")
    case OBJECT_STR:
        panic("No power for type 'str'!")
    }
    return GYSP_NIL
}

func Range(start, end, step int) *Object {
    list := make([]*Object, (end - start) / step)
    j := 0      // Just a counter
    for i := start; i < end; i += step {
        list[j] = NewInt(i)
        j ++
    }
    return NewObject(OBJECT_LIST, list)
}

func StandardEnv() *Env {
    return &Env{
        map[string]*Object {
        // Constants
        "nil": GYSP_NIL,
        "true": GYSP_TRUE,
        "false": GYSP_FALSE,

        // Primitives
        "+": NewPrim(func (args []*Object) *Object {
            sum := args[0]
            for i := 1; i < len(args); i ++ {
                sum = add(sum, args[i])
            }
            return sum
        }),
        "-": NewPrim(func (args []*Object) *Object {
            sum := args[0]
            for i := 1; i < len(args); i ++ {
                sum = sub(sum, args[i])
            }
            return sum
        }),
        "*": NewPrim(func (args []*Object) *Object {
            product := args[0]
            for i := 1; i < len(args); i ++ {
                product = mul(product, args[i])
            }
            return product
        }),
        "/": NewPrim(func (args []*Object) *Object {
            product := args[0]
            for i := 1; i < len(args); i ++ {
                product = div(product, args[i])
            }
            return product
        }),
        "%": NewPrim(func (args []*Object) *Object {
            if len(args) != 2 {
                panic("Invalid number of arguments passed to '%'!")
            }
            return mod(args[0], args[1])
        }),

        "range": NewPrim(func (args []*Object) *Object {
            switch len(args) {
            case 1:     // Only stop
                return Range(0, args[0].val.(int), 1)
            case 2:     // Start and stop
                return Range(args[0].val.(int), args[1].val.(int), 1)
            case 3:     // Start, stop and step
                return Range(args[0].val.(int), args[1].val.(int), args[2].val.(int))
            }
            panic("Invalid arguemnts for range()!")
        }),
        "print": NewPrim(func (args []*Object) *Object {
            converted := make([]interface{}, len(args))
            for i := 0; i < len(args); i ++ {
                converted[i] = args[i]
            }
            fmt.Print(converted...)
            return GYSP_NIL
        }),
        "println": NewPrim(func (args []*Object) *Object {
            converted := make([]interface{}, len(args))
            for i := 0; i < len(args); i ++ {
                converted[i] = args[i]
            }
            fmt.Println(converted...)
            return GYSP_NIL
        }),

        "if": NewMacro(func (args []parse.Node, env *Env) parse.Node {
            if eval(args[0], env) != GYSP_NIL {
                return args[1]
            }
            if len(args) > 2 {
                return args[2]
            }
            return parse.NIL_NODE
        }),

        "for": NewMacro(func (args []parse.Node, env *Env) parse.Node {
            _ranges, ok := args[0].(*parse.ListNode)
            if ! ok {
                panic("Expected list!")
            }

            ranges := _ranges.List
            range_len := len(ranges)
            if range_len % 2 != 0 {
                panic("Range list must have a even number of items!")
            }
            vars, lists := make([]string, range_len / 2), make([]*Object, range_len / 2)
            for i := 0; i < range_len / 2; i ++ {
                vars[i] = ranges[2 * i].(*parse.SymNode).Name
                lists[i] = eval(ranges[2 * i + 1], env)
                if lists[i].typ != OBJECT_LIST {
                    panic("Range source must be a list!")
                }
            }

            inner_env := NewEnv(env)
            gen := NewGenerator()
            for _, list := range lists {
                gen.AddSource(list.val.([]*Object))
            }
            gen.Generate(func(as []*Object) *Object {
                // Set loop variables
                for i := 0; i < len(vars); i ++ {
                    inner_env.SetVarX(vars[i], as[i])
                }
                // Run body
                for i := 1; i < len(args) - 1; i ++ {
                    eval(args[i], inner_env)
                }
                return eval(args[len(args) - 1], inner_env)     // TODO print results on repl
            })
            return parse.NIL_NODE
        }),

        "let": NewMacro(func (args []parse.Node, env *Env) parse.Node {
            _bindings, ok := args[0].(*parse.ListNode)
            if ! ok {
                panic("Expected list!")
            }

            bindings := _bindings.List
            bind_len := len(bindings)
            if bind_len % 2 != 0 {
                panic("Binding list must have a even number of items!")
            }

            // Set bindings
            inner_env := NewEnv(env)
            for i := 0; i < bind_len / 2; i ++ {
                inner_env.SetVarX(
                    bindings[2 * i].(*parse.SymNode).Name,
                    eval(bindings[2 * i + 1], env))
            }

            // Run body
            for i := 1; i < len(args) - 1; i ++ {       // Only the calls before the last
                eval(args[i], inner_env)
            }
            return WrapObject(eval(args[len(args) - 1], inner_env))
        }),

        "do": NewMacro(func (args []parse.Node, env *Env) parse.Node {
            // Create inner environment
            inner_env := NewEnv(env)
            // Run body
            for i := 0; i < len(args) - 1; i ++ {
                eval(args[i], inner_env)
            }
            return WrapObject(eval(args[len(args) - 1], inner_env))
        }),
    }, nil}
}
