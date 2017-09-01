package parse

import (
    re "regexp"
    "strings"
    "strconv"
    "fmt"
)

// code --- Lexer --> []*Node

type TokenType int

const (
    TOKEN_NONE TokenType = iota - 1 // -1
    FUNC_BEGIN      // 0
    FUNC_END
    LIST_BEGIN
    LIST_END
    DICT_BEGIN
    DICT_END

    STRING          // 6
    COMPLEX
    FLOAT
    INTEGER

    QUOTE           // 11
    QQUOTE
    UNQUOTESP       // Put unquote-splice before unquote for parser optimization
    UNQUOTE

    TOKEN
)

type Token struct {
    typ     TokenType
    cont    string
}

func NewToken(typ TokenType, tok string) *Token {
    return &Token{typ, tok}
}

func (t * Token) Typ() TokenType {
    return t.typ
}

func (t * Token) Cont() string {
    return t.cont
}

func (t * Token) String() string {
    return fmt.Sprintf("{%d %s}", t.typ, t.cont)
}

type Lexer struct {
    pats        []*re.Regexp
    types       []TokenType
    ignores     []*re.Regexp

    replacer    *strings.Replacer
    repats      []*re.Regexp        // Patterns to handle \777, \xFF, \uFFFF and \UFFFFFFFFs
    refuncs     []func(string)string
}

func NewLexer() *Lexer {
    l := new(Lexer)

    // Can be changed to adopt other situations
    regexs := []string{
        // Raw strings are better!
        `\(`, `\)`, `\[`, `\]`, `\{`, `\}`, // Brackets

        `"(?:[^"]|(?:\"))*?"`,              // String
        `[+-]?(?:\d*\.)?\d+[+-]?(?:\d*\.)?\d+j`,    // Complex
        `[+-]?\d*\.\d+`,                    // Float
        `[+-]?\d+`,                         // Integer

        `'(?:[^])}\s])`,                    // Quotes
        "`(?:[^])}\\s])",   // ... except this
        `~@(?:[^])}\s])`,
        `~(?:[^])}\s])`,

        `[^][(){}\s'";]+`,                  // Token
    }
    for i, pat := range regexs {
        l.AddRegexp(TokenType(i), pat)
    }
    l.Ignore("\\s+")    // Ignore spaces
    l.Ignore(";.*$")    // Ignore comments

    l.replacer = strings.NewReplacer(
        `\"`, `"`,      // Replace quote escapes
        `\a`, "\a",     // C escape characters
        `\b`, "\b",
        `\e`, "\x1b",   // GNU extension
        `\f`, "\f",
        `\n`, "\n",
        `\r`, "\r",
        `\t`, "\t",
        `\v`, "\v",
    )
    l.repats = []*re.Regexp{
        re.MustCompile(`\\x[0-9a-fA-F]{0,2}`),
        re.MustCompile(`\\u[0-9a-fA-F]{0,4}`),
        re.MustCompile(`\\U[0-9a-fA-F]{0,8}`),
        re.MustCompile(`\\[0-7]{0,3}`),
    }
    l.refuncs = []func(string)string {
        func(s string) string {
            if len(s) != 4 {
                panic(`\x must be followed by exactly 2 hexdigits!`)
            }
            i, _ := strconv.ParseInt(s[2:], 16, 8)
            return string(i)
        },

        func(s string) string {
            if len(s) != 6 {
                panic(`\u must be followed by exactly 4 hexdigits!`)
            }
            i, _ := strconv.ParseInt(s[2:], 16, 32)
            return string(i)
        },

        func(s string) string {
            if len(s) != 10 {
                panic(`\u must be followed by exactly 8 hexdigits!`)
            }
            i, _ := strconv.ParseInt(s[2:], 16, 32)
            return string(i)
        },

        func(s string) string {
            if len(s) != 4 {
                panic(`\ must be followed by exactly 3 octdigits!`)
            }
            i, _ := strconv.ParseInt(s[1:], 8, 8)
            return string(i)
        },
    }
    return l
}

func (l * Lexer) AddRegexp(typ TokenType, pat string) {
    l.pats = append(l.pats, re.MustCompile(pat))
    l.types = append(l.types, typ)
}

func (l * Lexer) Ignore(pat string) {
    l.ignores = append(l.ignores, re.MustCompile(pat))
}

func (l * Lexer) ProcessString(s string) string {
    s = l.replacer.Replace(s[1:len(s) - 1])
    for i := 0; i < 4; i ++ {
        s = l.repats[i].ReplaceAllStringFunc(s, l.refuncs[i])
    }
    return s
}

func (l * Lexer) Lex(code string) (toks []*Token) {
    for inds := []int{0, 0}; len(code) > 0; code = code[inds[1]:] {
        for i, pat := range l.pats {
            if inds = pat.FindStringIndex(code); inds != nil && inds[0] == 0 && inds[1] > 0 {
                typ, token := l.types[i], pat.FindString(code)
                if typ >= QUOTE && typ <= UNQUOTE {    // Belong to the quotes
                    inds[1] --      // Get rid of the last char
                    token = token[:len(token) - 1]
                } else if typ == STRING {
                    token = l.ProcessString(token)
                }
                toks = append(toks, NewToken(typ, token))
                goto next
            }
        }
        for _, pat := range l.ignores {
            if inds = pat.FindStringIndex(code); inds != nil && inds[0] == 0 && inds[1] > 0 {
                goto next
            }
        }
        if inds == nil {
            panic("Cannot identify the next token!")
        }
next:
    }
    return
}
