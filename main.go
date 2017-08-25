package main

import (
    "fmt"

    "github.com/Irides-Chromium/gysp/parse"
)

func main() {
    code := "(fn number? [n] (^ (int? `(n ~@a)) (is-float 'n) (complex? n))) (print 34 3.1415 4-1j 3. .5 \"\\\"Hello!\\033[0mtest\\\"\")"
    l := parse.NewLexer()
    fmt.Println(parse.Parse(l.Lex(code)))
}
