package server

import (
	"bufio"
	"os"

	"github.com/longhorn/longhorn-share-manager/pkg/rpc"
)

func (m *ShareManager) LogWatch(req *rpc.LogWatchRequest, srv rpc.ShareManagerService_LogWatchServer) error {
	m.logger.Debug("start streaming logs")
	file, err := os.OpenFile(m.logFile, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if err := srv.Send(&rpc.LogResponse{Line: line}); err != nil {
			return err
		}
	}

	m.logger.Debug("done streaming logs")
	return nil
}
