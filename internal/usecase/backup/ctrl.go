package backup

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/antonmisa/1cctl/internal/entity"
)

// CtrlBackup -.
type CtrlBackup struct {
	pathTo1C string
}

// New -.
func New(path string) (*CtrlBackup, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("1cv8 executable file does not exist: %w", err)
	}

	ctrl := &CtrlBackup{
		pathTo1C: path,
	}

	return ctrl, nil
}

// RunBackup -.
func (r *CtrlBackup) RunBackup(ctx context.Context,
	cl entity.Cluster, ib entity.Infobase,
	ibAdminName string, ibAdminPwd string,
	outputPath string) error {
	cmd := exec.CommandContext(ctx, r.pathTo1C, "CONFIG", "/S", fmt.Sprintf("%s:%s\\%s", cl.Host, cl.Port, ib.Name),
		"/N", ibAdminName, "/P", ibAdminPwd, "/UC", "12345", "/DisableStartupMessages",
		"/DumpIB", outputPath)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf(": %w", err)
	}
	defer cmd.Cancel()

	return nil
}
