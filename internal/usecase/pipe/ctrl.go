package pipe

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/antonmisa/1cctl/internal/entity"
	"github.com/antonmisa/1cctl/pkg/pipe"
)

const (
	initialSize int = 10
	notFound        = Error("key not found")
)

// CtrlPipe -.
type CtrlPipe struct {
	*pipe.Pipe
}

// New -.
func New(p *pipe.Pipe) *CtrlPipe {
	ctrl := &CtrlPipe{
		p,
	}

	return ctrl
}

// GetClusters -.
func (r *CtrlPipe) GetClusters(ctx context.Context) ([]entity.Cluster, error) {
	cmd, stdout, err := r.Pipe.Run("cluster", "list")
	if err != nil {
		return nil, fmt.Errorf("ctrlpipe - getclusters - error opening pipe: %w", err)
	}

	buf := bufio.NewReader(stdout)
	num := 0

	cmd.Start()
	defer cmd.Cancel()

	data := make([]entity.Cluster, initialSize)

	var cluster, emptyCluster entity.Cluster

	doWork := true
	for doWork {
		line, _, err := buf.ReadLine()
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("ctrlpipe - getclusters - error reading pipe: %w", err)
		}

		if err == io.EOF {
			doWork = false
		}

		sline := string(line)
		key, value, err := r.GetKeyValue(sline, ':')

		if err != nil && errors.Is(err, notFound) {
			continue
		}

		if err != nil && !errors.Is(err, notFound) {
			return nil, fmt.Errorf("ctrlpipe - getclusters - error parsing line: %w", err)
		}

		switch key {
		case "cluster":
			if cluster != emptyCluster {
				data = append(data, cluster)
				num += 1
			}

			cluster = entity.Cluster{}
			cluster.ID = value
		case "host":
			cluster.Host = value
		case "port":
			cluster.Port = value
		case "name":
			cluster.Name = value
		}
	}

	if cluster != emptyCluster {
		data = append(data, cluster)
		num += 1
	}

	return data[:num:num], nil
}

// GetInfobases -.
func (r *CtrlPipe) GetInfobases(ctx context.Context, cluster entity.Cluster) ([]entity.Infobase, error) {
	cmd, stdout, err := r.Pipe.Run("infobase", fmt.Sprintf("--cluster=%s", cluster.ID), "summary", "list")
	if err != nil {
		return nil, fmt.Errorf("ctrlpipe - getinfobases - error opening pipe: %w", err)
	}

	buf := bufio.NewReader(stdout)
	num := 0

	cmd.Start()
	defer cmd.Cancel()

	data := make([]entity.Infobase, initialSize)

	var ib, emptyIB entity.Infobase

	doWork := true
	for doWork {
		line, _, err := buf.ReadLine()
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("ctrlpipe - getinfobases - error reading pipe: %w", err)
		}

		if err == io.EOF {
			doWork = false
		}

		sline := string(line)
		key, value, err := r.GetKeyValue(sline, ':')

		if err != nil && errors.Is(err, notFound) {
			continue
		}

		if err != nil && !errors.Is(err, notFound) {
			return nil, fmt.Errorf("ctrlpipe - getinfobases - error parsing line: %w", err)
		}

		switch key {
		case "infobase":
			if ib != emptyIB {
				data = append(data, ib)
				num += 1
			}

			ib = entity.Infobase{ClusterID: cluster.ID}
			ib.ID = value
		case "descr":
			ib.Desc = value
		case "name":
			ib.Name = value
		}
	}

	if ib != emptyIB {
		data = append(data, ib)
		num += 1
	}

	return data[:num:num], nil
}

// GetSessions -.
func (r *CtrlPipe) GetSessions(ctx context.Context, cluster entity.Cluster, infobase entity.Infobase) ([]entity.Session, error) {
	filterByInfobase := ""
	if infobase != (entity.Infobase{}) {
		filterByInfobase = fmt.Sprintf("--infobase=%s", infobase.ID)
	}

	cmd, stdout, err := r.Pipe.Run("session", fmt.Sprintf("--cluster=%s", cluster.ID), "list", filterByInfobase)
	if err != nil {
		return nil, fmt.Errorf("ctrlpipe - getsessions - error opening pipe: %w", err)
	}

	buf := bufio.NewReader(stdout)
	num := 0

	cmd.Start()
	defer cmd.Cancel()

	data := make([]entity.Session, initialSize)

	var ses, emptySes entity.Session

	doWork := true
	for doWork {
		line, _, err := buf.ReadLine()
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("ctrlpipe - getsessions - error reading pipe: %w", err)
		}

		if err == io.EOF {
			doWork = false
		}

		sline := string(line)
		key, value, err := r.GetKeyValue(sline, ':')

		if err != nil && errors.Is(err, notFound) {
			continue
		}

		if err != nil && !errors.Is(err, notFound) {
			return nil, fmt.Errorf("ctrlpipe - getsessions - error parsing line: %w", err)
		}

		switch key {
		case "session":
			if ses != emptySes {
				data = append(data, ses)
				num += 1
			}

			ses = entity.Session{}
			ses.ID = value
		case "infobase":
			ses.InfobaseID = value
		case "connection":
			ses.ConnectionID = value
		case "process":
			ses.ProcessID = value
		case "user-name":
			ses.UserName = value
		case "host":
			ses.Host = value
		case "app-id":
			ses.AppID = value
		}
	}

	if ses != emptySes {
		data = append(data, ses)
		num += 1
	}

	return data[:num:num], nil
}

