package main

import (
    "fmt"
    "os"
    "bufio"

    "github.com/crides/gysp/parse"
    "github.com/crides/gysp/eval"
    "github.com/crides/gysp/color"
)

func main() {
    Repl()
    //code := `(println (+ 1 2 3) (* 2-3j 1+2j) (/ 2-3j 1+2j))`
    //l := parse.NewLexer()
    //fmt.Println(eval.Eval(parse.Parse(l.Lex(code)), eval.StandardEnv()))
}

// Wrapper for the eval.Eval function
func Eval(code string, lexer *parse.Lexer, env *eval.Env) *eval.Object {
    return eval.Eval(parse.Parse(lexer.Lex(code)), env)
}

func Repl() {
    // Repl constants
    header := "Gysp 1.0 by Steven."
    PS1 := " => "
    //PS2 := "... "

    // Environments
    input := bufio.NewScanner(os.Stdin)
    lexer := parse.NewLexer()
    env := eval.StandardEnv()

    fmt.Println(header)
    fmt.Print(PS1)
    for input.Scan() {
        fmt.Println(color.Yellow("output:"))
        ret_val := Eval(input.Text(), lexer, env)
        fmt.Print(color.Green("returned: "))
        fmt.Println(ret_val.GoString())

        fmt.Print(PS1)
    }
    fmt.Println("bye!")
}
