package pipe

import "os/exec"

//go:generate go run github.com/vektra/mockery/v2@v2.32.0 --all

type Commander interface {
	Start() error
	Wait() error
	Cancel() error
}

type Command struct {
	*exec.Cmd
}

func (c *Command) Cancel() error {
	return c.Cmd.Cancel()
}
