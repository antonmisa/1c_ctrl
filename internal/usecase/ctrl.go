package usecase

import (
	"context"
	"errors"
	"fmt"
	"go.opentelemetry.io/otel/trace"
	"reflect"

	"github.com/antonmisa/1cctl/internal/entity"
	"github.com/antonmisa/1cctl/internal/usecase/cache"
	"github.com/antonmisa/1cctl/internal/usecase/common"
)

// CtrlUseCase -.
type CtrlUseCase struct {
	cache  CtrlCache
	pipe   CtrlPipe
	backup CtrlBackup
}

var _ Ctrl = (*CtrlUseCase)(nil)

// New -.
func New(c CtrlCache, p CtrlPipe, b CtrlBackup) *CtrlUseCase {
	return &CtrlUseCase{
		cache:  c,
		pipe:   p,
		backup: b,
	}
}

// Clusters - getting clusters list in cache.
func (c *CtrlUseCase) Clusters(ctx context.Context, entrypoint string, args map[string]any) ([]entity.Cluster, error) {
	span := trace.SpanFromContext(ctx)

	if v, ok := args[common.UseCache]; ok && v.(bool) {
		span.AddEvent("GetClusters from cache")

		clusters, err := c.cache.GetClusters(ctx, entrypoint)

		if err != nil && !errors.Is(err, cache.ErrNotFound) {
			return nil, fmt.Errorf("CtrlUseCase - Clusters - c.cache.GetClusters: %w", err)
		} else if !reflect.DeepEqual(clusters, []entity.Cluster{}) {
			return clusters, nil
		}
	}

	span.AddEvent("GetClusters from 1C")

	clusters, err := c.pipe.GetClusters(ctx, entrypoint)
	if err != nil {
		return nil, fmt.Errorf("CtrlUseCase - Clusters - c.pipe.GetClusters: %w", err)
	}

	span.AddEvent("PutClusters to cache")

	err = c.cache.PutClusters(ctx, entrypoint, clusters)
	if err != nil {
		return nil, fmt.Errorf("CtrlUseCase - Clusters - c.pipe.PutClusters: %w", err)
	}

	return clusters, nil
}

// Infobases - getting infobases list for cluster.
func (c *CtrlUseCase) Infobases(ctx context.Context, entrypoint string, cluster entity.Cluster, clusterCred entity.Credentials, args map[string]any) ([]entity.Infobase, error) {
	span := trace.SpanFromContext(ctx)

	if v, ok := args[common.UseCache]; ok && v.(bool) {
		span.AddEvent("GetInfobases from cache")

		infobases, err := c.cache.GetInfobases(ctx, entrypoint, cluster)

		if err != nil && !errors.Is(err, cache.ErrNotFound) {
			return nil, fmt.Errorf("CtrlUseCase - Infobases - c.cache.GetInfobases: %w", err)
		} else if !reflect.DeepEqual(infobases, []entity.Infobase{}) {
			return infobases, nil
		}
	}

	span.AddEvent("GetInfobases from 1C")

	infobases, err := c.pipe.GetInfobases(ctx, entrypoint, cluster, clusterCred)
	if err != nil {
		return nil, fmt.Errorf("CtrlUseCase - Infobases - c.pipe.GetInfobases: %w", err)
	}

	span.AddEvent("PutInfobases to cache")

	err = c.cache.PutInfobases(ctx, entrypoint, cluster, infobases)
	if err != nil {
		return nil, fmt.Errorf("CtrlUseCase - Infobases - c.pipe.PutInfobases: %w", err)
	}

	return infobases, nil
}

// Sessions - getting sessions list for cluster.
func (c *CtrlUseCase) Sessions(ctx context.Context, entrypoint string, cluster entity.Cluster, clusterCred entity.Credentials, infobase entity.Infobase, args map[string]any) ([]entity.Session, error) {
	span := trace.SpanFromContext(ctx)

	if v, ok := args[common.UseCache]; ok && v.(bool) {
		span.AddEvent("GetSessions from cache")

		sessions, err := c.cache.GetSessions(ctx, entrypoint, cluster, infobase)

		if err != nil && !errors.Is(err, cache.ErrNotFound) {
			return nil, fmt.Errorf("CtrlUseCase - Sessions - c.cache.GetSessions: %w", err)
		} else if !reflect.DeepEqual(sessions, []entity.Session{}) {
			return sessions, nil
		}
	}

	span.AddEvent("GetSessions from 1C")

	sessions, err := c.pipe.GetSessions(ctx, entrypoint, cluster, infobase, clusterCred)
	if err != nil {
		return nil, fmt.Errorf("CtrlUseCase - Sessions - c.pipe.GetSessions: %w", err)
	}

	span.AddEvent("PutSessions to cache")

	err = c.cache.PutSessions(ctx, entrypoint, cluster, infobase, sessions)
	if err != nil {
		return nil, fmt.Errorf("CtrlUseCase - Sessions - c.pipe.PutSessions: %w", err)
	}

	return sessions, nil
}

// Connections - getting connections list for cluster.
func (c *CtrlUseCase) Connections(ctx context.Context, entrypoint string, cluster entity.Cluster, clusterCred entity.Credentials, infobase entity.Infobase, args map[string]any) ([]entity.Connection, error) {
	span := trace.SpanFromContext(ctx)

	if v, ok := args[common.UseCache]; ok && v.(bool) {
		span.AddEvent("GetConnections from cache")

		connections, err := c.cache.GetConnections(ctx, entrypoint, cluster, infobase)

		if err != nil && !errors.Is(err, cache.ErrNotFound) {
			return nil, fmt.Errorf("CtrlUseCase - Connections - c.cache.GetConnections: %w", err)
		} else if !reflect.DeepEqual(connections, []entity.Connection{}) {
			return connections, nil
		}
	}

	span.AddEvent("GetConnections from 1C")

	connections, err := c.pipe.GetConnections(ctx, entrypoint, cluster, infobase, clusterCred)
	if err != nil {
		return nil, fmt.Errorf("CtrlUseCase - Connections - uc.pipe.GetConnections: %w", err)
	}

	span.AddEvent("PutConnections to cache")

	err = c.cache.PutConnections(ctx, entrypoint, cluster, infobase, connections)
	if err != nil {
		return nil, fmt.Errorf("CtrlUseCase - Connections - c.pipe.PutConnections: %w", err)
	}

	return connections, nil
}
