package pipe

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/antonmisa/1cctl/internal/entity"
	uc "github.com/antonmisa/1cctl/internal/usecase"
	"github.com/antonmisa/1cctl/pkg/pipe"
	"golang.org/x/sync/errgroup"
)

const (
	initialSize int = 10

	defaultBlockTime time.Duration = 60

	initialBlockLimitSize int = 50

	formatDate string = "01-02-2006 15:04:05"

	keyInfobase   string = "infobase"
	keyConnection string = "connection"
	keyProcess    string = "process"
	keyHost       string = "host"
	keyName       string = "name"
)

var (
	ErrNotFound          = errors.New("key not found")
	ErrInfobaseIsEmpty   = errors.New("infobase is empty")
	ErrSessionIsEmpty    = errors.New("session is empty")
	ErrConnectionIsEmpty = errors.New("connection is empty")
)

// CtrlPipe -.
type CtrlPipe struct {
	h Helper

	pipe pipe.Piper
}

var _ uc.CtrlPipe = (*CtrlPipe)(nil)

// New -.
func New(p pipe.Piper) *CtrlPipe {
	ctrl := &CtrlPipe{
		h:    Helper{},
		pipe: p,
	}

	return ctrl
}

// GetClusters -.
func (r *CtrlPipe) GetClusters(ctx context.Context, entrypoint string) ([]entity.Cluster, error) {
	args := []string{entrypoint, "cluster", "list"}

	cmd, stdout, err := r.pipe.Run(ctx, args...)
	if err != nil {
		return nil, fmt.Errorf("ctrlpipe - getclusters - r.pipe.Run: %w", err)
	}

	defer cmd.Cancel()
	defer stdout.Close()

	var wg sync.WaitGroup

	datas := make(chan entity.Cluster)
	errs := make(chan error)
	quit := make(chan struct{})

	defer close(quit)
	defer close(errs)
	defer close(datas)

	if err = cmd.Start(); err != nil {
		//break all
		return nil, fmt.Errorf("ctrlpipe - getclusters - cmd.Start: %w", err)
	}

	wg.Add(1)

	go func() {
		defer wg.Done()

		var data, emptyData entity.Cluster

		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()

			key, value, err := r.h.GetKeyValue(line, ':')

			if err != nil && errors.Is(err, ErrNotFound) {
				continue
			}

			if err != nil && !errors.Is(err, ErrNotFound) {
				errs <- fmt.Errorf("ctrlpipe - getclusters - error parsing line: %w", err)

				return
			}

			switch key {
			case "cluster":
				if data != emptyData {
					datas <- data
				}

				data = entity.Cluster{}
				data.ID = value
			case keyHost:
				data.Host = value
			case "port":
				data.Port = value
			case keyName:
				data.Name = value
			}
		}

		if data != emptyData {
			datas <- data
		}

		if err := scanner.Err(); err != nil {
			errs <- fmt.Errorf("ctrlpipe - getclusters - scanner.Err: %w", err)
		}
	}()

	wg.Add(1)

	go func() {
		defer wg.Done()

		if err = cmd.Wait(); err != nil {
			errs <- fmt.Errorf("ctrlpipe - getclusters - cmd.Wait: %w", err)
		}
		quit <- struct{}{}
	}()

	var errg error

	rv := make([]entity.Cluster, 0, initialSize)

	num := 0

	for {
		select {
		case data := <-datas:
			rv = append(rv, data)
			num++
		case errg = <-errs:
			cmd.Cancel()
		case <-quit:
			wg.Wait()

			if errg != nil {
				return nil, errg
			}

			return rv[:num:num], errg
		}
	}
}

