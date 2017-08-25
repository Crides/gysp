package color

import (
    //"fmt"
    "strconv"
)

func csi(s string) string {
    return "\x1b[" + s
}

func color_8(n int, bold bool) string {
    a := ""
    if bold {
        a = "1;"
    }
    return csi(a + "3" + strconv.Itoa(n) + "m")
}

func Clear() string {
    return color_8(7, false)
}

func Teal(msg string) string {
    return color_8(0, false) + msg + Clear()
}

func Red(msg string) string {
    return color_8(1, false) + msg + Clear()
}

func Green(msg string) string {
    return color_8(2, false) + msg + Clear()
}

func Yellow(msg string) string {
    return color_8(3, false) + msg + Clear()
}

func Blue(msg string) string {
    return color_8(4, false) + msg + Clear()
}

func Mangenta(msg string) string {
    return color_8(5, false) + msg + Clear()
}

func Turqoise(msg string) string {
    return color_8(6, false) + msg + Clear()
}
