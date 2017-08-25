# Gysp

Gysp is lisp similar to [Hy](https://github.com/hylang/hy) but implemented in Go.

## Spec

### Goal

The goal of Gysp is to provide a simple implementation of lisp so that it is easier to write and read.

#### Features

1. Use `[]` literal to represent lists (or arrays or vectors) instead of quoted S-expressions (<code>&#96;(a b c)</code>)
2. Use `{}` literal to represent dictionaries (from Hy, or hashmaps or tables).
3. Macros are continue to be supported, because they are useful in making code more readable.
4. Lists are used in place of many S-expressions (from Hy). For example, in Hy, function definition would be:
```Hy
    (defn f [a b] (+ a b))
```
instead of:
```lisp
    (defun f (a b) (+ a b))
```
5. Use snake case, but with underscore replaced by hyphens. That makes it easier to type names, but readability stays the same. Earmuffs are still used.

#### Functions

1. (use ...)
Import a package. Arguments are strings. There are two kinds of import: 1) absolute import, and 2) relative import. Absolute import is used for packages installed under PATH, and relative import is used for files on a specified path. Examples:
```lisp
    (use "package.class.name" "./file/path/to/package")
```
2. (set ...)
Setting references' values (from Hy). The number of arguments must be multiples of 2. So the odd-numbered ones are the references, and the even-numbered ones are the values to set. If a variable is undefined, then it's created in the current scope and set. `setv`, `setq` and `setf` are not used because they can be confusing when used, and increased key-strokes. Example:
```lisp
    (set foo 1 bar 2)
```

3. (get coll ind)
Get a reference from a collection with the specified index. The collection can be a list, string, or a dictionary. The reference can then be set using the `set` function. Example:
```lisp
    (set (get arr 3) 4)
```
is equivalent to `arr[3] = 4` in Python.

4. (printf format ...)
Print the arguments with format to console, like the `printf` function in C.

5. (println ...)
Print the representation of the values to console, and append a newline. Like the builtin `print()` function in Python.

6. Some Python builtin functions:
```Python
    str()
    int()
    float()
    format()
```
(of course you need to call them in the lisp style). The `format()` function is a little bit different. It takes more than one argument and acts like the `sprintf` function in C.

7. Some functional functions:
```list
    (curry func arg)
    (first coll)
    (last coll)
    (second coll)
```

8. Macros like the ones in Hy:
```lisp
    (defm ...)
    (defr ...)
```

9. OOP functions:
```lisp
    (defc ...)
```
