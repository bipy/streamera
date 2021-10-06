package common

import (
    "errors"
    "fmt"
)

type Color struct {
    font       int
    background int
    display    int
}

var (
    RED = Color{font: 31, background: 40, display: 1}

    GREEN = Color{font: 32, background: 40, display: 1}

    YELLOW = Color{font: 33, background: 40, display: 1}

    BLUE = Color{font: 34, background: 40, display: 1}
)

func getColorStr(s string, color Color) (string, error) {
    if color.font < 30 || color.font > 37 {
        return "", errors.New("font [30, 37]")
    }
    if color.background < 40 || color.background > 47 {
        return "", errors.New("background [40, 47]")
    }
    if color.display < 0 || color.display > 8 || color.display == 2 || color.display == 3 || color.display == 6 {
        return "", errors.New("display {0, 1, 4, 5, 7, 8}")
    }
    rt := fmt.Sprintf("\033[%d;%d;%dm%s\033[0m", color.display, color.font, color.background, s)
    return rt, nil
}

func GetColorStr(s string, color Color) string {
    rt, err := getColorStr(s, color)
    if err != nil {
        panic(err)
    }
    return rt
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