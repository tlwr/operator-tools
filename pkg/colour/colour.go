package colour

import (
	"fmt"
)

func foreground(s string, colour int8) string {
	return fmt.Sprintf("\033[%dm%s\033[0m", colour, s)
}

func Red(s string) string {
	return foreground(s, 31)
}

func Green(s string) string {
	return foreground(s, 32)
}

func Yellow(s string) string {
	return foreground(s, 33)
}

func Blue(s string) string {
	return foreground(s, 34)
}

func Magenta(s string) string {
	return foreground(s, 35)
}

func Cyan(s string) string {
	return foreground(s, 36)
}
