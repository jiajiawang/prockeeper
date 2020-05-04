package prockeeper

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// Manager ...
type Manager struct {
	Services []*Service
	list     *tview.List
	logger   *log.Logger
}

var debug = false
var configFile string

var mu sync.Mutex

func init() {
	flag.StringVar(&configFile, "c", "./prockeeper.yml", "config file")
}

func (manager *Manager) updateListItem(index int) {
	s := manager.Services[index]
	title := s.NameWithPid()
	manager.logger.Println("Update list item: ", index, "-", title)
	manager.list.SetItemText(index, title, "")
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

// Run ...
func (manager *Manager) Run() {
	config := ParseConfig(configFile)

	app := tview.NewApplication()

	list := tview.NewList().ShowSecondaryText(false)
	manager.list = list
	list.SetTitle("Services (Press ? to show help)").SetBorder(true)

	debugger := tview.NewTextView().
		SetDynamicColors(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	debugger.SetTitle("debugger").SetBorder(true)

	serviceLog := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	serviceLog.SetBorder(true)

	help := HelpMenu()

	layout := tview.NewFlex().SetDirection(tview.FlexRow)
	pages := tview.NewPages().
		AddPage("app", layout, true, true).
		AddPage("help", modal(help, 40, 20), true, false)
	app.SetRoot(pages, true)

	appContainer := tview.NewFlex().SetDirection(tview.FlexRow)
	appContainer.AddItem(list, 0, 2, true)
	appContainer.AddItem(serviceLog, 0, 6, true)
	layout.AddItem(appContainer, 0, 5, false)

	debuggerContainer := tview.NewFlex()
	debuggerContainer.AddItem(debugger, 0, 1, true)

	appLog, err := os.OpenFile("/tmp/prockeeper.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	CheckError(err)

	defer appLog.Close()

	logger := log.New(io.MultiWriter(debugger, appLog), "", log.LstdFlags)
	manager.logger = logger

	updated := make(chan int)
	go func() {
		for id := range updated {
			manager.updateListItem(id)
		}
	}()

	for i, s := range config.Services {
		service := NewService(i, s.Name, s.Command, s.Dir, updated, logger, serviceLog)
		manager.Services = append(manager.Services, service)
		manager.list.AddItem(service.NameWithPid(), "", 0, nil)
	}
	currentService := manager.Services[0]
	serviceLog.SetTitle(fmt.Sprintf("%s: %s", currentService.Dir, currentService.Command))

	list.SetChangedFunc(func(i int, n string, v string, t rune) {
		currentService.PauseStdout()
		s := manager.Services[i]
		serviceLog.SetTitle(fmt.Sprintf("%s: %s", s.Dir, s.Command))
		serviceLog.SetText(s.History.String())
		s.ResumeStdout()
		currentService = s
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
				debugger.ScrollToEnd()
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
			app.SetFocus(exitMenu)
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

	if err := app.SetFocus(list).Run(); err != nil {
		panic(err)
	}
}
