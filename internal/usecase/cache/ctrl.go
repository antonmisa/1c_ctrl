package cache

import (
	"context"
	"fmt"

	"github.com/antonmisa/1cctl/internal/entity"
	"github.com/antonmisa/1cctl/pkg/cache"
)

const (
	_keyClusters  string = "clusters"
	_keyInfobases string = "%s:ibs"
	_keySessions  string = "%s:%s:ses"
)

// CtrlCache -.
type CtrlCache struct {
	*cache.Cache
}

// New -.
func New(c *cache.Cache) *CtrlCache {
	return &CtrlCache{c}
}

// GetClusters -.
func (cc *CtrlCache) GetClusters(ctx context.Context) ([]entity.Cluster, error) {
	var entities []entity.Cluster
	if err := cc.Cache.Get(ctx, _keyClusters, &entities); err != nil {
		return nil, fmt.Errorf("CtrlCache - GetClusters - cc.Cache.Get: : %w", err)
	}

	return entities, nil
}

// PutClusters -.
func (cc *CtrlCache) PutClusters(ctx context.Context, entities []entity.Cluster) error {
	if len(entities) == 0 {
		return nil
	}

	return cc.Store(ctx, _keyClusters, entities)
}

// GetInfobases -.
func (cc *CtrlCache) GetInfobases(ctx context.Context, cluster entity.Cluster) ([]entity.Infobase, error) {
	var entities []entity.Infobase

	key := fmt.Sprintf(_keyInfobases, cluster.ID)

	if err := cc.Cache.Get(ctx, key, &entities); err != nil {
		return nil, fmt.Errorf("CtrlCache - GetInfobases - cc.Cache.Get: : %w", err)
	}

	return entities, nil
}

// PutInfobases -.
func (cc *CtrlCache) PutInfobases(ctx context.Context, cluster entity.Cluster, entities []entity.Infobase) error {
	if len(entities) == 0 {
		return nil
	}

	key := fmt.Sprintf(_keyInfobases, cluster.ID)

	return cc.Store(ctx, key, entities)
}

// GetSessions -.
func (cc *CtrlCache) GetSessions(ctx context.Context, cluster entity.Cluster, ib entity.Infobase) ([]entity.Session, error) {
	var entities []entity.Session

	key := fmt.Sprintf(_keySessions, cluster.ID, ib.ID)

	if err := cc.Cache.Get(ctx, key, &entities); err != nil {
		return nil, fmt.Errorf("CtrlCache - GetSessions - cc.Cache.Get: : %w", err)
	}

	return entities, nil
}

// PutSessions -.
func (cc *CtrlCache) PutSessions(ctx context.Context, cluster entity.Cluster, ib entity.Infobase, entities []entity.Infobase) error {
	if len(entities) == 0 {
		return nil
	}

	key := fmt.Sprintf(_keySessions, cluster.ID, ib.ID)

	return cc.Store(ctx, key, entities)
}

// Store -.
func (r *CtrlCache) Store(ctx context.Context, k string, v interface{}) error {
	if err := cc.Cache.Set(&c.Item{
		Ctx:   ctx,
		Key:   k,
		Value: v,
	}); err != nil {
		return fmt.Errorf("CtrlCache - Store - cc.Cache.Set: : %w", err)
	}

	return nil
}
