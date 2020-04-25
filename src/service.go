package prockeeper

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os/exec"
	"syscall"

	"github.com/rivo/tview"
)

// Service ...
type Service struct {
	Name      string
	Command   string
	Cmd       *exec.Cmd `json:"-"`
	Logger    *log.Logger
	LogView   *tview.TextView
	LogWriter io.Writer
	Updated   chan struct{}
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
	s.LogWriter = tview.ANSIWriter(textView)
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

// Start ...
func (s *Service) Start() error {
	// stopped := make(chan struct{})
	if s.Cmd != nil {
		return errors.New("Error: service is running")
	}

	c := exec.Command("sh", "-c", s.Command)
	c.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	c.Stdout = s.LogWriter
	c.Stderr = s.LogWriter
	s.Cmd = c

	if err := s.Cmd.Start(); err != nil {
		return err
	}

	go func() {
		if err := s.Cmd.Wait(); err != nil {
			s.log(err)
		}
		s.Cmd = nil
		s.log("Job stopped", s.Name)
		s.Updated <- struct{}{}
	}()

	s.log("Started Job")
	s.Updated <- struct{}{}
	return nil
}

// Stop ...
func (s *Service) Stop() error {
	s.log("Stoping job", s.Name)

	if s.Cmd == nil {
		return errors.New("Not running")
	}

	pid, err := syscall.Getpgid(s.Cmd.Process.Pid)
	if err != nil {
		return err
	}
	err = syscall.Kill(-pid, syscall.SIGTERM)

	if err != nil {
		return err
	}
	return nil
}

// Toggle ...
func (s *Service) Toggle() {
	if s.Cmd != nil {
		if err := s.Stop(); err != nil {
			s.log(err)
		}
	} else {
		if err := s.Start(); err != nil {
			s.log(err)
		}
	}
}
