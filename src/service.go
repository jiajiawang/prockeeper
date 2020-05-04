package prockeeper

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os/exec"
	"syscall"

	"github.com/rivo/tview"
)

// Service ...
type Service struct {
	Id      int
	Name    string
	Command string
	Dir     string
	Cmd     *exec.Cmd `json:"-"`
	Logger  *log.Logger
	History *bytes.Buffer
	Updated chan int

	stdout       *PausableWriter
	cmdLogWriter io.Writer
}

// NewService ...
func NewService(
	id int,
	name, command, dir string,
	updated chan int,
	logger *log.Logger,
	out io.Writer,
) *Service {
	s := &Service{
		Id:      id,
		Name:    name,
		Command: command,
		Dir:     dir,
		Updated: updated,
	}
	s.Logger = logger
	s.History = new(bytes.Buffer)
	s.stdout = NewPausableWriter(out)
	s.stdout.Pause()
	s.cmdLogWriter = tview.ANSIWriter(io.MultiWriter(s.History, s.stdout))
	return s
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

// PauseStdout ...
func (s *Service) PauseStdout() {
	s.stdout.Pause()
}

// ResumeStdout ...
func (s *Service) ResumeStdout() {
	s.stdout.Resume()
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
		return fmt.Errorf("Error: %s is running", s.Name)
	}

	c := exec.Command("sh", "-c", s.Command)
	c.Dir = s.Dir
	s.Cmd = c
	c.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	c.Stdout = s.cmdLogWriter
	c.Stderr = s.cmdLogWriter

	if err := s.Cmd.Start(); err != nil {
		return err
	}

	go func() {
		if err := s.Cmd.Wait(); err != nil {
			s.log(err)
		}
		s.Cmd = nil
		s.log("Stopped service -", s.Name)
		s.Updated <- s.Id
	}()

	s.log("Started service -", s.Name)
	s.Updated <- s.Id
	return nil
}

// Stop ...
func (s *Service) Stop() error {
	s.log("Stopping service -", s.Name)

	if s.Cmd == nil {
		return fmt.Errorf("Error: %s is not running", s.Name)
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
