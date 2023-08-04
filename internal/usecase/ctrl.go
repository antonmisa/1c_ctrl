package usecase

import (
	"context"
	"fmt"

	"github.com/antonmisa/1cctl/internal/entity"
)

// TranslationUseCase -.
type CtrlUseCase struct {
	cache  CtrlCache
	pipe   CtrlPipe
	backup CtrlBackup
}

// New -.
func New(c CtrlCache, p CtrlPipe, b CtrlBackup) *CtrlUseCase {
	return &CtrlUseCase{
		cache:  c,
		pipe:   p,
		backup: b,
	}
}

// Clusters - getting clusters list.
func (uc *CtrlUseCase) Clusters(ctx context.Context) ([]entity.Cluster, error) {
	clusters, err := uc.cache.GetClusters(ctx)
	if err != nil {
		return nil, fmt.Errorf("CtrlUseCase - Clusters - uc.cache.GetClusters: %w", err)
	} else if clusters != nil {
		return clusters, nil
	}

	clusters, err = uc.pipe.GetClusters(ctx)
	if err != nil {
		return nil, fmt.Errorf("CtrlUseCase - Clusters - uc.pipe.GetClusters: %w", err)
	}

	err = uc.cache.PutClusters(ctx, clusters)
	if err != nil {
		return nil, fmt.Errorf("CtrlUseCase - Clusters - uc.pipe.PutClusters: %w", err)
	}

	return clusters, nil
}

// Infobases - getting infobases list for cluster.
func (uc *CtrlUseCase) Infobases(ctx context.Context, cluster entity.Cluster) ([]entity.Infobase, error) {
	infobases, err := uc.cache.GetInfobases(ctx, cluster)
	if err != nil {
		return nil, fmt.Errorf("CtrlUseCase - Infobases - uc.cache.GetInfobases: %w", err)
	} else if infobases != nil {
		return infobases, nil
	}

	infobases, err = uc.pipe.GetInfobases(ctx, cluster)
	if err != nil {
		return nil, fmt.Errorf("CtrlUseCase - Infobases - uc.pipe.GetInfobases: %w", err)
	}

	err = uc.cache.PutInfobases(ctx, cluster, infobases)
	if err != nil {
		return nil, fmt.Errorf("CtrlUseCase - Infobases - uc.pipe.PutInfobases: %w", err)
	}

	return infobases, nil
}

// Sessions - getting sessions list for cluster.
func (uc *CtrlUseCase) Sessions(ctx context.Context, cluster entity.Cluster, ib entity.Infobase) ([]entity.Session, error) {
	sessions, err := uc.cache.GetSessions(ctx, cluster, ib)
	if err != nil {
		return nil, fmt.Errorf("CtrlUseCase - Sessions - uc.cache.GetSessions: %w", err)
	} else if sessions != nil {
		return sessions, nil
	}

	sessions, err = uc.pipe.GetSessions(ctx, cluster, ib)
	if err != nil {
		return nil, fmt.Errorf("CtrlUseCase - Sessions - uc.pipe.GetSessions: %w", err)
	}

	err = uc.cache.PutSessions(ctx, cluster, ib, sessions)
	if err != nil {
		return nil, fmt.Errorf("CtrlUseCase - Sessions - uc.pipe.PutSessions: %w", err)
	}

	return sessions, nil
}

// Connections - getting connections list for cluster.
func (uc *CtrlUseCase) Connections(ctx context.Context, cluster entity.Cluster, ib entity.Infobase) ([]entity.Connection, error) {
	connections, err := uc.cache.GetConnections(ctx, cluster, ib)
	if err != nil {
		return nil, fmt.Errorf("CtrlUseCase - Connections - uc.cache.GetConnections: %w", err)
	} else if connections != nil {
		return connections, nil
	}

	connections, err = uc.pipe.GetConnections(ctx, cluster, ib)
	if err != nil {
		return nil, fmt.Errorf("CtrlUseCase - Connections - uc.pipe.GetConnections: %w", err)
	}

	err = uc.cache.PutConnections(ctx, cluster, ib, connections)
	if err != nil {
		return nil, fmt.Errorf("CtrlUseCase - Connections - uc.pipe.PutConnections: %w", err)
	}

	return connections, nil
}
