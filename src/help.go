package prockeeper

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

const usage = `usage: prockeeper [options]

  --help          Show this help
  -c path_to_yml  Specify the path of yaml file (default: './prockeeper.yml')

Example yaml:
  services:
    - name: "rails server"
      command: "rails s"
    - name: "node server"
      command: "npm start"
`

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

// Usage ...
func Usage() {
	fmt.Fprint(os.Stdout, usage)
	fmt.Fprint(os.Stdout, helpMessage)
	os.Exit(0)
}

// HelpMenu ...
func HelpMenu() *tview.TextView {
	help := tview.NewTextView()
	help.SetBorder(true).SetBackgroundColor(tcell.ColorDarkSlateGrey).SetTitle("Help")
	fmt.Fprintf(help, "%s", helpMessage)
	return help
}
