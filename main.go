package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"syscall"

	"github.com/rivo/tview"
	"gopkg.in/yaml.v2"
)

// Service ...
type Service struct {
	Name    string
	Command string
	Cmd     *exec.Cmd `json:"-"`
	Logger  *log.Logger
	LogView *tview.TextView
}

// Manager ...
type Manager struct {
	Services []*Service
}

// Prepare ...
func (s *Service) Prepare(app *tview.Application, logger *log.Logger) {
	s.Logger = logger

	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetChangedFunc(func() {
			app.Draw()
		})

	textView.SetBorder(true).SetTitle(s.Command)
	s.LogView = textView
}

func (s *Service) log(v ...interface{}) {
	s.Logger.Println(v...)
}

func (s *Service) pid() int {
	if s.Cmd != nil && s.Cmd.Process != nil {
		return s.Cmd.Process.Pid
	}
	return 0
}

// NameWithPid ...
func (s *Service) NameWithPid() string {
	pid := s.pid()
	if pid == 0 {
		return fmt.Sprintf("[      ] %s", s.Name)
	}
	return fmt.Sprintf("[%6d] %s", pid, s.Name)
}

func (s *Service) toggle() {
	if s.Cmd == nil {
		c := exec.Command("sh", "-c", s.Command)
		c.Stdout = s.LogView
		s.Cmd = c
		if err := c.Start(); err != nil {
			s.log(err)
		}
		s.log("Started job", s.NameWithPid())
		go func() {
			if err := c.Wait(); err != nil {
				s.log(err)
			}
			s.log("Job stopped", s.Name)
		}()
	} else {
		s.log("Stoping job", s.Name)
		s.Cmd.Process.Signal(syscall.SIGTERM)
		s.Cmd = nil
	}
}

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

func main() {
	manager := &Manager{}

	config, err := ioutil.ReadFile("./smgo.yml")
	checkError(err)

	err = yaml.Unmarshal(config, manager)
	checkError(err)

	app := tview.NewApplication()

	layout := tview.NewFlex().SetDirection(tview.FlexRow)
	app.SetRoot(layout, true)

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

	list := tview.NewList().ShowSecondaryText(false)
	list.SetBorder(true)
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
		manager.Services[i].toggle()
		list.Clear()
		for _, s := range manager.Services {
			list.AddItem(s.NameWithPid(), "", 0, nil)
		}
		list.SetCurrentItem(i)
	})

	if err := app.SetFocus(list).Run(); err != nil {
		panic(err)
	}
}
