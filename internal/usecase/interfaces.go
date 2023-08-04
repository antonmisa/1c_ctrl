// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"

	"github.com/antonmisa/1cctl/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

type (
	// Ctrl -.
	Ctrl interface {
		Clusters(context.Context) ([]entity.Cluster, error)
		Infobases(context.Context, entity.Cluster) ([]entity.Infobase, error)
		Sessions(context.Context, entity.Cluster, entity.Infobase) ([]entity.Session, error)
	}

	// CtrlCache -.
	CtrlCache interface {
		GetClusters(context.Context) ([]entity.Cluster, error)
		PutClusters(context.Context, []entity.Cluster) error

		GetInfobases(context.Context, entity.Cluster) ([]entity.Infobase, error)
		PutInfobases(context.Context, entity.Cluster, []entity.Infobase) error

		GetSessions(context.Context, entity.Cluster, entity.Infobase) ([]entity.Session, error)
		PutSessions(context.Context, entity.Cluster, entity.Infobase, []entity.Session) error

		GetConnections(context.Context, entity.Cluster, entity.Infobase) ([]entity.Connection, error)
		PutConnections(context.Context, entity.Cluster, entity.Infobase, []entity.Connection) error
	}

	// CtrlPipe -.
	CtrlPipe interface {
		GetClusters(context.Context) ([]entity.Cluster, error)
		GetInfobases(context.Context, entity.Cluster) ([]entity.Infobase, error)
		GetSessions(context.Context, entity.Cluster, entity.Infobase) ([]entity.Session, error)
		GetConnections(context.Context, entity.Cluster, entity.Infobase) ([]entity.Connection, error)

		DisableSessions(context.Context, entity.Infobase) error
		EnableSessions(context.Context, entity.Infobase) error

		DeleteSession(context.Context, entity.Cluster, entity.Session) error
		DeleteSessions(context.Context, entity.Cluster, []entity.Session) error

		DeleteConnection(context.Context, entity.Cluster, entity.Connection) error
		DeleteConnections(context.Context, entity.Cluster, []entity.Connection) error
	}

	// CtrlBackup -.
	CtrlBackup interface {
		RunBackup(context.Context, entity.Cluster, entity.Infobase, string, string, string) error
	}
)