func (r *CtrlPipe) DisableSessions(ctx context.Context, infobase entity.Infobase) error {
	cmd, _, err := r.Pipe.Run("session", fmt.Sprintf("--cluster=%s", infobase.ClusterID),
		"infobase", "update", fmt.Sprintf("--infobase=%s", infobase.ID),
		fmt.Sprintf("--denied-from=%s", infobase.ID),
		fmt.Sprintf("--denied-message=%s", "БАЗА ЗАКРЫТА НА СОЗДАНИЕ РЕЗЕРВНОЙ КОПИИ"),
		fmt.Sprintf("--denied-to=%s", infobase.ID),
		"--permission-code=12345",
		"--scheduled-jobs-deny=on",
		"--session-deny=on")
	if err != nil {
		return fmt.Errorf("ctrlpipe - DisableSessions - error opening pipe: %w", err)
	}

	cmd.Start()
	defer cmd.Cancel()

	return nil
}

func (r *CtrlPipe) EnableSessions(ctx context.Context, infobase entity.Infobase) error {
	cmd, _, err := r.Pipe.Run("session", fmt.Sprintf("--cluster=%s", infobase.ClusterID),
		"infobase", "update", fmt.Sprintf("--infobase=%s", infobase.ID),
		"--permission-code=12345",
		"--scheduled-jobs-deny=off",
		"--session-deny=off")
	if err != nil {
		return fmt.Errorf("ctrlpipe - EnableSessions - error opening pipe: %w", err)
	}

	cmd.Start()
	defer cmd.Cancel()

	return nil
}

func (r *CtrlPipe) DeleteSession(ctx context.Context, cluster entity.Cluster, session entity.Session) error {
	cmd, _, err := r.Pipe.Run("session", fmt.Sprintf("--cluster=%s", cluster.ID),
		"terminate", fmt.Sprintf("--session=%s", session.ID),
		"--permission-code=12345",
		"--scheduled-jobs-deny=off",
		"--session-deny=off")
	if err != nil {
		return fmt.Errorf("ctrlpipe - DeleteSession - error opening pipe: %w", err)
	}

	cmd.Start()
	defer cmd.Cancel()

	return nil
}

func (r *CtrlPipe) DeleteSessions(ctx context.Context, cluster entity.Cluster, sessions []entity.Session) error {
	//TODO: concurent run
	var errs error

	for i := range sessions {
		err := r.DeleteSession(ctx, cluster, sessions[i])
		if err != nil {
			errs = fmt.Errorf("%w - %w", errs, err)
		}
	}

	return errs
}

func (r *CtrlPipe) GetConnections(ctx context.Context, cluster entity.Cluster, infobase entity.Infobase) ([]entity.Connection, error) {
	return []entity.Connection{}, nil
}

func (r *CtrlPipe) DeleteConnection(ctx context.Context, cluster entity.Cluster, connection entity.Connection) error {
	cmd, _, err := r.Pipe.Run("connection", fmt.Sprintf("--cluster=%s", cluster.ID),
		"disconnect", fmt.Sprintf("--connection=%s", connection.ID))
	if err != nil {
		return fmt.Errorf("ctrlpipe - DeleteConnection - error opening pipe: %w", err)
	}

	cmd.Start()
	defer cmd.Cancel()

	return nil
}

func (r *CtrlPipe) DeleteConnections(ctx context.Context, cluster entity.Cluster, connections []entity.Connection) error {
	//TODO: concurent run
	var errs error

	for i := range connections {
		err := r.DeleteConnection(ctx, cluster, connections[i])
		if err != nil {
			errs = fmt.Errorf("%w - %w", errs, err)
		}
	}

	return errs
}

func (r *CtrlPipe) GetKeyValue(line string, delimeter rune) (string, string, error) {
	if pos := strings.IndexRune(line, delimeter); pos != -1 {
		return line[:pos], line[pos+1:], nil
	}

	return "", "", notFound
}
