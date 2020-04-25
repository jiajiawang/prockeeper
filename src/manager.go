package prockeeper

import (
	"flag"
	"log"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// Manager ...
type Manager struct {
	Services []*Service
	list     *tview.List
	logger   *log.Logger
}

// func NewManager() *Manager {
// list := tview.NewList().ShowSecondaryText(false)
// list.SetTitle("Services (Press ? to show help)").SetBorder(true)

// return &Manager{
// list: tview.NewList
// }
// }

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

var debug bool

func init() {
	flag.BoolVar(&debug, "debug", false, "debug mode")
	flag.Parse()
}

func (manager *Manager) refreshList() {
	currentSelection := manager.list.GetCurrentItem()
	manager.list.Clear()
	for _, s := range manager.Services {
		manager.list.AddItem(s.NameWithPid(), "", 0, nil)
	}
	manager.list.SetCurrentItem(currentSelection)
}

func (manager *Manager) startService(i int) {

}

func (manager *Manager) startAll() {
	for _, s := range manager.Services {
		go func(s *Service) {
			if err := s.Start(); err != nil {
				manager.logger.Println(err)
			}
		}(s)
	}
}

func (manager *Manager) stopAll() {
	for _, s := range manager.Services {
		go func(s *Service) {
			if err := s.Stop(); err != nil {
				manager.logger.Println(err)
			}
		}(s)
	}
}

func (manager *Manager) selectService(i int) {

}

// Run ...
func (manager *Manager) Run() {
	config := ParseConfig()

	updated := make(chan struct{})

	for _, s := range config.Services {
		service := &Service{
			Name:    s.Name,
			Command: s.Command,
			Updated: updated,
		}
		manager.Services = append(manager.Services, service)
	}

	app := tview.NewApplication()

	modal := func(p tview.Primitive, width, height int) tview.Primitive {
		return tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(p, height, 1, false).
				AddItem(nil, 0, 1, false), width, 1, false).
			AddItem(nil, 0, 1, false)
	}

	help := HelpMenu()

	layout := tview.NewFlex().SetDirection(tview.FlexRow)
	pages := tview.NewPages().
		AddPage("app", layout, true, true).
		AddPage("help", modal(help, 40, 20), true, false)
	app.SetRoot(pages, true)

	appContainer := tview.NewFlex().SetDirection(tview.FlexRow)
	debuggerContainer := tview.NewFlex()

	layout.AddItem(appContainer, 0, 5, false)

	if debug {
		layout.AddItem(debuggerContainer, 0, 1, false)
	}

	debugger := tview.NewTextView().
		SetDynamicColors(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	debugger.SetTitle("debugger").SetBorder(true)

	debuggerContainer.AddItem(debugger, 0, 1, true)
	logger := log.New(debugger, "", log.LstdFlags)
	manager.logger = logger

	list := tview.NewList().ShowSecondaryText(false)
	manager.list = list
	list.SetTitle("Services (Press ? to show help)").SetBorder(true)
	for _, s := range manager.Services {
		s.Prepare(app, logger)
		list.AddItem(s.NameWithPid(), "", 0, nil)
	}

	appContainer.AddItem(list, 0, 2, true)
	currentView := manager.Services[0].LogView
	appContainer.AddItem(currentView, 0, 6, true)

	list.SetChangedFunc(func(index int, n string, v string, t rune) {
		if currentView != nil {
			appContainer.RemoveItem(currentView)
		}

		view := manager.Services[index].LogView
		appContainer.AddItem(view, 0, 6, true)
		currentView = view
	})

	list.SetSelectedFunc(func(i int, n string, v string, t rune) {
		manager.Services[i].Toggle()
	})

	exitMenu := tview.NewModal().
		SetText("Running services!").
		AddButtons([]string{"Force Quit", "Cancel"}).
		SetFocus(1).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Force Quit" {
				app.Stop()
			} else {
				pages.RemovePage("exit")
				app.SetFocus(list)
			}
		})

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == '?' {
			if name, _ := pages.GetFrontPage(); name == "help" {
				pages.HidePage("help")
				app.SetFocus(list)
			} else {
				pages.ShowPage("help")
			}
			return nil
		}
		if event.Rune() == '.' {
			if debug {
				layout.RemoveItem(debuggerContainer)
				debug = false
			} else {
				layout.AddItem(debuggerContainer, 0, 1, false)
				debug = true
			}
		}
		if event.Rune() == 'j' {
			return tcell.NewEventKey(tcell.KeyDown, 'j', tcell.ModNone)
		}
		if event.Rune() == 'k' {
			return tcell.NewEventKey(tcell.KeyUp, 'k', tcell.ModNone)
		}
		if event.Key() == tcell.KeyCtrlC {
			allStopped := true
			for _, s := range manager.Services {
				if s.Cmd != nil {
					allStopped = false
					break
				}
			}
			if allStopped {
				return event
			}
			pages.AddPage("exit", exitMenu, true, true)
			return nil
		}
		if event.Rune() == 'u' {
			manager.startAll()
			return nil
		}
		if event.Rune() == 'd' {
			manager.stopAll()
			return nil
		}
		return event
	})

	go func() {
		for range updated {
			logger.Println("refresh list")
			manager.refreshList()
		}
	}()

	if err := app.SetFocus(list).Run(); err != nil {
		panic(err)
	}
}
