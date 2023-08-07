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
		Clusters(ctx context.Context, entrypoint string, args map[string]any) ([]entity.Cluster, error)

		Infobases(ctx context.Context, entrypoint string, cluster entity.Cluster, clusterCred entity.Credentials, args map[string]any) ([]entity.Infobase, error)

		Sessions(ctx context.Context, entrypoint string, cluster entity.Cluster, clusterCred entity.Credentials, infobase entity.Infobase, args map[string]any) ([]entity.Session, error)

		Connections(ctx context.Context, entrypoint string, cluster entity.Cluster, clusterCred entity.Credentials, infobase entity.Infobase, args map[string]any) ([]entity.Connection, error)
	}

	// CtrlCache -.
	CtrlCache interface {
		GetClusters(ctx context.Context, entrypoint string) ([]entity.Cluster, error)
		PutClusters(ctx context.Context, entrypoint string, clusters []entity.Cluster) error

		GetInfobases(ctx context.Context, entrypoint string, cluster entity.Cluster) ([]entity.Infobase, error)
		PutInfobases(ctx context.Context, entrypoint string, cluster entity.Cluster, infobases []entity.Infobase) error

		GetSessions(ctx context.Context, entrypoint string, cluster entity.Cluster, infobase entity.Infobase) ([]entity.Session, error)
		PutSessions(ctx context.Context, entrypoint string, cluster entity.Cluster, infobase entity.Infobase, sessions []entity.Session) error

		GetConnections(ctx context.Context, entrypoint string, cluster entity.Cluster, infobase entity.Infobase) ([]entity.Connection, error)
		PutConnections(ctx context.Context, entrypoint string, cluster entity.Cluster, infobase entity.Infobase, connections []entity.Connection) error
	}

	// CtrlPipe -.
	CtrlPipe interface {
		GetClusters(ctx context.Context, entrypoint string) ([]entity.Cluster, error)
		GetInfobases(ctx context.Context, entrypoint string, cluster entity.Cluster, clusterCred entity.Credentials) ([]entity.Infobase, error)
		GetSessions(ctx context.Context, entrypoint string, cluster entity.Cluster, infobase entity.Infobase, clusterCred entity.Credentials) ([]entity.Session, error)
		GetConnections(ctx context.Context, entrypoint string, cluster entity.Cluster, infobase entity.Infobase, clusterCred entity.Credentials) ([]entity.Connection, error)

		DisableSessions(ctx context.Context, entrypoint string, cluster entity.Cluster, infobase entity.Infobase, clusterCred entity.Credentials, infobaseCred entity.Credentials, code string) error
		EnableSessions(ctx context.Context, entrypoint string, cluster entity.Cluster, infobase entity.Infobase, clusterCred entity.Credentials, infobaseCred entity.Credentials, code string) error

		DeleteSession(ctx context.Context, entrypoint string, cluster entity.Cluster, session entity.Session, clusterCred entity.Credentials) error
		DeleteSessions(ctx context.Context, entrypoint string, cluster entity.Cluster, sessions []entity.Session, clusterCred entity.Credentials) error

		DeleteConnection(ctx context.Context, entrypoint string, cluster entity.Cluster, connection entity.Connection, clusterCred entity.Credentials) error
		DeleteConnections(ctx context.Context, entrypoint string, cluster entity.Cluster, connections []entity.Connection, clusterCred entity.Credentials) error
	}

	// CtrlBackup -.
	CtrlBackup interface {
		RunBackup(ctx context.Context, cluster entity.Cluster, infobase entity.Infobase, infobaseCred entity.Credentials, lockCode string, outputPath string) error
	}
)
