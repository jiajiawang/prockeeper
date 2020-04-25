package prockeeper

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

const helpMessage = `
Keyboard commands

j      - Select previous item
k      - Select next item
Enter  - Start/stop selected service
u      - Start all services
d      - Stop all services

?      - Show/hide help menu
.      - Show/hide debugger
Ctrl-C - Exit app
`

// HelpMenu ...
func HelpMenu() *tview.TextView {
	help := tview.NewTextView()
	help.SetBorder(true).SetBackgroundColor(tcell.ColorDarkSlateGrey).SetTitle("Help")
	fmt.Fprintf(help, "%s", helpMessage)
	return help
}
