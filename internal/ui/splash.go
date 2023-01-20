package ui

import (
	"fmt"
	"strings"

	"github.com/derailed/tview"
)

// LogoSmall hcloud small log.
var LogoSmall = []string{
	`.__           .__                   .___`,
	`|  |__   ____ |  |   ____  __ __  __| _/`,
	`|  |  \_/ ___\|  |  /  _ \|  |  \/ __ | `,
	`|   Y  \  \___|  |_(  <_> )  |  / /_/ | `,
	`|___|  /\___  >____/\____/|____/\____ | `,
	`     \/     \/                       \/ `,
}

// LogoBig hcloud big logo for splash page.
var LogoBig = []string{
	`.__           .__                  __           .__                   .___`,
	`|  |__   ____ |  |   _____  __ ___/  |_    ____ |  |   ____  __ __  __| _/`,
	`|  |  \_/ __ \|  |  /     \|  |  \   __\ _/ ___\|  |  /  _ \|  |  \/ __ | `,
	`|   Y  \  ___/|  |_|  Y Y  \  |  /|  |   \  \___|  |_(  <_> )  |  / /_/ | `,
	`|___|  /\___  >____/__|_|  /____/ |__|   /\___  >____/\____/|____/\____ | `,
	`    \/     \/           \/              \/   \/                       \/`,
}

// Splash represents a splash screen.
type Splash struct {
	*tview.Flex
}

// NewSplash instantiates a new splash screen with product and company info.
func NewSplash() *Splash {
	s := Splash{Flex: tview.NewFlex()}

	logo := tview.NewTextView()
	logo.SetDynamicColors(true)
	logo.SetTextAlign(tview.AlignCenter)
	s.layoutLogo(logo)

	s.SetDirection(tview.FlexRow)
	s.AddItem(logo, 10, 1, false)

	return &s
}

func (s *Splash) layoutLogo(t *tview.TextView) {
	logo := strings.Join(LogoBig, fmt.Sprint("\n"))
	fmt.Fprintf(t, "%s%s\n", strings.Repeat("\n", 2), logo)
}
