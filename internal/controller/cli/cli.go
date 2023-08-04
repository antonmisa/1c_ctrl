package cli

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/antonmisa/1cctl/internal/entity"
	"github.com/antonmisa/1cctl/internal/usecase"
	"github.com/antonmisa/1cctl/pkg/logger"
)

type ctrl1CCLI struct {
	cp usecase.CtrlPipe
	cb usecase.CtrlBackup

	l logger.Interface
}

func New(cp usecase.CtrlPipe, cb usecase.CtrlBackup, l logger.Interface) *ctrl1CCLI {
	return &ctrl1CCLI{
		cp: cp,
		cb: cb,

		l: l,
	}
}

func (cc *ctrl1CCLI) Backup(clusterID string, infobaseID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check cluster exists
	cls, err := cc.cp.GetClusters(ctx)
	if err != nil {
		return fmt.Errorf("cli - Process - cc.cp.GetClusters: %w", err)
	}

	pos := -1
	for i := range cls {
		if cls[i].ID == clusterID {
			pos = i
			break
		}
	}

	if pos == -1 {
		return fmt.Errorf("cluster %s not found", clusterID)
	}

	cl := cls[pos]

	// Check infobase exists in cluster
	cl := entity.Cluster{ID: clusterID}
	ibs, err := cc.cp.GetInfobases(ctx, cl)
	if err != nil {
		return fmt.Errorf("cli - Process - cc.cp.GetInfobases: %w", err)
	}

	pos = -1
	for i := range ibs {
		if ibs[i].ID == infobaseID {
			pos = i
			break
		}
	}

	if pos == -1 {
		return fmt.Errorf("infobase %s not found in cluster %s", infobaseID, clusterID)
	}

	ib := ibs[pos]

	// Block all new sessions in infobase
	err = cc.cp.DisableSessions(ctx, ib)
	if err != nil {
		return fmt.Errorf("cli - Process - cc.cp.DisableSessions: %w", err)
	}

	// Get all sessions in infobase
	sessions, err := cc.cp.GetSessions(ctx, cl, ib)
	if err != nil {
		return fmt.Errorf("cli - Process - cc.cp.GetSessions: %w", err)
	}

	// Drop all sessions in infobase
	_ = cc.cp.DeleteSessions(ctx, cl, sessions)

	// Get all connections in infobase
	connections, err := cc.cp.GetConnections(ctx, cl, ib)
	if err != nil {
		return fmt.Errorf("cli - Process - cc.cp.GetConnections: %w", err)
	}

	// Drop all connections in infobase
	_ = cc.cp.DeleteConnections(ctx, cl, connections)

	// Run backup
	cx, cancel := context.WithTimeout(ctx, 60*time.Minute)
	err = cc.cb.RunBackup(cx, cl, ib, "admin", "pwd", "output")
	defer cancel()

	if err != nil {
		return fmt.Errorf("cli - Process - cc.cb.RunBackup: %w", err)
	}

	// Check backup exists
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("backup file does not exist at: %s", "output")
	}

	// UnBlock all sessions in infobase
	err = cc.cp.EnableSessions(ctx, ib)
	if err != nil {
		return fmt.Errorf("cli - Process - cc.cp.EnableSessions: %w", err)
	}

	return nil
}