// GetInfobases -.
func (r *CtrlPipe) GetInfobases(ctx context.Context, entrypoint string, cluster entity.Cluster, clusterCred entity.Credentials) ([]entity.Infobase, error) {
	args := []string{entrypoint, "infobase", "summary", "list", "--cluster", cluster.ID}

	if clusterCred != (entity.Credentials{}) {
		args = append(args, []string{"--cluster-user", clusterCred.Name, "--cluster-pwd", clusterCred.Pwd}...)
	}

	cmd, stdout, err := r.pipe.Run(ctx, args...)
	if err != nil {
		return nil, fmt.Errorf("ctrlpipe - getinfobases - error opening pipe: %w", err)
	}

	defer cmd.Cancel()
	defer stdout.Close()

	var wg sync.WaitGroup

	datas := make(chan entity.Infobase)
	errs := make(chan error)
	quit := make(chan struct{})

	defer close(quit)
	defer close(errs)
	defer close(datas)

	if err = cmd.Start(); err != nil {
		//break all
		return nil, fmt.Errorf("ctrlpipe - getinfobases - cmd.Start: %w", err)
	}

	wg.Add(1)

	go func() {
		defer wg.Done()

		if err := cmd.Wait(); err != nil {
			errs <- err
		}
		quit <- struct{}{}
	}()

	wg.Add(1)

	go func() {
		defer wg.Done()

		var data, emptyData entity.Infobase

		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()

			key, value, err := r.h.GetKeyValue(line, ':')

			if err != nil && errors.Is(err, ErrNotFound) {
				continue
			}

			if err != nil && !errors.Is(err, ErrNotFound) {
				errs <- fmt.Errorf("ctrlpipe - getinfobases - error parsing line: %w", err)

				return
			}

			switch key {
			case keyInfobase:
				if data != emptyData {
					datas <- data
				}

				data = entity.Infobase{}
				data.ID = value
			case "descr":
				data.Desc = value
			case keyName:
				data.Name = value
			}
		}

		if data != emptyData {
			datas <- data
		}

		if err := scanner.Err(); err != nil {
			errs <- fmt.Errorf("ctrlpipe - getinfobases - scanner.Err: %w", err)
		}
	}()

	var errg error

	rv := make([]entity.Infobase, 0, initialSize)

	num := 0

	for {
		select {
		case data := <-datas:
			rv = append(rv, data)
			num++
		case errg = <-errs:
			cmd.Cancel()
		case <-quit:
			wg.Wait()

			if errg != nil {
				return nil, errg
			}

			return rv[:num:num], errg
		}
	}
}

// GetSessions -.
func (r *CtrlPipe) GetSessions(ctx context.Context, entrypoint string, cluster entity.Cluster, infobase entity.Infobase, clusterCred entity.Credentials) ([]entity.Session, error) {
	args := []string{entrypoint, "session", "list", "--cluster", cluster.ID}

	if clusterCred != (entity.Credentials{}) {
		args = append(args, []string{"--cluster-user", clusterCred.Name, "--cluster-pwd", clusterCred.Pwd}...)
	}

	if infobase != (entity.Infobase{}) {
		args = append(args, []string{"--infobase", infobase.ID}...)
	}

	cmd, stdout, err := r.pipe.Run(ctx, args...)
	if err != nil {
		return nil, fmt.Errorf("ctrlpipe - getsessions - error opening pipe: %w", err)
	}

	defer cmd.Cancel()
	defer stdout.Close()

	var wg sync.WaitGroup

	datas := make(chan entity.Session)
	errs := make(chan error)
	quit := make(chan struct{})

	defer close(quit)
	defer close(errs)
	defer close(datas)

	if err = cmd.Start(); err != nil {
		//break all
		return nil, fmt.Errorf("ctrlpipe - getsessions - cmd.Start: %w", err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := cmd.Wait(); err != nil {
			errs <- err
		}
		quit <- struct{}{}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		var data, emptyData entity.Session

		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()

			key, value, err := r.h.GetKeyValue(line, ':')

			if err != nil && errors.Is(err, ErrNotFound) {
				continue
			}

			if err != nil && !errors.Is(err, ErrNotFound) {
				errs <- fmt.Errorf("ctrlpipe - getsessions - error parsing line: %w", err)

				return
			}

			switch key {
			case "session":
				if data != emptyData {
					datas <- data
				}

				data = entity.Session{}
				data.ID = value
			case keyInfobase:
				data.InfobaseID = value
			case keyConnection:
				data.ConnectionID = value
			case keyProcess:
				data.ProcessID = value
			case "user-name":
				data.UserName = value
			case keyHost:
				data.Host = value
			case "app-id":
				data.AppID = value
			}
		}

		if data != emptyData {
			datas <- data
		}

		if err := scanner.Err(); err != nil {
			errs <- fmt.Errorf("ctrlpipe - getsessions - scanner.Err: %w", err)
		}
	}()

	var errg error

	rv := make([]entity.Session, 0, initialSize)

	num := 0

	for {
		select {
		case data := <-datas:
			rv = append(rv, data)
			num++
		case errg = <-errs:
			cmd.Cancel()
		case <-quit:
			wg.Wait()

			if errg != nil {
				return nil, errg
			}

			return rv[:num:num], errg
		}
	}
}

