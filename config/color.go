package config

import "github.com/gdamore/tcell"

type (
	Color string
)

const (
	DefaultColor Color = "default"
)

func (c Color) Color() tcell.Color {
	if c == DefaultColor {
		return tcell.ColorDefault
	}

	return tcell.GetColor(string(c))
}
