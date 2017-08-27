package main

import (
    "fmt"

    "github.com/crides/gysp/parse"
    "github.com/crides/gysp/eval"
)

func main() {
    code := `(debug (if nil {"yes" true} {"no" false}))`
    l := parse.NewLexer()
    fmt.Println(eval.Eval(parse.Parse(l.Lex(code)), eval.StandardEnv()))
}
