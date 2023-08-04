package pipe

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
)

var (
	ErrNoFile = errors.New("file does not exist")
)

//go:generate go run github.com/vektra/mockery/v2@v2.32.0 --all

type Piper interface {
	Run(ctx context.Context, arg ...string) (Commander, io.ReadCloser, error)
}

type Pipe struct {
	pathToRAC string
}

var _ Piper = (*Pipe)(nil)

func New(path string) (*Pipe, error) {
	_, err := os.Stat(path)
	if _, ok := err.(*fs.PathError); ok || os.IsNotExist(err) {
		return nil, fmt.Errorf("%w: %w", ErrNoFile, err)
	}

	return &Pipe{
		pathToRAC: path,
	}, nil
}

func (p Pipe) Run(ctx context.Context, arg ...string) (Commander, io.ReadCloser, error) {
	cmd := exec.CommandContext(ctx, p.pathToRAC, arg...) //nolint:gosec // it is normal

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}

	return &Command{cmd}, stdout, nil
}
