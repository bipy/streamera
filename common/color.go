package common

import (
    "errors"
    "fmt"
)

type Color int

const (
    BLACK  = 30
    RED    = 31
    GREEN  = 32
    YELLOW = 33
    BLUE   = 34
    PURPLE = 35
    CYAN   = 36
    GRAY   = 37
)

func getColorStr(s string, color Color) (string, error) {
    if color < 30 || color > 37 {
        return "", errors.New("font [30, 37]")
    }
    rt := fmt.Sprintf("\033[%dm%s\033[0m", color, s)
    return rt, nil
}

func GetColorStr(s string, color Color) string {
    rt, err := getColorStr(s, color)
    if err != nil {
        panic(err)
    }
    return rt
}

func Black(s string) string {
    return GetColorStr(s, BLACK)
}

func Green(s string) string {
    return GetColorStr(s, GREEN)
}

func Red(s string) string {
    return GetColorStr(s, RED)
}

func Yellow(s string) string {
    return GetColorStr(s, YELLOW)
}

func Blue(s string) string {
    return GetColorStr(s, BLUE)
}

func Purple(s string) string {
    return GetColorStr(s, PURPLE)
}

func Cyan(s string) string {
    return GetColorStr(s, CYAN)
}

func Gray(s string) string {
    return GetColorStr(s, GRAY)
}
