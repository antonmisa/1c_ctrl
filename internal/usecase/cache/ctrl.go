package cache

import (
	"context"
	"errors"
	"fmt"

	"github.com/antonmisa/1cctl/internal/entity"
	uc "github.com/antonmisa/1cctl/internal/usecase"
	"github.com/antonmisa/1cctl/pkg/cache"
)

const (
	_keyClusters    string = "%s:clusters"
	_keyInfobases   string = "%s:clusters:%s:ibs"
	_keySessions    string = "%s:clusters:%s:ibs:%s:ses"
	_keyConnections string = "%s:clusters:%s:ibs:%s:conns"
)

var (
	ErrNotFound = errors.New("key not found")
)

// CtrlCache -.
type CtrlCache struct {
	cache cache.Cacher
}

var _ uc.CtrlCache = (*CtrlCache)(nil)

// New -.
func New(c cache.Cacher) *CtrlCache {
	return &CtrlCache{cache: c}
}

// GetClusters -.
func (cc *CtrlCache) GetClusters(ctx context.Context, entrypoint string) ([]entity.Cluster, error) {
	key := fmt.Sprintf(_keyClusters, entrypoint)

	v, ok := cc.cache.Get(key)

	if !ok {
		return []entity.Cluster{}, ErrNotFound
	}

	return v.([]entity.Cluster), nil
}

// PutClusters -.
func (cc *CtrlCache) PutClusters(ctx context.Context, entrypoint string, entities []entity.Cluster) error {
	key := fmt.Sprintf(_keyClusters, entrypoint)

	cc.cache.Set(key, entities)

	return nil
}

// GetInfobases -.
func (cc *CtrlCache) GetInfobases(ctx context.Context, entrypoint string, cluster entity.Cluster) ([]entity.Infobase, error) {
	key := fmt.Sprintf(_keyInfobases, entrypoint, cluster.ID)

	v, ok := cc.cache.Get(key)

	if !ok {
		return []entity.Infobase{}, ErrNotFound
	}

	return v.([]entity.Infobase), nil
}

// PutInfobases -.
func (cc *CtrlCache) PutInfobases(ctx context.Context, entrypoint string, cluster entity.Cluster, entities []entity.Infobase) error {
	key := fmt.Sprintf(_keyInfobases, entrypoint, cluster.ID)

	cc.cache.Set(key, entities)

	return nil
}

// GetSessions -.
func (cc *CtrlCache) GetSessions(ctx context.Context, entrypoint string, cluster entity.Cluster, ib entity.Infobase) ([]entity.Session, error) {
	key := fmt.Sprintf(_keySessions, entrypoint, cluster.ID, ib.ID)

	v, ok := cc.cache.Get(key)

	if !ok {
		return []entity.Session{}, ErrNotFound
	}

	return v.([]entity.Session), nil
}

// PutSessions -.
func (cc *CtrlCache) PutSessions(ctx context.Context, entrypoint string, cluster entity.Cluster, ib entity.Infobase, entities []entity.Session) error {
	key := fmt.Sprintf(_keySessions, entrypoint, cluster.ID, ib.ID)

	cc.cache.Set(key, entities)

	return nil
}

// GetConnections -.
func (cc *CtrlCache) GetConnections(ctx context.Context, entrypoint string, cluster entity.Cluster, ib entity.Infobase) ([]entity.Connection, error) {
	key := fmt.Sprintf(_keyConnections, entrypoint, cluster.ID, ib.ID)

	v, ok := cc.cache.Get(key)

	if !ok {
		return []entity.Connection{}, ErrNotFound
	}

	return v.([]entity.Connection), nil
}

// PutConnection -.
func (cc *CtrlCache) PutConnections(ctx context.Context, entrypoint string, cluster entity.Cluster, ib entity.Infobase, entities []entity.Connection) error {
	key := fmt.Sprintf(_keyConnections, entrypoint, cluster.ID, ib.ID)

	cc.cache.Set(key, entities)

	return nil
}
