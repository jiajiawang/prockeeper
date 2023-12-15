package prockeeper

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const usage = `usage: prockeeper [options]

  --help          Show this help
  --debug         Write application log to ./prockeeper.log
  -c path_to_yml  Specify the path of yaml file (default: './prockeeper.yml')

Service Options:
  [name]    Specify the name of the service
  [command] Specify the exec command
  [dir]     Specify the working directory

Example yaml:
  services:
    - name: "rails server"
      command: "rails s"
    - name: "node server"
      command: "npm start"
      dir: "./client"
`

const helpMessage = `
Keyboard commands

j      - Select previous item
k      - Select next item
Enter  - Start/stop selected service
u      - Start all services
d      - Stop all services

?      - Show/hide help menu
.      - Show/hide application log
,      - Show/hide process log
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
