package ui

import (
	"fmt"

	"github.com/derailed/tview"
	"github.com/gdamore/tcell/v2"
)

// Logo represents a K9s logo.
type Logo struct {
	*tview.Flex

	logo, status *tview.TextView
}

// NewLogo returns a new logo.
func NewLogo() *Logo {
	l := Logo{
		Flex: tview.NewFlex(),
		logo: logo(),
	}
	l.SetDirection(tview.FlexRow)
	l.AddItem(l.logo, 6, 1, false)

	for i, s := range LogoSmall {
		fmt.Fprintf(l.logo, "[#%06x::b]%s", tcell.ColorDarkOrange.TrueColor().Hex(), s)
		if i+1 < len(LogoSmall) {
			fmt.Fprintf(l.logo, "\n")
		}
	}

	return &l
}

func logo() *tview.TextView {
	v := tview.NewTextView()
	v.SetWordWrap(false)
	v.SetWrap(false)
	v.SetTextAlign(tview.AlignLeft)
	v.SetDynamicColors(true)

	return v
}
