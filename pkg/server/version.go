package server

import (
	"github.com/longhorn/longhorn-share-manager/pkg/meta"
	"github.com/longhorn/longhorn-share-manager/pkg/rpc"
	"golang.org/x/net/context"
)

func (m *ShareManager) VersionGet(ctx context.Context, req *rpc.VersionRequest) (*rpc.VersionResponse, error) {
	v := meta.GetVersion()
	return &rpc.VersionResponse{
		Version:   v.Version,
		GitCommit: v.GitCommit,
		BuildDate: v.BuildDate,

		ApiVersion:    int64(v.APIVersion),
		ApiMinVersion: int64(v.APIMinVersion),
	}, nil
}