func (r *CtrlPipe) DisableSessions(ctx context.Context, entrypoint string, cluster entity.Cluster, infobase entity.Infobase, clusterCred entity.Credentials, infobaseCred entity.Credentials, code string) error {
	now := time.Now()

	if infobase == (entity.Infobase{}) {
		return fmt.Errorf("ctrlpipe - disablesessions: %w", ErrInfobaseIsEmpty)
	}

	args := []string{entrypoint, "infobase", "update",
		"--cluster", cluster.ID,
		"--infobase", infobase.ID,
		"--denied-from", now.Format(formatDate),
		"--denied-message", "БАЗА ЗАКРЫТА НА СОЗДАНИЕ РЕЗЕРВНОЙ КОПИИ",
		"--denied-to", now.Add(defaultBlockTime * time.Minute).Format(formatDate),
		"--permission-code", code,
		"--scheduled-jobs-deny", "on",
		"--sessions-deny", "on"}

	if clusterCred != (entity.Credentials{}) {
		args = append(args, []string{"--cluster-user", clusterCred.Name, "--cluster-pwd", clusterCred.Pwd}...)
	}

	if infobaseCred != (entity.Credentials{}) {
		args = append(args, []string{"--infobase-user", infobaseCred.Name, "--infobase-pwd", infobaseCred.Pwd}...)
	}

	cmd, _, err := r.pipe.Run(ctx, args...)

	if err != nil {
		return fmt.Errorf("ctrlpipe - disablesessions - error opening pipe: %w", err)
	}

	defer cmd.Cancel()

	var wg sync.WaitGroup

	errs := make(chan error)
	quit := make(chan struct{})

	defer close(quit)
	defer close(errs)

	if err = cmd.Start(); err != nil {
		//break all
		return fmt.Errorf("ctrlpipe - disablesessions - cmd.Start: %w", err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := cmd.Wait(); err != nil {
			errs <- err
		}
		quit <- struct{}{}
	}()

	var errg error

	for {
		select {
		case errg = <-errs:
			cmd.Cancel()
		case <-quit:
			wg.Wait()

			return errg
		}
	}
}

func (r *CtrlPipe) EnableSessions(ctx context.Context, entrypoint string, cluster entity.Cluster, infobase entity.Infobase, clusterCred entity.Credentials, infobaseCred entity.Credentials, code string) error {
	if infobase == (entity.Infobase{}) {
		return fmt.Errorf("ctrlpipe - enablesessions: %w", ErrInfobaseIsEmpty)
	}

	args := []string{entrypoint, "infobase", "update",
		"--cluster", cluster.ID,
		"--infobase", infobase.ID,
		"--permission-code", code,
		"--scheduled-jobs-deny", "off",
		"--sessions-deny", "off"}

	if clusterCred != (entity.Credentials{}) {
		args = append(args, []string{"--cluster-user", clusterCred.Name, "--cluster-pwd", clusterCred.Pwd}...)
	}

	if infobaseCred != (entity.Credentials{}) {
		args = append(args, []string{"--infobase-user", infobaseCred.Name, "--infobase-pwd", infobaseCred.Pwd}...)
	}

	cmd, _, err := r.pipe.Run(ctx, args...)

	if err != nil {
		return fmt.Errorf("ctrlpipe - enablesessions - error opening pipe: %w", err)
	}

	defer cmd.Cancel()

	var wg sync.WaitGroup

	errs := make(chan error)
	quit := make(chan struct{})

	defer close(quit)
	defer close(errs)

	if err = cmd.Start(); err != nil {
		//break all
		return fmt.Errorf("ctrlpipe - enablesessions - cmd.Start: %w", err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := cmd.Wait(); err != nil {
			errs <- err
		}
		quit <- struct{}{}
	}()

	var errg error

	for {
		select {
		case errg = <-errs:
			cmd.Cancel()
		case <-quit:
			wg.Wait()

			return errg
		}
	}
}

func (r *CtrlPipe) DeleteSession(ctx context.Context, entrypoint string, cluster entity.Cluster, session entity.Session, clusterCred entity.Credentials) error {
	if session == (entity.Session{}) {
		return fmt.Errorf("ctrlpipe - deletesession: %w", ErrSessionIsEmpty)
	}

	args := []string{entrypoint, "session", "terminate",
		"--cluster", cluster.ID,
		"--session", session.ID}

	if clusterCred != (entity.Credentials{}) {
		args = append(args, []string{"--cluster-user", clusterCred.Name, "--cluster-pwd", clusterCred.Pwd}...)
	}

	cmd, _, err := r.pipe.Run(ctx, args...)

	if err != nil {
		return fmt.Errorf("ctrlpipe - deletesession - error opening pipe: %w", err)
	}

	defer cmd.Cancel()

	var wg sync.WaitGroup

	errs := make(chan error)
	quit := make(chan struct{})

	defer close(quit)
	defer close(errs)

	if err = cmd.Start(); err != nil {
		//break all
		return fmt.Errorf("ctrlpipe - deletesession - cmd.Start: %w", err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := cmd.Wait(); err != nil {
			errs <- err
		}
		quit <- struct{}{}
	}()

	var errg error

	for {
		select {
		case errg = <-errs:
			cmd.Cancel()
		case <-quit:
			wg.Wait()

			return errg
		}
	}
}

func (r *CtrlPipe) DeleteSessions(ctx context.Context, entrypoint string, cluster entity.Cluster, sessions []entity.Session, clusterCred entity.Credentials) error {
	g, ctx := errgroup.WithContext(ctx)

	g.SetLimit(initialBlockLimitSize)

	for i := range sessions {
		i := i

		g.Go(func() error {
			return r.DeleteSession(ctx, entrypoint, cluster, sessions[i], clusterCred)
		})
	}

	err := g.Wait()
	return err
}

func (r *CtrlPipe) GetConnections(ctx context.Context, entrypoint string, cluster entity.Cluster, infobase entity.Infobase, clusterCred entity.Credentials) ([]entity.Connection, error) {
	args := []string{entrypoint, "connection", "list", "--cluster", cluster.ID}

	if clusterCred != (entity.Credentials{}) {
		args = append(args, []string{"--cluster-user", clusterCred.Name, "--cluster-pwd", clusterCred.Pwd}...)
	}

	if infobase != (entity.Infobase{}) {
		args = append(args, []string{"--infobase", infobase.ID}...)
	}

	cmd, stdout, err := r.pipe.Run(ctx, args...)
	if err != nil {
		return nil, fmt.Errorf("ctrlpipe - getconnections - error opening pipe: %w", err)
	}

	defer cmd.Cancel()
	defer stdout.Close()

	var wg sync.WaitGroup

	datas := make(chan entity.Connection)
	errs := make(chan error)
	quit := make(chan struct{})

	defer close(quit)
	defer close(errs)
	defer close(datas)

	if err = cmd.Start(); err != nil {
		//break all
		return nil, fmt.Errorf("ctrlpipe - getconnections - cmd.Start: %w", err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := cmd.Wait(); err != nil {
			errs <- err
		}
		quit <- struct{}{}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		var data, emptyData entity.Connection

		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()

			key, value, err := r.h.GetKeyValue(line, ':')

			if err != nil && errors.Is(err, ErrNotFound) {
				continue
			}

			if err != nil && !errors.Is(err, ErrNotFound) {
				errs <- fmt.Errorf("ctrlpipe - getconnections - error parsing line: %w", err)

				return
			}

			switch key {
			case keyConnection:
				if data != emptyData {
					datas <- data
				}

				data = entity.Connection{}
				data.ID = value
			case keyHost:
				data.Host = value
			case keyProcess:
				data.ProcessID = value
			case keyInfobase:
				data.InfobaseID = value
			case "application":
				data.AppID = value
			}
		}

		if data != emptyData {
			datas <- data
		}

		if err := scanner.Err(); err != nil {
			errs <- fmt.Errorf("ctrlpipe - getconnections - scanner.Err: %w", err)
		}
	}()

	var errg error

	rv := make([]entity.Connection, 0, initialSize)

	num := 0

	for {
		select {
		case data := <-datas:
			rv = append(rv, data)
			num++
		case errg = <-errs:
			cmd.Cancel()
		case <-quit:
			wg.Wait()

			if errg != nil {
				return nil, errg
			}

			return rv[:num:num], errg
		}
	}
}

func (r *CtrlPipe) DeleteConnection(ctx context.Context, entrypoint string, cluster entity.Cluster, connection entity.Connection, clusterCred entity.Credentials) error {
	if connection == (entity.Connection{}) {
		return fmt.Errorf("ctrlpipe - deleteconnection: %w", ErrConnectionIsEmpty)
	}

	args := []string{entrypoint, "connection", "disconnect",
		"--cluster", cluster.ID,
		"--conection", connection.ID}

	if clusterCred != (entity.Credentials{}) {
		args = append(args, []string{"--cluster-user", clusterCred.Name, "--cluster-pwd", clusterCred.Pwd}...)
	}

	cmd, _, err := r.pipe.Run(ctx, args...)

	if err != nil {
		return fmt.Errorf("ctrlpipe - deleteconnection - error opening pipe: %w", err)
	}

	defer cmd.Cancel()

	var wg sync.WaitGroup

	errs := make(chan error)
	quit := make(chan struct{})

	defer close(quit)
	defer close(errs)

	if err = cmd.Start(); err != nil {
		//break all
		return fmt.Errorf("ctrlpipe - deleteconnection - cmd.Start: %w", err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := cmd.Wait(); err != nil {
			errs <- err
		}
		quit <- struct{}{}
	}()

	var errg error

	for {
		select {
		case errg = <-errs:
			cmd.Cancel()
		case <-quit:
			wg.Wait()

			return errg
		}
	}
}

func (r *CtrlPipe) DeleteConnections(ctx context.Context, entrypoint string, cluster entity.Cluster, connections []entity.Connection, clusterCred entity.Credentials) error {
	g, ctx := errgroup.WithContext(ctx)

	g.SetLimit(initialBlockLimitSize)

	for i := range connections {
		i := i

		g.Go(func() error {
			return r.DeleteConnection(ctx, entrypoint, cluster, connections[i], clusterCred)
		})
	}

	err := g.Wait()
	return err
}
